package main

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/compozed/travel-agent/models"
)

func main() {
	var buf bytes.Buffer

	config, err := LoadFromFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nCould not parse travel-agent config file.")
		panic(err)
	}

	err = ManifestTmpl(&buf, config)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}
