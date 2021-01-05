package main

import (
	"fmt"
	"os"

	"github.com/navidnabavi/jsonrpc2rest/pkg/configmanager"
	"github.com/navidnabavi/jsonrpc2rest/pkg/proxy"
)

func main() {
	if len(os.Args) != 2 {
		panic("Bad Params: expected config file address as the input parameter: ./jsonrpc2rest config.json")
	}

	fmt.Println("Starting Up")
	config := configmanager.LoadConfiguration(os.Args[1])
	fmt.Printf("Listening on %s\n", config.Bind)
	proxy := proxy.NewProxy(&config)
	proxy.Serve()
}
