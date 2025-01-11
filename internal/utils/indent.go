package utils

import (
	"strings"
)

func YamlIndent(yaml []byte, indent int) string {
	indentSpace := strings.Repeat(" ", indent)
	yamlstring := string(yaml)
	indentedYaml := strings.ReplaceAll(yamlstring, "    ", indentSpace)

	if len(indentedYaml) > 2 && indentedYaml[2:] == "\n" {
		indentedYaml = indentedYaml[:len(indentedYaml)-1]
	}
	return indentedYaml
}
