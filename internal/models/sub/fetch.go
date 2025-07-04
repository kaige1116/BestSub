package sub

type FetchResult struct {
	Type      string `json:"sub_type"`
	NodeCount int    `json:"node_count"`
	Size      int64  `json:"size"`
	Duration  string `json:"duration"`
	Error     string `json:"error"`
}
