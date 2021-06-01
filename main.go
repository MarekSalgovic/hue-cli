package main

import (
	"github.com/MarekSalgovic/hue-go/cli"
)

func main() {
	cmd, err := cli.NewCLI()
	if err != nil {
		panic(err)
	}
	cmd.Run()
}
