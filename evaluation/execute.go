package evaluation

import (
	"fmt"
	"mage/shell"
	"mage/typedefs"
)

func Execute(tasks []typedefs.Task) ([]string, error) {
	outputs := []string{}
	for _, tsk := range tasks {
		for _, cmd := range tsk.GetCommands() {
			command := shell.NewCommand([]string{cmd})
			out, err := command.GetStdout()
			if err != nil {
				return nil, fmt.Errorf("error executing command [%s] for task [%s]: %v", cmd, tsk.GetName(), err)
			}
			outputs = append(outputs, out)
			tsk.RecordTime()
		}
	}
	return outputs, nil
}
