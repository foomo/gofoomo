package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/foomo/gofoomo/proxy"
)

var flagConfig = flag.String("config", "", "config file")

func main() {
	fmt.Println("I am a little demo program")
	flag.Parse()
	proxyServer, err := proxy.NewServerWithConfig(*flagConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(proxyServer.ListenAndServe().Error())
}
