package usecase

import (
	"context"
	"fmt"
	"strings"
	"wan-system/internal/entity"
	"wan-system/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TelegrafUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	TelegrafRepository *repository.TelegrafRepository
}

func NewTelegrafUseCase(db *gorm.DB, logger *logrus.Logger, repo *repository.TelegrafRepository) *TelegrafUseCase {
	return &TelegrafUseCase{
		DB:                 db,
		Log:                logger,
		TelegrafRepository: repo,
	}
}

// GenerateSnmpConfig generates the [[inputs.snmp]] block for Telegraf
func (c *TelegrafUseCase) GenerateSnmpConfig(ctx context.Context) (string, error) {
	agents, err := c.TelegrafRepository.ListActive(c.DB)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("########################################\n")
	sb.WriteString("# INPUT SNMP (AUTOMATICALLY GENERATED)\n")
	sb.WriteString("########################################\n")
	sb.WriteString("[[inputs.snmp]]\n")
	sb.WriteString("  agents = [\n")

	for i, agent := range agents {
		comma := ","
		if i == len(agents)-1 {
			comma = ""
		}
		sb.WriteString(fmt.Sprintf("    \"%s://%s:%d\"%s\n", agent.Protocol, agent.IPAddress, agent.Port, comma))
	}

	sb.WriteString("  ]\n")
	
	// Add other standard SNMP settings here if needed
	sb.WriteString("  version = 2\n")
	sb.WriteString("  community = \"public\"\n")

	return sb.String(), nil
}

func (c *TelegrafUseCase) CreateAgent(ctx context.Context, req *model.TelegrafAgentRequest) (*model.TelegrafAgentResponse, error) {
	agent := &entity.TelegrafAgent{
		IPAddress:   req.IPAddress,
		Port:        req.Port,
		Protocol:    req.Protocol,
		Description: req.Description,
		IsActive:    true,
	}

	if agent.Port == 0 {
		agent.Port = 161
	}
	if agent.Protocol == "" {
		agent.Protocol = "udp"
	}

	if err := c.TelegrafRepository.Create(c.DB, agent); err != nil {
		return nil, err
	}

	return &model.TelegrafAgentResponse{
		ID:          agent.ID,
		IPAddress:   agent.IPAddress,
		Port:        agent.Port,
		Protocol:    agent.Protocol,
		Description: agent.Description,
		IsActive:    agent.IsActive,
	}, nil
}

func (c *TelegrafUseCase) ListAgents(ctx context.Context) ([]model.TelegrafAgentResponse, error) {
	agents, err := c.TelegrafRepository.List(c.DB)
	if err != nil {
		return nil, err
	}

	var responses []model.TelegrafAgentResponse
	for _, a := range agents {
		responses = append(responses, model.TelegrafAgentResponse{
			ID:          a.ID,
			IPAddress:   a.IPAddress,
			Port:        a.Port,
			Protocol:    a.Protocol,
			Description: a.Description,
			IsActive:    a.IsActive,
		})
	}

	return responses, nil
}
