package main

import (
	"fmt"

	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/config"
)

func main() {
	cfg := config.MustGetConfig()
	fmt.Println(cfg)
}
