# Simple magefile
root: dependency2 dependency1
	# commands to build root task
dependency2: dep-dependency1
	# subtask 2 dependes on dep-dependency1
dependency1: parser.go
	# dependency1 depends on parser.go
parser.go:
	# no dependencies here, but needs to be created
dep-dependency1: tmp/not-created-yet.go
	# no dependencies, but has to be done to allow dep2
tmp/not-created-yet.go: tmp
	touch tmp/not-created-yet.go
tmp:
	mkdir tmp
