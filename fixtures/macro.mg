macro NAME whatever goes here, ends up here

root: tmp/whatever.test
	echo $(NAME)

tmp/whatever.test:
	echo nay nya $(NAME)
