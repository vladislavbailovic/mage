macro NAME whatever goes here, ends up here
macro OTHER $(NAME)

root: tmp/whatever.test
	echo $(NAME)

tmp/whatever.test:
	echo nay nya $(OTHER)
