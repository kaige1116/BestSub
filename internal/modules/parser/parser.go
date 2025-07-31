package parser

import (
	"fmt"

	"github.com/bestruirui/bestsub/internal/models/node"
	parserModel "github.com/bestruirui/bestsub/internal/models/parser"
	"github.com/bestruirui/bestsub/internal/modules/parser/mihomo"
	"github.com/bestruirui/bestsub/internal/modules/parser/singbox"
	"github.com/bestruirui/bestsub/internal/modules/parser/v2ray"
)

func Parse(content *[]byte, subType parserModel.ParserType, sublinkID uint16) (*[]node.Data, error) {
	if content == nil || len(*content) == 0 {
		return nil, fmt.Errorf("content is empty")
	}
	switch subType {
	case parserModel.ParserTypeAuto:
		return auto(content, sublinkID)
	case parserModel.ParserTypeMihomo:
		return mihomo.Parse(content, sublinkID)
	case parserModel.ParserTypeSingbox:
		return singbox.Parse(content, sublinkID)
	case parserModel.ParserTypeV2ray:
		return v2ray.Parse(content, sublinkID)
	default:
		return nil, fmt.Errorf("unknown subscription format")
	}
}
