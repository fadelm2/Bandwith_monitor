package usecase

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"wan-system/internal/entity"
	"wan-system/internal/model"
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
	agents, err := c.TelegrafRepository.List(c.DB)
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

	// Wait, we should sync to file automatically
	go func() {
		err := c.SyncConfigToFile(context.Background())
		if err != nil {
			c.Log.Errorf("Failed to automatically sync Telegraf config: %v", err)
		}
	}()

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

func (c *TelegrafUseCase) SyncConfigToFile(ctx context.Context) error {
	config, err := c.GenerateSnmpConfig(ctx)
	if err != nil {
		return err
	}

	// Path to telegraf.conf - typically /etc/telegraf/telegraf.conf
	// FOR SAFETY: You might want to use a specific include directory like /etc/telegraf/telegraf.d/snmp.conf
	configPath := "/etc/telegraf/telegraf.conf"

	// If running on Windows (Local Development), write to a local temp file instead
	if os.PathSeparator == '\\' {
		configPath = "telegraf_test.conf"
	}

	err = os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		c.Log.Errorf("Failed to write telegraf config: %v", err)
		return err
	}

	c.Log.Infof("Telegraf config updated successfully at %s", configPath)

	// Reload Telegraf
	// Note: Application needs sudo permissions without password for this to work
	cmd := exec.Command("sudo", "systemctl", "reload", "telegraf")
	if os.PathSeparator == '\\' {
		return nil // Skip reload on windows
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.Log.Errorf("Failed to reload telegraf: %v, output: %s", err, string(output))
		return err
	}

	c.Log.Info("Telegraf service reloaded successfully")
	return nil
}
