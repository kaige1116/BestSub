package task

import (
	"github.com/bestruirui/bestsub/internal/models/detector"
)

type DetectConfig struct {
	Type    detector.DetectorType `json:"type"`
	Url     string                `json:"url"`
	Enabled bool                  `json:"enable" default:"false" example:"false"`
}
