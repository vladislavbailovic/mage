package evaluation

import (
	"fmt"
	"mage/shell"
	"mage/typedefs"
	"strings"
)

func Compile(tasks []typedefs.Task) string {
	return strings.Join(GetCompiledStatements(tasks), "\n")
}

func GetCompiledStatements(tasks []typedefs.Task) []string {
	if len(tasks) == 0 {
		return []string{}
	}
	outputs := []string{
		fmt.Sprintf("#!%s", shell.GetShellBinary()),
		"",
	}
	for _, tsk := range tasks {
		outputs = append(outputs, fmt.Sprintf("# task [%s]", tsk.GetName()))
		for _, cmd := range tsk.GetCommands() {
			outputs = append(outputs, cmd)
		}
		outputs = append(outputs, "")
	}
	return outputs
}
