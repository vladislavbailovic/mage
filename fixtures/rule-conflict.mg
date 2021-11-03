:include ../fixtures/included.mg

root: feat3 feat2 feat1
	echo "All done!"

feat3:
	echo $(Feature3)
	echo "Last, obviously"

feat1:
	echo $(Feature1)
	echo "This comes 1st"

feat2:
	echo $(Feature2)
	echo "This is in the middle either way"

conflict-rule:
	echo "This should barf!"
