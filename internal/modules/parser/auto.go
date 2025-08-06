package parser

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/bestruirui/bestsub/internal/models/node"
	"github.com/bestruirui/bestsub/internal/modules/parser/mihomo"
	"github.com/bestruirui/bestsub/internal/modules/parser/singbox"
	"github.com/bestruirui/bestsub/internal/modules/parser/v2ray"
	"github.com/bestruirui/bestsub/internal/utils"
)

func auto(content *[]byte ) (*[]node.Data, error) {
	utils.RemoveAllControlCharacters(content)

	contentStr := strings.TrimSpace(string(*content))
	if len(contentStr) == 0 {
		return nil, fmt.Errorf("content is empty after cleaning")
	}

	if isSingBoxFormat(&contentStr) {
		return singbox.Parse(content)
	}

	if isMihomoFormat(&contentStr) {
		return mihomo.Parse(content)
	}

	if isBase64Encoded(&contentStr) {
		decoded, err := base64.StdEncoding.DecodeString(cleanBase64String(contentStr))
		if err == nil {
			decodedStr := strings.TrimSpace(string(decoded))
			if isV2rayFormat(&decodedStr) {
				*content = decoded
				return v2ray.Parse(content)
			}
		}
	}

	if isV2rayFormat(&contentStr) {
		return v2ray.Parse(content)
	}

	return nil, fmt.Errorf("unknown subscription format")
}

func isSingBoxFormat(content *string) bool {
	trimmed := strings.TrimSpace(*content)
	if !strings.HasPrefix(trimmed, "{") || !strings.HasSuffix(trimmed, "}") {
		return false
	}

	if strings.Contains(trimmed, `"inbounds"`) ||
		strings.Contains(trimmed, `"outbounds"`) ||
		strings.Contains(trimmed, `"route"`) ||
		strings.Contains(trimmed, `"dns"`) {
		return true
	}

	return false
}

func isMihomoFormat(content *string) bool {
	if strings.Contains(*content, "proxies:") ||
		strings.Contains(*content, "proxy-groups:") ||
		strings.Contains(*content, "rules:") ||
		strings.Contains(*content, "port:") ||
		strings.Contains(*content, "socks-port:") ||
		strings.Contains(*content, "mixed-port:") {
		return true
	}

	return false
}

func isBase64Encoded(content *string) bool {
	if len(*content) == 0 {
		return false
	}

	cleanContent := cleanBase64String(*content)

	for _, c := range cleanContent {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '+' || c == '/' || c == '=') {
			return false
		}
	}

	if len(cleanContent)%4 != 0 {
		return false
	}

	_, err := base64.StdEncoding.DecodeString(cleanContent)
	return err == nil
}

func cleanBase64String(s string) string {
	var builder strings.Builder
	builder.Grow(len(s))

	for _, c := range s {
		if c != '\n' && c != '\r' && c != ' ' && c != '\t' {
			builder.WriteRune(c)
		}
	}
	return builder.String()
}

func isV2rayFormat(content *string) bool {
	trimmed := strings.TrimSpace(*content)
	if len(trimmed) == 0 {
		return false
	}

	protocols := [...]string{
		"vmess://",
		"vless://",
		"trojan://",
		"ss://",
		"ssr://",
		"hysteria://",
		"hysteria2://",
		"tuic://",
		"http://",
		"https://",
		"socks://",
		"socks5://",
	}

	validLines := 0

	lines := strings.Split(trimmed, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		for _, protocol := range protocols {
			if strings.HasPrefix(line, protocol) {
				validLines++
				break
			}
		}
	}

	return validLines > 0
}
