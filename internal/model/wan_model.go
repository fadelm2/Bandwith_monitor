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
	InterfaceName      string    `json:"interface_name"`
	RxMbps             float64   `json:"rx_mbps"`
	TxMbps             float64   `json:"tx_mbps"`
	CapacityMbps       float64   `json:"capacity_mbps"`
	UtilizationPercent float64   `json:"utilization_percent"`
	CreatedAt          time.Time `json:"created_at"`
}
