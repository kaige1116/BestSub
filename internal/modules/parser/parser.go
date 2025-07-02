package parser

import (
	"fmt"

	modparser "github.com/bestruirui/bestsub/internal/models/parser"
	"github.com/bestruirui/bestsub/internal/modules/parser/mihomo"
	"github.com/bestruirui/bestsub/internal/modules/parser/singbox"
	"github.com/bestruirui/bestsub/internal/modules/parser/v2ray"
)

func Parse(content *[]byte, subType modparser.ParserType, sublinkID int64) (modparser.ParserType, int, error) {
	if content == nil || len(*content) == 0 {
		return "", 0, fmt.Errorf("content is empty")
	}
	var addedCount int
	var err error
	switch subType {
	case modparser.ParserTypeAuto:
		return auto(content, sublinkID)
	case modparser.ParserTypeMihomo:
		addedCount, err = mihomo.Parse(content, sublinkID)
	case modparser.ParserTypeSingbox:
		addedCount, err = singbox.Parse(content, sublinkID)
	case modparser.ParserTypeV2ray:
		addedCount, err = v2ray.Parse(content, sublinkID)
	default:
		return subType, 0, fmt.Errorf("unknown subscription format")
	}
	return subType, addedCount, err
}
