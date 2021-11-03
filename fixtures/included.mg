:macro Feature1 Has macros
:macro Feature2 Has task definitions
:macro Feature3 Has includes?
:macro Features $(Feature1); $(Feature2); $(Feature3) Now it does

conflict-rule:
	echo $(Features)
	# note, rule name should trigger conflict.
	echo Erase ^that and me when done
	echo This line should be kept though
