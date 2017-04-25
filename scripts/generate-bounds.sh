#!/bin/bash
#
# Take the tsl parameter map and find boundaries to keep in a compact setting.
# $1 = tolerance for wasted bytes
# $2 = tsl parameter map CSV file

OLDIFS=$IFS
IFS=";"

T=$1
START=-1
END=-1

function print_ln(){
    echo "$START;$END"
}

while read tlsname offset size relevant
 do
    # For entries that will be used in the Katana.
    if(($relevant == 1)); then
        
        # Means we haven't started on a boundary yet.
        if(($START == -1)); then
            START=$offset
            END=$(($offset + $size - 1))
        
        # Moving the boundary end, since we're within tolerance.
        elif(($END + $T >= $offset)); then
            END=$(($offset + $size - 1))
        
        # That makes it out of bounds.
        else
            print_ln
            START=$offset
            END=$(($offset + $size - 1))
        fi
    fi
 done < $2

(($START > 0)) && print_ln
IFS=$OLDIFS