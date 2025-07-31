package mihomo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bestruirui/bestsub/internal/models/node"
	"github.com/bestruirui/bestsub/internal/utils/log"
	"github.com/cespare/xxhash/v2"
	"github.com/goccy/go-yaml"
)

func Parse(content *[]byte, sublinkID uint16) (*[]node.Data, error) {

	var inProxiesSection bool
	var yamlBuffer bytes.Buffer
	var indent int
	var isFirst bool = true
	yamlBuffer.Grow(1024)

	contentReader := bytes.NewReader(*content)
	scanner := bufio.NewScanner(contentReader)
	var nodes []node.Data

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "proxies:" {
			inProxiesSection = true
			continue
		}

		if !inProxiesSection {
			continue
		}

		if isFirst {
			indent = len(line) - len(trimmedLine)
			isFirst = false
		}

		if len(line)-len(trimmedLine) == 0 && !strings.HasPrefix(trimmedLine, "-") && trimmedLine != "" {
			break
		}

		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		if strings.HasPrefix(trimmedLine, "-") && len(line)-len(trimmedLine) == indent {
			if yamlBuffer.Len() > 0 {
				yamlBytes := yamlBuffer.Bytes()
				yamlBytes[indent] = ' '
				if err := parseProxyNode(&yamlBytes, &nodes); err != nil {
					log.Errorf("parseProxyNode error: %v", err)
					break
				}

				yamlBuffer.Reset()
			}
			yamlBuffer.WriteString(line + "\n")
		} else if yamlBuffer.Len() > 0 {
			yamlBuffer.WriteString(line + "\n")
		}
	}

	if yamlBuffer.Len() > 0 {
		yamlBytes := yamlBuffer.Bytes()
		yamlBytes[indent] = ' '
		if err := parseProxyNode(&yamlBytes, &nodes); err != nil {
			log.Errorf("parseProxyNode error: %v", err)
		}
	}

	return &nodes, nil
}

func parseProxyNode(nodeYAML *[]byte, nodes *[]node.Data) error {
	mihomoConfig := map[string]any{}
	if err := yaml.Unmarshal(*nodeYAML, &mihomoConfig); err != nil {
		return fmt.Errorf("failed to unmarshal to config struct: %v", err)
	}
	jsonBytes, err := json.Marshal(&mihomoConfig)
	if err != nil {
		return fmt.Errorf("failed to convert to JSON: %v", err)
	}
	nodeData := &node.Data{
		Raw: jsonBytes,
		Info: node.Info{
			UniqueKey: generateUniqueKey(&mihomoConfig),
			AddTime:   uint64(time.Now().Unix()),
		},
	}
	*nodes = append(*nodes, *nodeData)
	return nil
}

func generateUniqueKey(mihomoConfig *map[string]any) uint64 {
	h := xxhash.New()
	h.Write([]byte(fmt.Sprintf("%v%v%v%v%v%v%v",
		(*mihomoConfig)["server"],
		(*mihomoConfig)["servername"],
		(*mihomoConfig)["port"],
		(*mihomoConfig)["type"],
		(*mihomoConfig)["uuid"],
		(*mihomoConfig)["username"],
		(*mihomoConfig)["password"],
	)))
	return h.Sum64()
}
