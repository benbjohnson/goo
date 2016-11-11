package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/atotto/clipboard"
)

// re matches FILE:LINE.
var re = regexp.MustCompile(`(.+?)(\S+\.go):(\d+)(.+)`)

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

	// Print a visual beep to cause the dock icon to bounce.
	fmt.Print("\x07")

	if err := cmd.Wait(); err != nil {
		os.Exit(1)
	}
}

// processPipe scans the src by line and attempts to match the first FILE:LINE.
func processPipe(dst io.Writer, src io.Reader) {
	// Find absolute path of present wording directory.
	pwd, _ := os.Getwd()
	pwd, _ = filepath.Abs(pwd)

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		func() {
			mu.Lock()
			defer mu.Unlock()

			// Ignore if already matched a line.
			if matched {
				fmt.Fprintln(dst, line)
				return
			}

			// Find .go path.
			m := re.FindStringSubmatch(line)
			if len(m) == 0 {
				fmt.Fprintln(dst, line)
				return
			}
			prefix, path, num, suffix := m[1], m[2], m[3], m[4]

			// Determine absolute path of Go file.
			path, _ = filepath.Abs(path)

			// Ignore if path is not relative to pwd.
			rel, err := filepath.Rel(pwd, path)
			if err != nil || strings.HasPrefix(rel, "..") {
				fmt.Fprintln(dst, line)
				return
			}

			// Copy match.
			clipboard.WriteAll(rel + ":" + num)

			// Bold line.
			line = prefix + pwd + "/\033[7m" + rel + ":" + num + "\033[0m" + suffix
			fmt.Fprintln(dst, line)

			// Ensure no more lines match.
			matched = true
		}()
	}
}
