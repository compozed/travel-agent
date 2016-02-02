package main

import (
	"bytes"
	"fmt"
	. "github.com/compozed/travel-agent/models"
	"os"
)

func main() {
	var buf bytes.Buffer

	config, err := LoadFromFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = ManifestTmpl(&buf, config)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}
