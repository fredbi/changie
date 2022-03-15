package main

import (
	"fmt"
	"os"

	"github.com/fredbi/changie/cmd"
)

// goreleaser injected values
var version = "dev"

func main() {
	if err := cmd.Execute("v" + version); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
