#!/bin/bash
IFS_old=$IFS

# set btfs cmd 
btfsCmd="btfs"
if [ $# -ge 1 ]; then
    btfsCmd=$1
fi
if ! [ -x "$(command -v $btfsCmd)" ]; then
    echo "btfs command not found, usage: $0 <btfs-command>"
    echo ".e.g.\"$0 ./btfs\""
    exit 1
fi


# Stage 1: cash out all cheques
echo "> Stage 1: cash out all your cheques"

# 1.1 get current balance and it will be used to check if cheque cashing out transaction completed
echo "*Get current vault balance..."
balanceRsp=$($btfsCmd vault balance)
if [ $? -ne 0 ]; then
    echo -e "get vault balance fail!"
    exit 1
fi
IFS=':'
balanceRspArr=($balanceRsp)
willBalance=${balanceRspArr[1]}

# 1.2 get receive cheques list
echo "*Get receive cheques list..."
receivelistRsp=$(./btfs cheque receivelist 0 100000)
if [ $? -ne 0 ]; then
    echo -e "get cheque receive list fail!"
    exit 1
fi
IFS=$' \t\n'
receivelistRspArr=($receivelistRsp)

# 1.3 cash out receive cheques
echo "*Cash out all cheques..."
i=5
totalCash=0
succCash=0
failCash=0
while [ $i -lt ${#receivelistRspArr[@]} ]
do
    peerID=${receivelistRspArr[$i]}
    vault=${receivelistRspArr[$(($i+1))]}
    beneficiary=${receivelistRspArr[$(($i+2))]}
    cashout_amount=${receivelistRspArr[$(($i+3))]}
    amount=${receivelistRspArr[$(($i+4))]}
    remain=$((cashout_amount-amount))
    echo -e "--------------\n perrID: $peerID\n cashout_amount: $cashout_amount\n amount: $amount\n remain: $remain"
    echo "--------------"
    if [ $amount -lt $cashout_amount ]; then
        ((totalCash++))
        cashRsp=$(./btfs cheque cash $peerID)
        if [ $? -ne 0 ]; then
            echo -e "cash out fail!"
            ((failCash++))
        else
            ((willBalance+=remain))
            ((succCash++))
            echo "cash out success: $cashRsp"
        fi
    else
        echo "no amount to cash out!"
    fi
    ((i+=5))
done

# 1.4 waiting for all cheques withdrawal transactions to complete
echo "*Waiting for all cheques withdrawal transactions to complete..."
IFS=':'
balance=0
retries=1
while [ $retries -le 10 -a $balance -lt $willBalance ]
do
    balanceRsp=$(./btfs vault balance)
    if [ $? -ne 0 ]; then
        echo -e "get vault balance fail, retry: $retries"
    else
        balanceRspArr=($balanceRsp)
        balance=${balanceRspArr[1]}
        echo "retry: $retries, balance: $balance, should: $willBalance"
    fi
    sleep 2
    ((retries++))
done

# Stage 2: withdraw all balance
echo ""
echo "> Stage 2: withdraw from vault to your bttc"

# 2.1 withdraw all balance
echo "*Withdraw all balance..."
echo "balance need to withdraw: $balance"
if [ $balance -gt 0 ]; then
    withdrawRsp=$(./btfs vault withdraw $balance)
    if [ $? -ne 0 ]; then
        echo -e "withdraw fail!"
        exit 1
    else
        echo "withdraw success: $withdrawRsp"
    fi
fi

# Stage 3: print result
echo ""
echo "> Result"
echo "cheques cash out:"
echo " total: $totalCash"
echo " success: $succCash"
echo " fail: $failCash"
echo "blanche withdraw:"
echo " should: $willBalance"
echo " withdraw: $balance"
if [ $failCash -gt 0  -o  $balance -lt $willBalance ]; then
    echo "some stages failed, please retry!"
else
    echo "all stage success, congratulate!"
fi

IFS=$IFS_old
exit 0
