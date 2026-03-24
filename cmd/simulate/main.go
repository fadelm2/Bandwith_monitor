package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"wan-system/internal/config"

	"github.com/nats-io/nats.go"
)

type Payload struct {
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

func main() {
	log := config.NewLogger()
	viper := config.NewViper()

	url := viper.GetString("nats.url")
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Infof("Simulation started. Targeting 99%% for first 6 sites... Publishing to %s...", url)

	// Counters for each WAN
	counters := make(map[string]*struct{ rx, tx int64 })
	for i := 1; i <= 40; i++ {
		wanID := fmt.Sprintf("WAN-%03d", i)
		counters[wanID] = &struct{ rx, tx int64 }{rx: 1000000, tx: 500000}
	}

	for {
		for i := 1; i <= 40; i++ {
			wanID := fmt.Sprintf("WAN-%03d", i)
			c := counters[wanID]

			var rxDelta, txDelta int64

			if i <= 6 {
				// Target 99% of capacity on RX directly
				// Capacity: WAN-00i = (100 + i*10) Mbps
				capacityMbps := float64(100 + i*10)
				targetRxMbps := 0.99 * capacityMbps
				// bytes = Mbps * 1024 * 1024 / 8
				rxDelta = int64(targetRxMbps * 1024 * 1024 / 8)
				txDelta = int64(targetRxMbps * 0.3 * 1024 * 1024 / 8)
			} else {
				// Normal random traffic
				rxDelta = rand.Int63n(5_000_000)
				txDelta = rand.Int63n(2_500_000)
			}

			c.rx += rxDelta
			c.tx += txDelta

			payload := Payload{}
			payload.Tags.Hostname = "sim-host"
			payload.Tags.IfName = "eth0"
			payload.Tags.IfAlias = wanID
			payload.Fields.RxBytes = c.rx
			payload.Fields.TxBytes = c.tx

			data, _ := json.Marshal(payload)
			nc.Publish("network.wan", data)
		}

		time.Sleep(1 * time.Second)
	}
}
