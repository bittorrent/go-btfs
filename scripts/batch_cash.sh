#!/bin/bash

IFS_old=$IFS

cash_fee=27000000000000000000
withdraw_fee=15000000000000000000
wait_seconds=15

# for test
#cash_fee=27000000000
#withdraw_fee=15000000000
#wait_seconds=4


# check curl
if ! [ -x "$(command -v curl)" ]; then
    echo "curl is not found"
    echo "for details, please refer to this document: https://github.com/bittorrent/go-btfs/blob/master/docs/tutorial-upgrade-to-v2.1.2.md"
    exit 1
fi

# check bc
if ! [ -x "$(command -v bc)" ]; then
    echo "bc is not found"
    echo "for details, please refer to this document: https://github.com/bittorrent/go-btfs/blob/master/docs/tutorial-upgrade-to-v2.1.2.md"
    exit 1
fi

# set api host
echo -n ">> Btfs api host: "
btfs_api_host="127.0.0.1:5001"
if [ $# -gt 0 ]; then
    btfs_api_host=$1
fi
echo "$btfs_api_host"

# all apis
node_info_api="$btfs_api_host/api/v1/id"
cheque_stats_api="$btfs_api_host/api/v1/cheque/stats"
receive_cheques_api="$btfs_api_host/api/v1/cheque/receivelist"
vault_balance_api="$btfs_api_host/api/v1/vault/balance"
bttc_balance_api="$btfs_api_host/api/v1/cheque/bttbalance"
cheque_cash_api="$btfs_api_host/api/v1/cheque/cash"
withdraw_api="$btfs_api_host/api/v1/vault/withdraw"

# convert bigint to btt
function human_btt() {
    l=${#1}
    lct=$((l - 18))
    if [ $lct -gt 0 ]; then
        left=${1:0:$lct}
        rigt=${1:$lct:9}
        echo "$left.$rigt"
        return 0
    fi
    zct=$((18 - l))
    if [ $zct -ge 9 ]; then
        echo "0.000000000"
        return 0
    fi
    i=0
    rzeo=""
    while [ $i -lt $zct ]; do
        rzeo=$(echo "${rzeo}0")
        ((i++))
    done
    rnum=${1:0:$((9 - zct))}
    echo "0.$rzeo$rnum"
    return 0
}

# send post request
function post() {
    rsp=$(curl -s -XPOST $1)
    if [ $? -ne 0 ]; then
        echo "failed to access the host"
        return 1
    fi
    IFS=$'{},'
    lines=($rsp)
    i=0
    IFS=$':'
    unset vals
    while [ $i -lt ${#lines[@]} ]; do
        pairs=(${lines[$i]})
        j=2
        while [ $j -le $# ]; do
            field=$(eval echo '$'${j})
            if [ "${pairs[0]}" = "\"${field}\"" ]; then
                val=${pairs[1]//\"/}
                if [ "${val[0]}" = "\"" ]; then
                    val=${val:1:$((${#val} - 2))}
                fi
                if [ -n "$vals" ]; then
                    vals=$(echo "$vals\"$val")
                else
                    vals=$val
                fi
            fi
            ((j++))
        done
        ((i++))
    done
    if [ -z "$vals" ]; then
        echo "$rsp"
        return 1
    fi
    echo $vals
    return 0
}

# get node info
echo ""
echo ">> Node info: "
IFS=$' '
id_rsp=$(post $node_info_api "ID" "BttcAddress")
if [ "$?" -ne 0 ]; then
    echo -e "   error: $id_rsp"
    IFS=$IFS_old
    exit 1
fi
IFS=$'\"'
id_arr=($id_rsp)
my_peer_id=${id_arr[0]}
bttc_addr=${id_arr[1]}
echo -e "   peer_id: $my_peer_id"
echo -e "   bttc_address: $bttc_addr"

# stat received cheques
echo ""
echo ">> Received cheques stats: "
IFS=$' '
cash_count=0
cash_amount=0
cash_peers=()
cash_balances=()
stats_rsp=$(post $cheque_stats_api "total_received_count")
if [ $? -ne 0 ]; then
    echo -e "   error: $stat_rsp"
    IFS=$IFS_old
    exit 1
fi
list_rsp=$(post "$receive_cheques_api?arg=0&arg=$stats_rsp" "PeerID" "Payout" "CashedAmount")
if [ $? -ne 0 ]; then
    echo -e "   error: $list_rsp"
    IFS=$IFS_old
    exit 1
fi
IFS=$'\"'
list_arr=($list_rsp)
i=0
while [ $i -lt ${#list_arr[@]} ]; do
    peer_id=${list_arr[$((i++))]}
    payout=${list_arr[$((i++))]}
    cashed_amount=${list_arr[$((i++))]}
    balance=$(echo "$payout - $cashed_amount" | bc)
    echo -e "   ----------------------------------------"
    echo -e "   peer_id: $peer_id"
    echo -e "   payout: $(human_btt $payout) WBTT"
    echo -e "   cashed_amount: $(human_btt $cashed_amount) WBTT"
    echo -e "   balance: $(human_btt $balance) WBTT"
    balance_gt_0=$(echo "$balance > 0" | bc)
    if [ $balance_gt_0 -eq 1 ]; then
        balance_gt_fee=$(echo "$balance > $cash_fee" | bc)
        if [ $balance_gt_fee -eq 1 ]; then
            cash_peers[$cash_count]=$peer_id
            cash_balances[$cash_count]=$balance
            ((cash_count++))
            cash_amount=$(echo "$cash_amount + $balance" | bc)
            echo -e "   need_to_cash: 'Yes'"
        else
            echo -e "   need_to_cash: 'No, fee is too high'"
        fi
    else
        echo -e "   need_to_cash: 'No, zero balance'"
    fi
done

# assessment
echo ""
echo ">> Cash amounts and handling fee assessment: "
IFS=$' '
vault_balance_rsp=$(post $vault_balance_api "balance")
if [ $? -ne 0 ]; then
    echo -e "   error: $vault_balance_rsp"
    IFS=$IFS_old
    exit 1
fi
IFS=$' '
org_vault_balance=$vault_balance_rsp
bttc_balance_rsp=$(post "$bttc_balance_api?arg=$bttc_addr" "balance")
if [ $? -ne 0 ]; then
    echo -e "   error: $bttc_balance_rsp"
    IFS=$IFS_old
    exit 1
fi
org_bttc_balance=$bttc_balance_rsp
echo -e "   to_be_cashed_cheques: $cash_count"
echo -e "   to_be_cashed_cheques_amount: $(human_btt $cash_amount) WBTT"
es_withdraw=$(echo "$cash_amount + $org_vault_balance" | bc)
echo -e "   total_withdrawal_amount: $(human_btt $es_withdraw) WBTT"
es_with_gt_0=$(echo "$es_withdraw > 0" | bc)
if [ $es_with_gt_0 -eq 0 ]; then
    echo ""
    echo ">> Result: "
    echo -e "   Success, nothing to do!"
    IFS=$IFS_old
    exit 0
fi
echo -e "   --------"
handle_fee=$(echo "$cash_fee * $cash_count + $withdraw_fee" | bc)
echo -e "   estimated_handling_fee: $(human_btt $handle_fee) BTT"
echo -e "   current_bttc_balance: $(human_btt $org_bttc_balance) BTT"
echo ""
need_btt=$(echo "$handle_fee - $org_bttc_balance" | bc)
need_btt_gt_0=$(echo "$need_btt > 0" | bc)
if [ $need_btt_gt_0 -eq 1 ]; then
    echo "Your bttc balance is insufficient, please recharge $(human_btt $need_btt) BTT and try again!"
    echo "BTTC: $bttc_addr"
    IFS=$IFS_old
    exit 1
fi

# wait to start
echo ""
echo ">> Your balance is sufficient. Actions will start after 15 seconds. To cancel, press CTRL+C"
i=0
while [ $i -lt $wait_seconds ]; do
    c=$((wait_seconds - i))
    if [ $c -lt 10 ]; then
        c=$(echo "0$c")
    fi
    echo -ne "   seconds to start: ${c}s"
    sleep 1
    ((i++))
    echo -ne "\r"
done
echo ""

# cash cheques
echo ""
cash_vault_balance=$org_vault_balance
if [ ${#cash_peers[@]} -gt 0 ]; then
    echo ">> Cheques cashing: "
    i=0
    while [ $i -lt ${#cash_peers[@]} ]; do
        IFS=$' '
        cash_rsp=$(post "$cheque_cash_api?arg=${cash_peers[$i]}" "TxHash")
        if [ $? -eq 0 ]; then
            echo -e "   peer_id: ${cash_peers[$i]}, transacation: $cash_rsp"
        else
            echo -e "   peer_id: ${cash_peers[$i]}, error: $cash_rsp"
            IFS=$IFS_old
            exit 1
        fi
        ((i++))
    done
    i=0
    echo -e "   waiting for all transacations completed..."
    while [ $i -lt 10 ]; do
        IFS=$' '
        vault_balance_rsp=$(post $vault_balance_api "balance")
        if [ "$?" -eq 0 ]; then
            cash_vault_balance=$vault_balance_rsp
        else
            echo -e "   error: $vault_balance_rsp"
        fi
        cash_vault_balance_lt_wth=$(echo "$cash_vault_balance < $es_withdraw" | bc)
        if [ $cash_vault_balance_lt_wth -eq 1 ]; then
            ((i++))
            sleep 1
        else
            i=10
        fi
    done
fi

# withdraw from vault
echo ""
with_vault_balance=$cash_vault_balance
with_vault_balance_gt_0=$(echo "$with_vault_balance > 0" | bc)
if [ $with_vault_balance_gt_0 -eq 1 ]; then
    echo ">> Vault balance withdraw: "
    echo -e "   withdraw_amount: $(human_btt $cash_vault_balance) WBTT"
    withdraw_rsp=$(post "$withdraw_api?arg=$cash_vault_balance" "hash")
    if [ $? -eq 0 ]; then
        echo -e "   transacation: $withdraw_rsp"
    else
        echo -e "   error: $withdraw_rsp"
        IFS=$IFS_old
        exit 1
    fi
    i=0
    echo -e "   waiting for the transacation completed..."
    while [ $i -lt 10 ]; do
        IFS=$' '
        vault_balance_rsp=$(post $vault_balance_api "balance")
        if [ "$?" -eq 0 ]; then
            with_vault_balance=$vault_balance_rsp
        else
            echo -e "   error: $vault_balance_rsp"
        fi
        with_vault_balance_gt_0=$(echo "$with_vault_balance > 0" | bc)
        if [ $with_vault_balance_gt_0 -eq 1 ]; then
            ((i++))
            sleep 1
        else
            i=10
        fi
    done
fi

# print result
echo ""
echo ">> Result: "
not_cash_all=$(echo "$cash_vault_balance < $es_withdraw" | bc)
not_with_all=$(echo "$with_vault_balance > 0" | bc)
if [ $not_cash_all -eq 1 -o $not_with_all -eq 1 ]; then
    echo "   Some transactions are still not completed, you can try again!"
else
    echo "   Success, all tasks completed!"
fi

IFS=$IFS_old
exit 0
