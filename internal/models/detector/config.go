package detector

type DetectorConfig struct {
	Type    DetectorType `json:"type" example:"alive"`
	Url     string       `json:"url"`
	Enabled bool         `json:"enable" default:"false" example:"false"`
}
