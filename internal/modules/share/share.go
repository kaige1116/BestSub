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
	"github.com/bestruirui/bestsub/internal/models/setting"
	"github.com/bestruirui/bestsub/internal/models/share"
	"github.com/bestruirui/bestsub/internal/modules/subcer"
	"github.com/bestruirui/bestsub/internal/utils/country"
	"github.com/google/go-querystring/query"
)

func GenSubData(genConfigStr string, userAgent string, token string, extraQuery string) []byte {
	var genConfig share.GenConfig
	if err := json.Unmarshal([]byte(genConfigStr), &genConfig); err != nil {
		return nil
	}
	subUrlParam, _ := query.Values(genConfig.SubConverter)
	if genConfig.Proxy {
		subUrlParam.Add("config_proxy", op.GetSettingStr(setting.PROXY_URL))
	}
	subUrlParam.Add("url", fmt.Sprintf("http://127.0.0.1:%d/api/v1/share/node/%s", config.Base().Server.Port, token))
	subUrlParam.Add("remove_emoji", "false")
	subcer.RLock()
	defer subcer.RUnlock()
	requestUrl := fmt.Sprintf("%s/sub?%s&%s", subcer.GetBaseUrl(), subUrlParam.Encode(), extraQuery)
	client := mihomo.Default(false)
	if client == nil {
		return nil
	}
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
	result.Write(nodeData)
	tmpl, err := template.New("node").Parse(genConfig.Rename)
	if err != nil {
		return nil
	}
	var newName bytes.Buffer
	for i, node := range *nodes {
		newName.Reset()
		result.Write(dash)
		simpleInfo := renameTmpl{
			SpeedUp:   node.Info.SpeedUp.Average(),
			SpeedDown: node.Info.SpeedDown.Average(),
			Delay:     node.Info.Delay.Average(),
			Risk:      node.Info.Risk,
			Count:     uint32(i + 1),
			Country:   country.GetCountry(node.Info.Country),
		}
		tmpl.Execute(&newName, simpleInfo)
		result.Write(rename(node.Base.Raw, newName.Bytes()))
		result.Write(newLine)
	}
	return result.Bytes()
}

func rename(raw []byte, newName []byte) []byte {
	idx := bytes.Index(raw, serverDelim)
	if idx < 0 {
		return raw
	}
	out := make([]byte, 0, len(name)+len(newName)+len(raw)-idx)
	out = append(out, name...)
	out = append(out, newName...)
	out = append(out, raw[idx:]...)
	return out
}

var (
	name        = []byte("{name: ")
	serverDelim = []byte(", server:")

	nodeData = []byte("proxies:\n")
	newLine  = []byte("\n")
	dash     = []byte(" - ")
)

type renameTmpl struct {
	SpeedUp   uint32
	SpeedDown uint32
	Delay     uint16
	Risk      uint8
	Country   country.Country
	Count     uint32
}
