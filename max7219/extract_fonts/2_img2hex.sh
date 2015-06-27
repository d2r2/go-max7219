#!/bin/bash
i=0
k=0
char=""
for line in `./img2bin.sh $1 | sort -k1 -n`
do
    #echo $line
    a=(${line//;/ })
    #echo "key: ${a[0]}, value: ${a[1]}"
    hex=$(printf "%02X" "$((2#${a[1]}))")
    if [ -z "$char" ]; then
        char="0x$hex"
    else
        char="$char, 0x$hex"
    fi
    if [ $k -eq 7 ]; then
        printf "{%s}, // 0x%02X\n" "$char" $i
        char=""
        i=$(( i+1 ))
        k=0
    else
        k=$(( k+1 ))
    fi

done
