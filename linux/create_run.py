#!/usr/bin/env python

import sys
i=0
Max= 60*10*1000
Range = 6*10*1000
Qps = 30

index = -1
if len(sys.argv) != 1:
    index = int(sys.argv[1])

printindex = 0
while True:
    if i >= Max-1:
        break
    if i+Range>= Max-1:

        if index == -1 or index == printindex:
            print "./client -min", i, "-max", Max-1, "-total",Max, "-qps",Qps
        i = Max -1
        printindex += 1
    else:
        if index == -1 or index == printindex:
            print "./client -min", i, "-max", i+Range -1 , "-total",Max, "-qps", Qps
        i = i + Range
        printindex += 1

