package nats

import (
	"context"
	"encoding/json"
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

	nc.Subscribe("network.wan", func(msg *nats.Msg) {
		var payload struct {
			Tags struct {
				Hostname string `json:"hostname"`
				IfName   string `json:"ifName"`
				IfAlias  string `json:"ifAlias"`
			} `json:"tags"`
			Fields struct {
				RxBytes int64 `json:"rx_bytes"`
				TxBytes int64 `json:"tx_bytes"`
			} `json:"fields"`
		}

		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			c.Log.Warnf("Failed to unmarshal NATS message: %v", err)
			return
		}

		wanID := payload.Tags.IfAlias
		if wanID == "" {
			return
		}

		now := time.Now()
		last, ok := c.LastData[wanID]
		if !ok {
			c.LastData[wanID] = struct {
				Rx   int64
				Tx   int64
				Time time.Time
			}{payload.Fields.RxBytes, payload.Fields.TxBytes, now}
			return
		}

		duration := now.Sub(last.Time).Seconds()
		if duration <= 0 {
			return
		}

		rxDelta := payload.Fields.RxBytes - last.Rx
		txDelta := payload.Fields.TxBytes - last.Tx

		rxMbps := float64(rxDelta*8) / duration / 1_000_000
		txMbps := float64(txDelta*8) / duration / 1_000_000

		traffic := &entity.WanTraffic{
			WanID:         wanID,
			Hostname:      payload.Tags.Hostname,
			InterfaceName: payload.Tags.IfName,
			RxMbps:        rxMbps,
			TxMbps:        txMbps,
			CreatedAt:     now,
		}

		if err := c.UseCase.ProcessTraffic(context.TODO(), traffic); err != nil {
			c.Log.Warnf("Failed to process traffic: %v", err)
		}

		c.LastData[wanID] = struct {
			Rx   int64
			Tx   int64
			Time time.Time
		}{payload.Fields.RxBytes, payload.Fields.TxBytes, now}
	})
}
