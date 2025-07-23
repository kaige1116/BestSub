package task

import (
	"time"
)

type Data struct {
	ID     uint16 `db:"id"`
	Name   string `db:"name"`
	Enable bool   `db:"enable"`
	System bool   `db:"system"`
	Config string `db:"config"` // 以json格式存储
	Extra  string `db:"extra"`  // 以json格式存储
	Result string `db:"result"` // 以json格式存储
}

type Config struct {
	Cron          string `json:"cron"`
	Type          string `json:"type"`
	LogLevel      string `json:"log_level"`
	Timeout       int    `json:"timeout"`
	Notify        bool   `json:"notify"`
	NotifyChannel string `json:"notify_channel"`
}

type Result struct {
	Success         int       `json:"success"`
	Failed          int       `json:"failed"`
	LastRunResult   string    `json:"last_run_result"`
	LastRunTime     time.Time `json:"last_run_time"`
	LastRunDuration int       `json:"last_run_duration"`
}

type Response struct {
	ID     uint16 `json:"id"`
	Enable bool   `json:"enable"`
	Config Config `json:"config"`
	Extra  string `json:"extra"`
	Status string `json:"status"`
	Result Result `json:"result"`
}

type CreateRequest struct {
	Config
	Extra string
}
type UpdateRequest struct {
	ID uint16 `json:"id"`
	Config
	Extra string
}
