#
# Plot pie chart of portfolio currency allocation by value.
# This script is called from `examples/plot-valuation.sh`
# Data file name passed as variable `data`.
# Plot title passed as variable `title`.
#

# set terminal qt size 1024,768
set datafile separator comma
date_fmt="%Y-%m-%d"
stats data using (strptime(date_fmt,stringcolumn(2))) noout
max_date=STATS_max # The most recent dataset date.

cost = system('cat "'.data.'" | tail -1 | awk -F "," "{print \$3}"') + 0
value = system('cat "'.data.'" | tail -1 | awk -F "," "{print \$4}"') + 0
roi = (value - cost)/cost*100

set xdata time
set timefmt date_fmt
set format x "%d-%b-%y"
set xtics nomirror rotate by 45 right font ", 8"
set datafile missing NaN
set decimalsign locale  # To ensure thousands separator in formated strings.
set multiplot layout 2,1

set title '{/:Bold Value '.sprintf("$%'d", value).' '.strftime("%d-%b-%Y", max_date).'}'
set ylabel 'Value (USD)'
set yrange [0:]
set ytics 0, 100 format "$%'.0fK" nomirror font ", 8"
plot \
    data \
        using 2:($4/1000) \
        with linespoints pointtype 7 title ''

set title '{/:Bold ROI '.sprintf("%.0f%% ($%'d)",roi, value - cost).' '.strftime("%d-%b-%Y", max_date).'}'
set ylabel 'Percent ROI'
set yrange [-100:]
set ytics -100, 50 format "%.0f%%" nomirror
plot \
    0 notitle linetype 8 dashtype 3, \
\
    data \
        using 2:(($4-$3)/$3*100) \
        with linespoints pointtype 7 title ''