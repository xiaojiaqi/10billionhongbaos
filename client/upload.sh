#!/bin/bash

index=0
for i in 182 183 184 185 186 187 188 189 190 191 193 194 195 155 156 157 158
do
    echo $i
    scp client 10.29.101.$i:~/
done
