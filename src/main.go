package main

import (
	"log"
	"github.com/tomhjx/network-probe/service"
)

func main() {
	if err := service.NewProcessor().Run();err != nil {
		log.Fatal(err)
	}
}