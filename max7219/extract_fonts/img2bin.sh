#!/bin/bash

width=$(identify -format "%w" $1)
declare -A a
for line in `./img2txt.sh $1`
do
    xyv=(${line//;/ })
    if [ ${#xyv[@]} -eq 3 ]; then
        y=${xyv[1]}
        x=${xyv[0]}
        val=${xyv[2]}
        # echo "x=$x, y=$y, val=$val"
        ind=$(( y/8*width+x ))
        index=$(printf "%04d" $ind)
        # echo "index: $index"
        if [ -z ${a[$index]} ]; then
            a[$index]=$val
        else
            a[$index]=$val${a[$index]}
        fi
        #echo "value: ${a[$index]}"
    fi
done

for bt in "${!a[@]}"
do
    echo "$bt;${a[$bt]}"
done
