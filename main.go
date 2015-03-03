package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	log.SetFlags(0)

	// Read all of STDIN in.
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// Match the first filename+line.
	var match string
	buf = regexp.MustCompile(`\w+\.go:\d+`).ReplaceAllFunc(buf, func(b []byte) []byte {
		if match == "" && !bytes.HasPrefix(b, []byte("testing.go:")) {
			match = string(b)
			b = append([]byte("\033[1m"), b...)
			b = append(b, []byte("\033[0m")...)
		}
		return b
	})

	// Write out stdin back to stdout.
	os.Stdout.Write(buf)

	// Copy match to clipboard.
	cmd := exec.Command("pbcopy")
	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else if _, err = in.Write([]byte(match)); err != nil {
		log.Fatal(err)
	} else if err = in.Close(); err != nil {
		log.Fatal(err)
	} else if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
