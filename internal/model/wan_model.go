package model

import "time"

type WanCapacityRequest struct {
	WanID            string  `json:"wan_id"`
	CapacityMbps     float64 `json:"capacity_mbps"`
	ThresholdPercent float64 `json:"threshold_percent"`
	Description      string  `json:"description"`
}

type WanCapacityResponse struct {
	WanID            string    `json:"wan_id"`
	CapacityMbps     float64   `json:"capacity_mbps"`
	ThresholdPercent float64   `json:"threshold_percent"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
}

type SearchTrafficRequest struct {
	WanID        string `query:"wan_id"`
	Page         int    `query:"page"`
	Size         int    `query:"size"`
	SinceMinutes int    `query:"since_minutes"` // filter to last N minutes
}

type WanTrafficResponse struct {
	ID                 uint      `json:"id"`
	WanID              string    `json:"wan_id"`
	Hostname           string    `json:"hostname"`
	AgentHost          string    `json:"agent_host"`
	InterfaceName      string    `json:"interface_name"`
	RxMbps             float64   `json:"rx_mbps"`
	TxMbps             float64   `json:"tx_mbps"`
	CapacityMbps       float64   `json:"capacity_mbps"`
	UtilizationPercent float64   `json:"utilization_percent"`
	CreatedAt          time.Time `json:"created_at"`
}

type WanAlertResponse struct {
	WanID            string  `json:"wan_id"`
	Description      string  `json:"description"`
	AvgRxMbps        float64 `json:"avg_rx_mbps"`
	AvgTxMbps        float64 `json:"avg_tx_mbps"`
	CurrentCapacity  float64 `json:"current_capacity"`
	CurrentThreshold float64 `json:"current_threshold"`
	AvgUtilization   float64 `json:"avg_utilization"`
}
