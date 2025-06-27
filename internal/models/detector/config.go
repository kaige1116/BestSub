package detector

type DetectorConfig struct {
	Type     DetectorType `json:"type" example:"alive"`
	Url      *string      `json:"url"`
	CronExpr string       `json:"cron_expr" default:"*/30 * * * *" example:"*/30 * * * *"`
	Enabled  *bool        `json:"enabled" default:"false" example:"false"`
}
