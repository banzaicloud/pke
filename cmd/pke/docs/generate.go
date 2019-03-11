package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/banzaicloud/pke/cmd/pke/app/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	c := cmd.NewPKECommand(os.Stdin, ioutil.Discard, "", "", "", "")
	err := doc.GenMarkdownTree(c, ".")
	if err != nil {
		log.Fatal(err)
	}
}
