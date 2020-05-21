package main

import (
	"os"

	"github.com/hashicorp/tfproviderdocsgen/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Args[1:]))
}
