package setting

import "github.com/bestruirui/bestsub/internal/utils/desc"

type SettingAdvance = desc.Data

type GroupSettingAdvance struct {
	GroupName   string           `json:"group_name"`
	Description string           `json:"description"`
	Data        []SettingAdvance `json:"data"`
}

type Setting struct {
	Key   string `json:"key" example:"proxy.enabled"`
	Value string `json:"value" example:"true"`
}

const (
	PROXY_ENABLE = "proxy_enable"
	PROXY_URL    = "proxy_url"

	LOG_RETENTION_DAYS = "log_retention_days"

	FRONTEND_URL           = "frontend_url"
	FRONTEND_URL_PROXY     = "frontend_url_proxy"
	SUBCONVERTER_URL       = "subconverter_url"
	SUBCONVERTER_URL_PROXY = "subconverter_url_proxy"

	SUB_DISABLE_AUTO = "sub_disable_auto"

	NODE_POOL_SIZE    = "node_pool_size"
	NODE_TEST_URL     = "node_test_url"
	NODE_TEST_TIMEOUT = "node_test_timeout"

	TASK_MAX_THREAD  = "task_max_thread"
	TASK_MAX_TIMEOUT = "task_max_timeout"
	TASK_MAX_RETRY   = "task_max_retry"

	NOTIFY_OPERATION = "notify_operation"
	NOTIFY_ID        = "notify_id"
)
