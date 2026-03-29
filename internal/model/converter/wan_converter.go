package converter

import (
	"wan-system/internal/entity"
	"wan-system/internal/model"
)

func WanCapacityToResponse(wan *entity.WanCapacity) *model.WanCapacityResponse {
	return &model.WanCapacityResponse{
		WanID:            wan.WanID,
		CapacityMbps:     wan.CapacityMbps,
		ThresholdPercent: wan.ThresholdPercent,
		Description:      wan.Description,
		CreatedAt:        wan.CreatedAt,
	}
}

func WanTrafficToResponse(traffic *entity.WanTraffic) *model.WanTrafficResponse {
	return &model.WanTrafficResponse{
		ID:                 traffic.ID,
		WanID:              traffic.WanID,
		Hostname:           traffic.Hostname,
		AgentHost:          traffic.AgentHost,
		InterfaceName:      traffic.InterfaceName,
		RxMbps:             traffic.RxMbps,
		TxMbps:             traffic.TxMbps,
		CapacityMbps:       traffic.CapacityMbps,
		UtilizationPercent: traffic.UtilizationPercent,
		CreatedAt:          traffic.CreatedAt,
	}
}
