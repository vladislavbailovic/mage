Mage - task runner in go
========================


TODO:
-----

- [ ] Implement CLI runner
- [ ] Implement result reporting
- [ ] Spaces are significant (at least in command definitions)


Features:
---------

- Has includes
- Has macros
- Macros are recursive _in definition_
- Macros can call shell commands in macro definitions
- Has tasks
- Tasks have dependencies on other tasks or files
- Tasks cannot be redefined
- Execution order:
	1. Select a task:
		- If called with task argument, use that
		- Otherwise, use first task in the file
	2. Determine our epoch time:
		- Specific time
		- Task time (last called time for this specific task)
		- Build time (last time we built anything)
	3. Build evaluation stack
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
- Can compile scripts
- Selects initial task from command line, or first task in (preprocessed) source file
