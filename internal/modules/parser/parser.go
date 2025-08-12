package parser

import (
	"bytes"

	"github.com/bestruirui/bestsub/internal/models/node"
	"gopkg.in/yaml.v3"
)

func Parse(content *[]byte, subID uint16) (*[]node.Base, error) {
	var nodes []node.Base
	var unique node.UniqueKey
	lines := bytes.Split(*content, []byte("\n"))
	lines = lines[1:]
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		line = line[4:]
		if err := yaml.Unmarshal(line, &unique); err != nil {
			continue
		}
		nodes = append(nodes, node.Base{
			Raw:       line,
			SubId:     subID,
			UniqueKey: unique.Gen(),
		})
	}
	return &nodes, nil
}
