
# qps.
set terminal png size 1024,768
set output 'qps.png'
set autoscale
set title "QPS"
set grid
set xlabel 'Time stamp'
set ylabel 'Requests'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "qps.txt" using 1:2 title "Clinet QPS", '' using 1:3 title "Server QPS"
replo

set terminal png size 1024,768
set output 'clientqps.png'
set autoscale
set title "QPS"
set grid
set xlabel 'Time stamp'
set ylabel 'Requests'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "qps.txt" using 1:2 title "Clinet QPS" 
replo

set terminal png size 1024,768
set output 'serverqps.png'
set autoscale
set title "QPS"
set grid
set xlabel 'Time stamp'
set ylabel 'Requests'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "qps.txt" using 1:3 title "Server QPS"
replo


set terminal png size 1024,768
set output 'hongbaoAmount.png'
set autoscale
set title "hongbao Amount"
set grid
set xlabel 'Time stamp'
set ylabel 'hongbao Amount'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "hongbaoAmount.txt" using 1:2 title "get hongbao Amount",'' using 1:3 title "create hongbao Amount" 
replo

set terminal png size 1024,768
set output 'createhongbao.png'

set autoscale
set title "createhongbao"
set grid
set xlabel 'Time stamp'
set ylabel 'create hongbao number'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "createhongbao.txt" using 1:2 title "create hongbao num"
replo

#
set terminal png size 1024,768
set output 'gethongbao.png'
set autoscale
set title "gethongbao"
set grid
set xlabel 'Time stamp'
set ylabel 'get hongbao number'
set xdata time
set timefmt "%H:%M:%S"
set style data lines
plot "gethongbao.txt" using 1:2 title "get hongbao num"
replo

