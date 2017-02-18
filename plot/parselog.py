#!/usr/bin/env python

import sys
import datetime

maps1 = {} # for client
maps2 = {} # for server
times = []

def getTime(i):
    return datetime.datetime.fromtimestamp(int(i)).strftime('%H:%M:%S')


for line in sys.stdin:
    v = line.split()
    print v
    if v[1] == "0":
          maps1[v[0]] = v
    if v[1] == "1":
          maps2[v[0]] = v
          times.append(v[0])

# client and server qps
f = open('qps.txt', 'w')
for i in times:
    s = "" + getTime(i) + " "
    if maps1.has_key(i):
        v = maps1[i]
        s += v[3] + " "
    else:
        s += "0 "
    if maps2.has_key(i):
        v = maps2[i]
        s += v[3] + " "
    else:
        s += "0 "
    print >> f, s
f.close()

#hongbaoAmount.png
f = open('hongbaoAmount.txt', 'w')
for i in times:
    s = "" + getTime(i) + " "
    if maps1.has_key(i):
        v = maps1[i]
        s += v[17] + " "
    else:
        s += "0 "
    if maps2.has_key(i):
        v = maps1[i]
        s += v[17] + " "
    else:
        s += "0 "
    print >> f, s
f.close()

# create hongbao
hongbaonum = 0
f = open('createhongbao.txt', 'w')
for i in times:
    s = "" + getTime(i) + " "
    if maps2.has_key(i):
        v = maps2[i]
        #num = int(v[15].split("/")[0])
        #strnum = "0"
        #if num != 0:
        #   strnum = str(num - hongbaonum)
        #   hongbaonum = num
        #s += strnum + " "
        v = maps2[i]
        s += v[7] + " "
    else:
        s += "0 "
    print >> f, s
f.close()

# get hongbao 
f = open('gethongbao.txt', 'w')
for i in times:
    s = "" + getTime(i) + " "
    if maps1.has_key(i):
        v = maps1[i]
        s += v[7] + " "
    else:
        s += "0 "
    print >> f, s
f.close()

