#!/bin/bash
IFS_old=$IFS

#1.Set the btfs command, the default is "btfs"
echo ">> 1.Set the btfs command, the default is 'btfs'"
btfs_cmd="btfs"
if [ $# -ge 1 ]; then
    btfs_cmd=$1
fi
if ! [ -x "$(command -v $btfs_cmd)" ]; then
    echo "btfs command not found"
    exit 1
else
    echo "current btfs command: '$btfs_cmd'"
fi
echo ""

#2.Query all received cheques 
echo ">> 2.Query all received cheques"
receive_rsp=$($btfs_cmd cheque receivelist 0 100000)
if [ $? -ne 0 ]; then
    echo  "query failed"
    exit 1
else
    IFS=$' \t\n'
    receive_arr=($receive_rsp)
    cheques_count=$((${#receive_arr[@]}/5-1))
    echo "received cheques: $cheques_count"
fi
echo ""

#3.Query vault balance before cashing out" 
echo ">> 3.Query vault balance before cashing out"
balance_rsp=$($btfs_cmd vault balance)
if [ $? -ne 0 ]; then
    echo -e "get vault balance failed"
    exit 1
fi
IFS=':'
balance_rsp_arr=($balance_rsp)
before_balance=$((${balance_rsp_arr[1]}+0))
echo "before balance:$before_balance"
echo ""

#4.Cash out all outstanding cheques
echo ">> 4.Cash out all outstanding cheques"
i=5
total_cash=0
succ_cash=0
fail_cash=0
cashed=0
while [ $i -lt ${#receive_arr[@]} ]
do
    peer_id=${receive_arr[i++]}
    vault=${receive_arr[$((i++))]}
    beneficiary=${receive_arr[$((i++))]}
    cashout_amount=${receive_arr[$((i++))]}
    amount=${receive_arr[$((i++))]}
    uncashed=$((cashout_amount-amount))
    echo "-----"
    echo -e "peer id: $peer_id"
    echo -e "cheque amount: $cashout_amount"
    echo -e "cashed: $amount"
    echo -e "uncashed: $uncashed"
    if [ $amount -lt $cashout_amount ]; then
        ((total_cash++))
        cash_rsp=$($btfs_cmd cheque cash $peer_id)
        if [ $? -ne 0 ]; then
            echo -e "cash out failed"
            ((fail_cash++))
        else
            ((cashed+=uncashed))
            ((succ_cash++))
            echo "cash out success: $cash_rsp"
        fi
    else
        echo "no amount to cash"
    fi
done
echo ""

#5.Wait for all cheques cashing out amount to arrive
echo ">> 5.Wait for all cheques cashing out amount to arrive"
IFS=':'
retries=1
current_balance=0
should_balance=$((before_balance+cashed))
while [ $retries -le 10 -a $current_balance -lt $should_balance ]
do
    balance_rsp=$($btfs_cmd vault balance)
    if [ $? -ne 0 ]; then
        echo -e "retries: $retries, get vault balance failed"
    else
        balance_rsp_arr=($balance_rsp)
        current_balance=$((${balance_rsp_arr[1]}+0))
        echo "retries: $retries, current balance: $current_balance, should be: $should_balance"
    fi
    sleep 2
    ((retries++))
done
echo ""

#6.Withdraw current balance
echo ">> 6.Withdraw current balance"
echo "balance: $current_balance"
if [ $current_balance -gt 0 ]; then
    withdraw_rsp=$($btfs_cmd vault withdraw $current_balance)
    if [ $? -ne 0 ]; then
        echo -e "withdraw fail"
        exit 1
    else
        echo "withdraw success: $withdraw_rsp"
    fi
else
    echo "no balance to withdraw"
fi
echo ""

#7.Wait for all cheques cashing out amount to arrive
echo ">> 7.Wait for withdraw completed"
IFS=':'
retries=1
after_balance=$current_balance
while [ $retries -le 10 -a $after_balance -gt 0 ]
do
    balance_rsp=$($btfs_cmd vault balance)
    if [ $? -ne 0 ]; then
        echo -e "retries: $retries, get vault balance failed"
    else
        balance_rsp_arr=($balance_rsp)
        after_balance=$((${balance_rsp_arr[1]}+0))
        echo "retries: $retries, after balance: $after_balance"
    fi
    sleep 2
    ((retries++))
done
echo ""

#8.Print result
echo ">> 8.Result"
echo -e "cheques to be cashed: $total_cash"
echo -e "cheques successfully cashed: $succ_cash"
echo -e "cheques failed to cash: $fail_cash"
echo -e "cashed amount: $cashed"
echo "-----"
echo -e "balance to withdraw: $should_balance"
echo -e "balance withdrawn: $current_balance"
echo -e "balance unwithdrawn: $after_balance"
echo ""

if [ $fail_cash -gt 0  -o  $current_balance -lt $should_balance -o $after_balance -gt 0 ]; then
    echo ">> some stages failed, please retry!"
else
    echo ">> all stages success, congratulate!"
fi
echo ""

IFS=$IFS_old
exit 0
