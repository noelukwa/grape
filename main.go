package main

import (
	"os"

	"github.com/noelukwa/grape/cmd"
)

func main() {

	if err := cmd.RootCmd().Execute(); err != nil {
		os.Exit(1)
	}

}
