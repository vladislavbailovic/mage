# should be 2021-11-04
:macro DATE $(!date -d@1636047105 +%Y-%m-%d)
:macro SEQ $(!seq 13 12 161)
:macro FU $(!for i in `seq 13 12 161`; do echo -n "number $i;"; done)

root:
	echo $(DATE)
	# $(FU)
	# $(SEQ)
