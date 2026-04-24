package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"wan-system/internal/entity"
	"wan-system/internal/usecase"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type WanConsumer struct {
	Log      *logrus.Logger
	Config   *viper.Viper
	UseCase  *usecase.WanUseCase
	LastData map[string]struct {
		Rx   int64
		Tx   int64
		Time time.Time
	}
}

type WanPayload struct {
	Tags struct {
		Hostname  string `json:"hostname"`
		AgentHost string `json:"agent_host"`
		IfName    string `json:"ifName"`
		IfAlias   string `json:"ifAlias"`
	} `json:"tags"`
	Fields struct {
		RxBytes int64 `json:"rx_bytes"`
		TxBytes int64 `json:"tx_bytes"`
	} `json:"fields"`
}

func NewWanConsumer(logger *logrus.Logger, config *viper.Viper, useCase *usecase.WanUseCase) *WanConsumer {
	return &WanConsumer{
		Log:     logger,
		Config:  config,
		UseCase: useCase,
		LastData: make(map[string]struct {
			Rx   int64
			Tx   int64
			Time time.Time
		}),
	}
}

func (c *WanConsumer) Start() {
	url := c.Config.GetString("nats.url")
	nc, err := nats.Connect(url)
	if err != nil {
		c.Log.Fatalf("Failed to connect to NATS: %v", err)
	}

	c.Log.Infof("NATS Consumer started on topic 'network.wan' at %s", url)

	nc.Subscribe("network.wan", func(msg *nats.Msg) {
		data := string(msg.Data)
		if len(data) == 0 {
			return
		}

		var payload WanPayload
		var err error

		if data[0] == '{' {
			// Handle JSON format
			err = json.Unmarshal(msg.Data, &payload)
		} else {
			// Handle Influx Line Protocol format
			payload, err = parseLineProtocol(data)
		}

		if err != nil {
			c.Log.Debugf("Failed to parse NATS message: %v. Message data: %s", err, data)
			return
		}

		// Use ifAlias as WanID. If empty, skip.
		wanID := payload.Tags.IfAlias
		if wanID == "" {
			return
		}

		now := time.Now()
		last, ok := c.LastData[wanID]
		if !ok {
			// First data point for this WAN ID
			c.LastData[wanID] = struct {
				Rx   int64
				Tx   int64
				Time time.Time
			}{payload.Fields.RxBytes, payload.Fields.TxBytes, now}
			c.Log.Infof("Discovered new WAN ID from Telegraf: %s (Host: %s, Interface: %s)", wanID, payload.Tags.Hostname, payload.Tags.IfName)
			return
		}

		duration := now.Sub(last.Time).Seconds()
		if duration <= 0 {
			return
		}

		// Calculate Delta (Byte counts are cumulative)
		rxDelta := payload.Fields.RxBytes - last.Rx
		txDelta := payload.Fields.TxBytes - last.Tx

		// Basic sanity check to handle counter resets
		if rxDelta < 0 || txDelta < 0 {
			c.Log.Infof("Counter reset detected for %s. Resetting baseline.", wanID)
			c.LastData[wanID] = struct {
				Rx   int64
				Tx   int64
				Time time.Time
			}{payload.Fields.RxBytes, payload.Fields.TxBytes, now}
			return
		}

		// Convert Bytes to Mbps (Bytes * 8 / Duration / 1,000,000)
		rxMbps := float64(rxDelta*8) / duration / 1_000_000
		txMbps := float64(txDelta*8) / duration / 1_000_000

		traffic := &entity.WanTraffic{
			WanID:         wanID,
			Hostname:      payload.Tags.Hostname,
			AgentHost:     payload.Tags.AgentHost,
			InterfaceName: payload.Tags.IfName,
			RxMbps:        rxMbps,
			TxMbps:        txMbps,
			CreatedAt:     now,
		}

		if err := c.UseCase.ProcessTraffic(context.TODO(), traffic); err != nil {
			c.Log.Warnf("Failed to process traffic for %s: %v", wanID, err)
		}

		c.LastData[wanID] = struct {
			Rx   int64
			Tx   int64
			Time time.Time
		}{payload.Fields.RxBytes, payload.Fields.TxBytes, now}
	})
}

// Simple parser for Influx Line Protocol specialized for Telegraf's WAN output
func parseLineProtocol(line string) (WanPayload, error) {
	var p WanPayload
	line = strings.TrimSpace(line)

	// Format: measurement,tags fields timestamp
	// Example: Harapan-Mulya,agent_host=...,hostname=...,ifAlias=---WAN-123-TIS\ FO---,ifName=ether1 rx_bytes=39064919446004i,tx_bytes=3566278155161i 1774595030000000000

	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return p, errors.New("invalid line protocol format")
	}

	// 1. Tags part
	tagSection := parts[0]
	tagParts := strings.Split(tagSection, ",")
	for _, tp := range tagParts {
		kv := strings.Split(tp, "=")
		if len(kv) == 2 {
			k := kv[0]
			v := strings.ReplaceAll(kv[1], "\\ ", " ") // Unescape spaces
			switch k {
			case "hostname":
				p.Tags.Hostname = v
			case "agent_host":
				p.Tags.AgentHost = v
			case "ifName":
				p.Tags.IfName = v
			case "ifAlias":
				p.Tags.IfAlias = v
			}
		}
	}

	// 2. Fields part
	fieldSection := parts[1]
	fieldParts := strings.Split(fieldSection, ",")
	for _, fp := range fieldParts {
		kv := strings.Split(fp, "=")
		if len(kv) == 2 {
			k := kv[0]
			v := strings.TrimSuffix(kv[1], "i") // Remove Influx 'i' suffix for integers
			var val int64
			fmt.Sscanf(v, "%d", &val)
			switch k {
			case "rx_bytes":
				p.Fields.RxBytes = val
			case "tx_bytes":
				p.Fields.TxBytes = val
			}
		}
	}

	return p, nil
}
