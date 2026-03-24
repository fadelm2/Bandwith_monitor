package consumer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"wan-system/internal/database"
	"wan-system/internal/models"
)

func Start() {

	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}

	lastData := make(map[string]struct {
		Rx   int64
		Tx   int64
		Time time.Time
	})

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
			return
		}

		wanID := payload.Tags.IfAlias
		if wanID == "" {
			return
		}

		now := time.Now()
		last, ok := lastData[wanID]
		if !ok {
			lastData[wanID] = struct {
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

		var cap models.WanCapacity
		database.DB.First(&cap, "wan_id = ?", wanID)

		util := 0.0
		if cap.CapacityMbps > 0 {
			util = (rxMbps / cap.CapacityMbps) * 100
		}

		database.DB.Create(&models.WanTraffic{
			WanID:              wanID,
			Hostname:           payload.Tags.Hostname,
			InterfaceName:      payload.Tags.IfName,
			RxMbps:             rxMbps,
			TxMbps:             txMbps,
			CapacityMbps:       cap.CapacityMbps,
			UtilizationPercent: util,
			CreatedAt:          now,
		})

		lastData[wanID] = struct {
			Rx   int64
			Tx   int64
			Time time.Time
		}{payload.Fields.RxBytes, payload.Fields.TxBytes, now}
	})
}
