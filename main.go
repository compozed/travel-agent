package main

import (
	"bytes"
	"fmt"
	. "github.com/compozed/travel-agent/models"
	"os"
)

func main() {
	var buf bytes.Buffer

	envs, err := LoadFromFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = ManifestTmpl(&buf, envs)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}
