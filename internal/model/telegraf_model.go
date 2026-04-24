package model

type TelegrafAgentRequest struct {
	IPAddress   string `json:"ip_address" validate:"required"`
	Port        int    `json:"port" default:"161"`
	Protocol    string `json:"protocol" default:"udp"`
	Description string `json:"description"`
}

type TelegrafAgentResponse struct {
	ID          int64  `json:"id"`
	IPAddress   string `json:"ip_address"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}
