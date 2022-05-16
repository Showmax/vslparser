package vslparser_test

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/Showmax/vslparser"
)

func Example() {
	cmd := exec.Command("varnishlog") // Add "-b" to only process back-end requests.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	parser := vslparser.NewEntryParser(stdout)
	for {
		entry, err := parser.Parse()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(entry)
	}
}
