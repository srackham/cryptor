#
# Plot pie chart of portfolio currency allocation by value.
# This script is called from `examples/plot-valuation.sh`
# Data file name passed as variable `data`.
# Plot title passed as variable `title`.
#

set datafile separator comma
date_fmt="%Y-%m-%d"
set size square
set style fill solid 1
unset border
unset tics
unset key
stats data using 4 noout
set title title font 'Helvetica,16' offset -4,-3
set yrange [-1.25:1.25]
ang(x) = x*360.0/100.0
Ai = 0.0; Bi = 0.0;             # init angle
mid = 0.0;                      # mid angle
i = 0; j = 0;                   # color
yi  = 0.0; yi2 = 0.0;           # label position
plot data using (0):(0):(1):(Ai):(Ai=Ai+ang($4)):(i=i+1) with circle linecolor var,\
     data using (1.5):(yi=yi+0.5/STATS_records):(sprintf('%s (%.0f\%)', stringcolumn(3), $4)) w labels left offset -3,0,\
     data using (1.3):(yi2=yi2+0.5/STATS_records):(j=j+1) w p pt 5 ps 2 linecolor var,\
     data using (mid=Bi+ang($4)*pi/360.0, Bi=2.0*mid-Bi, 0.5*cos(mid)):(0.5*sin(mid)):($4<1.0 ? "" : sprintf('%s (%.0f\%)', stringcolumn(3), $4)) w labels