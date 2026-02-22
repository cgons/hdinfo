package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cgons/hdinfo/internal/cli"
	"github.com/cgons/hdinfo/internal/lib"
)

func main() {
	lib.InitLogger()
	cmd := cli.Register()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: Unable to run hdinfo")
		os.Exit(1)
	}
}
