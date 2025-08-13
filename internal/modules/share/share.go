package share

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/core/mihomo"
	"github.com/bestruirui/bestsub/internal/core/node"
	"github.com/bestruirui/bestsub/internal/database/op"
	"github.com/bestruirui/bestsub/internal/models/share"
	"github.com/bestruirui/bestsub/internal/modules/subconverter"
	"github.com/bestruirui/bestsub/internal/utils/country"
	"github.com/google/go-querystring/query"
)

var (
	NodeData = []byte("proxies:\n")
	NewLine  = []byte("\n")
	Dash     = []byte(" - ")
)

type renameTmpl struct {
	SpeedUp   uint32
	SpeedDown uint32
	Delay     uint16
	Risk      uint8
	Country   country.Country
	Count     uint32
}

func GenSubData(genConfigStr string, userAgent string, token string, extraQuery string) []byte {
	var genConfig share.GenConfig
	if err := json.Unmarshal([]byte(genConfigStr), &genConfig); err != nil {
		return nil
	}
	nodeUrl := fmt.Sprintf("http://127.0.0.1:%d/api/v1/share/node/%s", config.Base().Server.Port, token)
	subUrlParam, _ := query.Values(genConfig.SubConverter)
	if genConfig.Proxy {
		subUrlParam.Add("proxy", op.GetConfigStr("proxy.url"))
	}
	subUrlParam.Add("url", nodeUrl)
	subUrlParam.Add("remove_emoji", "false")
	requestUrl := fmt.Sprintf("%s/sub?%s&%s", subconverter.GetBaseUrl(), subUrlParam.Encode(), extraQuery)
	client := mihomo.Default(false)
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("User-Agent", userAgent)
	response, err := client.Do(request)
	if err != nil {
		return nil
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	return body
}

func GenNodeData(config string) []byte {
	var genConfig share.GenConfig
	if err := json.Unmarshal([]byte(config), &genConfig); err != nil {
		return nil
	}
	nodes := node.GetByFilter(genConfig.Filter)
	var result bytes.Buffer
	result.Write(NodeData)
	tmpl, err := template.New("node").Parse(genConfig.Rename)
	if err != nil {
		return nil
	}
	var newName bytes.Buffer
	for i, node := range *nodes {
		newName.Reset()
		result.Write(Dash)
		simpleInfo := renameTmpl{
			SpeedUp:   node.Info.SpeedUp.Average(),
			SpeedDown: node.Info.SpeedDown.Average(),
			Delay:     node.Info.Delay.Average(),
			Risk:      node.Info.Risk,
			Count:     uint32(i + 1),
			Country:   country.GetCountry(node.Info.Country),
		}
		tmpl.Execute(&newName, simpleInfo)
		result.Write(replaceNameEnsureDQuoted(node.Base.Raw, newName.Bytes()))
		result.Write(NewLine)
	}
	return result.Bytes()
}

func yamlDQuoteBytes(b []byte) []byte {
	out := make([]byte, 0, len(b)+2)
	out = append(out, '"')
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case '\\':
			out = append(out, '\\', '\\')
		case '"':
			out = append(out, '\\', '"')
		case '\n':
			out = append(out, '\\', 'n')
		case '\t':
			out = append(out, '\\', 't')
		default:
			out = append(out, b[i])
		}
	}
	out = append(out, '"')
	return out
}

func replaceNameEnsureDQuoted(raw []byte, newName []byte) []byte {
	const key = "name:"
	i := bytes.Index(raw, []byte(key))
	if i < 0 {
		return raw
	}
	j := i + len(key)

	k := j
	for k < len(raw) && (raw[k] == ' ' || raw[k] == '\t') {
		k++
	}
	if k >= len(raw) {
		return raw
	}

	valStart := k
	valEnd := k

	switch raw[k] {
	case '"':
		p := k + 1
		for p < len(raw) {
			if raw[p] == '"' {
				bs := 0
				for q := p - 1; q >= k+1 && raw[q] == '\\'; q-- {
					bs++
				}
				if bs%2 == 0 {
					valEnd = p + 1
					break
				}
			}
			p++
		}
		if valEnd == k {
			return raw
		}
	case '\'':
		p := k + 1
		for p < len(raw) {
			if raw[p] == '\'' {
				if p+1 < len(raw) && raw[p+1] == '\'' {
					p += 2
					continue
				}
				valEnd = p + 1
				break
			}
			p++
		}
		if valEnd == k {
			return raw
		}
	default:
		rel := bytes.IndexByte(raw[k:], ',')
		if rel < 0 {
			return raw
		}
		valEnd = k + rel
	}

	quoted := yamlDQuoteBytes(newName)
	out := make([]byte, 0, len(raw)-(valEnd-valStart)+len(quoted))
	out = append(out, raw[:valStart]...)
	out = append(out, quoted...)
	out = append(out, raw[valEnd:]...)
	return out
}
