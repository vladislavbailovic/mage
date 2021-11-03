:macro NAME whatever goes here, ends up here
:macro OTHER $(M3)
:macro M3 $(M4)
:macro M4 $(NAME)
:macro M5 $(NAME)

root: tmp/whatever.test
	echo $(NAME)

tmp/whatever.test:
	echo nay nya $(OTHER)
	sed -e 's/$(M5)/nana/g'
