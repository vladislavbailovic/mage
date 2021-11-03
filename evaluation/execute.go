package evaluation

import (
	"fmt"
	"mage/shell"
)

func Execute(tasks []task) ([]string, error) {
	outputs := []string{}
	for _, tsk := range tasks {
		for _, cmd := range tsk.getCommands() {
			command := shell.NewCommand([]string{cmd})
			out, err := command.GetStdout()
			if err != nil {
				return nil, fmt.Errorf("error executing command [%s] for task [%s]: %v", cmd, tsk.getName(), err)
			}
			outputs = append(outputs, out)
		}
	}
	return outputs, nil
}
