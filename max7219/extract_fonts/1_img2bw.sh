#!/bin/bash
# convert $1 -level 75%,100% -type bilevel $1_bw.png
convert $1 -negate -level 75%,100% -type bilevel $1_bw.png
