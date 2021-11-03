Mage - task runner in go
========================


TODO:
-----

- [X] Implement evaluation stack creation
- [X] Implement evaulation stack processing
- [ ] Implement initial task selection
- [ ] Implement CLI runner
- [ ] Implement result reporting
- [ ] Spaces are significant (at least in command definitions)


Features:
---------

- Has includes
- Has macros
- Macros are recursive _in definition_
- Has tasks
- Tasks have dependencies on other tasks or files
- Tasks cannot be redefined
- Execution order:
	0. Select a task:
		- If called with task argument, use that
		- Otherwise, use first task in the file
	1. Determine our epoch time:
		- Specific time
		- Task time (last called time for this specific task)
		- Build time (last time we built anything)
	2. Build evaluation stack
- Task evaluation order:
	1. Get a list of dependencies:
		- Evaluate dependencies in reverse order they're listed
	2. Determine dependency milestone time:
		- If the dependency is another task, use its last run time
		- If the dependency is a file, use its last modified time
	3. Determine dependency age compared with our epoch time
	4. If dependency is new in our epoch, select it as task and goto 1
	5. Add task content to execution stack
- Execute stack
