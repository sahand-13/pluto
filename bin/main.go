package main

import (
	"fmt"
	"os"
	"pluto"
	_ "pluto/panel"
)

func init() {
	pluto.ApplicationLogger.Channel.Join(pluto.BaseJoinable{
		Identifier: pluto.ExternalIdentifier{
			Name: "STD_OUT",
			Kind: "LocalStream",
		},
		Processor: pluto.IOWriter{Writer: os.Stdout},
	})
}

func main() {
	fmt.Printf("%s %s\n", pluto.Name, pluto.Version)
	select {}
}
