package vslparser_test

import (
	"bufio"
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

	scanner := bufio.NewScanner(stdout)
	for {
		entry, err := vslparser.Parse(scanner)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(entry)
	}
}
