#!/bin/bash


for i in 182 183 184 185 186 187 188 189 190 191 193 194 195 155 157 158
do
    echo $i
    curl  "http://10.29.101.$i:9090/qps?qps=$1"
done
