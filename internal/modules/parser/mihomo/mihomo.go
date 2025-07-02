package mihomo

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/bestruirui/bestsub/internal/core/nodepool"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// Parse 解析mihomo配置内容
func Parse(content *[]byte, sublinkID int64) (int, error) {

	var inProxiesSection bool
	var yamlBuffer bytes.Buffer
	var indent int
	var isFirst bool = true
	var addedCount int
	yamlBuffer.Grow(1024)

	contentReader := bytes.NewReader(*content)
	scanner := bufio.NewScanner(contentReader)
	collection := &nodepool.Collection{}

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
				if err := parseProxyNode(&yamlBytes, collection); err != nil {
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
		if err := parseProxyNode(&yamlBytes, collection); err != nil {
			log.Errorf("parseProxyNode error: %v", err)
		}
	}

	addedCount, err := nodepool.Add(collection, sublinkID)

	return addedCount, err
}
