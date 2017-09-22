package main

import (
	"fmt"
	"os"

	"./config_manager"
	"./proxy"
)

func main() {
	if len(os.Args) != 2 {
		panic("Bad Params: expected config file address as the input parameter: ./jsonrpc2rest config.json")
	}

	fmt.Println("Starting Up")
	config := ConfigManager.LoadConfiguration(os.Args[1])
	fmt.Printf("Listening on %s\n", config.Bind)
	proxy := Proxy.NewProxy(&config)
	proxy.Serve()
}
