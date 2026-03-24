package input

type StatusInput struct {
	HookEventName  string         `json:"hook_event_name"`
	SessionID      string         `json:"session_id"`
	TranscriptPath string         `json:"transcript_path"`
	CWD            string         `json:"cwd"`
	Model          *ModelInfo     `json:"model"`
	Workspace      *WorkspaceInfo `json:"workspace"`
	Version        string         `json:"version"`
	OutputStyle    *OutputStyle   `json:"output_style"`
	Cost           *CostInfo      `json:"cost"`
	ContextWindow  *ContextWindow `json:"context_window"`
	Exceeds200k    bool           `json:"exceeds_200k_tokens"`
	RateLimits     *RateLimits    `json:"rate_limits"`
	Vim            *VimInfo       `json:"vim"`
}

type ModelInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type WorkspaceInfo struct {
	CurrentDir string   `json:"current_dir"`
	ProjectDir string   `json:"project_dir"`
	AddedDirs  []string `json:"added_dirs"`
}

type OutputStyle struct {
	Name string `json:"name"`
}

type CostInfo struct {
	TotalCostUSD       float64 `json:"total_cost_usd"`
	TotalDurationMs    int64   `json:"total_duration_ms"`
	TotalAPIDurationMs int64   `json:"total_api_duration_ms"`
	TotalLinesAdded    int     `json:"total_lines_added"`
	TotalLinesRemoved  int     `json:"total_lines_removed"`
}

type ContextWindow struct {
	TotalInputTokens  int           `json:"total_input_tokens"`
	TotalOutputTokens int           `json:"total_output_tokens"`
	ContextWindowSize int           `json:"context_window_size"`
	UsedPercentage    float64       `json:"used_percentage"`
	RemainingPct      float64       `json:"remaining_percentage"`
	CurrentUsage      *CurrentUsage `json:"current_usage"`
}

type CurrentUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

type RateLimits struct {
	FiveHour *RateLimitPeriod `json:"five_hour"`
	SevenDay *RateLimitPeriod `json:"seven_day"`
}

type RateLimitPeriod struct {
	UsedPercentage *float64 `json:"used_percentage"`
	ResetsAt       *int64   `json:"resets_at"`
}

type VimInfo struct {
	Mode string `json:"mode"`
}
