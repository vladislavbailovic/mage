package processor

// Processes a set of pre-processed tokens into
// intermediate ruleset representation

import (
	"strings"

	"mage/typedefs"
)

func process(tokens []typedefs.Token) (map[string]typedefs.TaskDefinition, error) {
	result := map[string]typedefs.TaskDefinition{}

	commands := []string{}
	currentCommand := []string{}
	rulePos := typedefs.SourcePosition{}
	ruleName := ""
	dependencies := []string{}
	for i := 0; i < len(tokens); i++ {

		if tokens[i].Kind == typedefs.TOKEN_COMMAND_OPEN {
			i += 1
			currentCommand = []string{}
			for j := i; j < len(tokens); j++ {
				if tokens[j].Kind == typedefs.TOKEN_COMMAND_CLOSE {
					cmd := strings.Join(currentCommand, " ")
					commands = append(commands, cmd)
					break
				}
				if tokens[j].Kind != typedefs.TOKEN_WORD {
					return nil, tokenError(tokens[j], "only words allowed in commands")
				}
				currentCommand = append(currentCommand, tokens[j].Value)
			}
			i += 1
			continue
		}

		if tokens[i].Kind == typedefs.TOKEN_RULE_OPEN {
			if len(ruleName) > 0 && (len(commands) > 0 || len(dependencies) > 0) {
				// Add old rule
				result[ruleName] = typedefs.TaskDefinition{
					rulePos,
					ruleName, ruleName,
					dependencies,
					commands,
				}
				commands = []string{}
				dependencies = []string{}
			}

			rulePos = tokens[i].Pos
			i += 1
			if tokens[i].Kind != typedefs.TOKEN_WORD {
				return nil, tokenError(tokens[i], "rule name has to be a word")
			}
			ruleName = tokens[i].Value
			i += 1
			for j := i; j < len(tokens); j++ {
				if tokens[j].Kind == typedefs.TOKEN_RULE_CLOSE {
					break
				}
				if tokens[j].Kind != typedefs.TOKEN_WORD {
					return nil, tokenError(tokens[j], "dependency not a word")
				}
				dependencies = append(dependencies, tokens[j].Value)
			}
			continue
		}
	}

	if len(ruleName) > 0 && (len(commands) > 0 || len(dependencies) > 0) {
		// Add old rule
		result[ruleName] = typedefs.TaskDefinition{
			rulePos,
			ruleName, ruleName,
			dependencies,
			commands,
		}
	}

	return result, nil
}
