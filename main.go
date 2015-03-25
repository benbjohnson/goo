package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"sync"

	"github.com/atotto/clipboard"
)

// re matches FILE:LINE.
var re = regexp.MustCompile(`\S+\.go:\d+`)

var matched bool
var mu sync.Mutex

func main() {
	log.SetFlags(0)

	// Execute "go" command with the same arguments.
	cmd := exec.Command("go", os.Args[1:]...)

	// Pass through standard input.
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(stdin, os.Stdin)

	// Create a wait group for stdout/stderr.
	var wg sync.WaitGroup
	wg.Add(2)

	// Pass through standard out.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		processPipe(os.Stdout, stdout)
		wg.Done()
	}()

	// Read through stderr and decorate.
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		processPipe(os.Stderr, stderr)
		wg.Done()
	}()

	// Execute command.
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// Forward signals to command.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c)
		for sig := range c {
			cmd.Process.Signal(sig)
		}
	}()

	// Wait for pipes to finish reading and then wait for command to exit.
	wg.Wait()
	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

// processPipe scans the src by line and attempts to match the first FILE:LINE.
func processPipe(dst io.Writer, src io.Reader) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		func() {
			mu.Lock()
			defer mu.Unlock()

			if !matched {
				if m := re.FindString(line); m != "" && !strings.Contains(m, "testing.go") {
					// Remove "./" prefix.
					m = strings.TrimPrefix(m, "./")

					// Remove present working directory prefix.
					if pwd, _ := os.Getwd(); pwd != "" {
						m = strings.TrimPrefix(m, pwd+"/")
					}

					// Copy match.
					clipboard.WriteAll(m)

					// Bold line.
					line = "\033[1m" + line + "\033[0m"
					matched = true
				}
			}
			fmt.Fprintln(dst, line)
		}()
	}
}
