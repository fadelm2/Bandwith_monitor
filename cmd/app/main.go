package main

import (
	"wan-system/internal/api"
	"wan-system/internal/consumer"
	"wan-system/internal/database"
)

func main() {
	database.Init()
	go consumer.Start()
	go api.StartPublic()
	api.StartInternal()
}
