package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

func main() {
	flag.Parse()
	if fileBytes, err := os.ReadFile(flag.Arg(0)); err == nil {

		env := make(map[string]string)
		for _, v := range os.Environ() {
			split := strings.Split(v, "=")

			if len(split) < 2 {
				continue
			}

			env[split[0]] = strings.Join(split[1:], "=")
		}

		parsed := regexp.MustCompile(`\${[^}]+}`).ReplaceAllFunc(fileBytes, func(bytes []byte) []byte {
			// Replace with lookup in environment file
			return []byte(env[string(bytes[2:len(bytes)-1])])
		})

		fmt.Println("Generated configfile\n", string(parsed))

		fmt.Println("With environment:")
		for k, v := range env {
			fmt.Println("\t", k, "\t", v)
		}

		// Write the generated file to the desired location
		fmt.Println("Writing configfile error", os.WriteFile(flag.Arg(1), parsed, os.ModePerm))

		done := make(chan bool, 1)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGABRT, syscall.SIGALRM, syscall.SIGBUS, syscall.SIGCHLD, syscall.SIGCONT, syscall.SIGFPE, syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT, syscall.SIGIO, syscall.SIGIOT, syscall.SIGPIPE, syscall.SIGPROF, syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGSYS, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGTSTP, syscall.SIGTTIN, syscall.SIGTTOU, syscall.SIGURG, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGVTALRM, syscall.SIGWINCH, syscall.SIGXCPU, syscall.SIGXFSZ)

		go func() {
			defer func() { done <- true }()

			if len(flag.Args()) <= 2 {
				return
			}

			rest := flag.Args()[2:]

			cmd := exec.Command(rest[0], rest[1:]...)
			if cmd == nil {
				fmt.Println("Error starting command", rest[0])
				cleanAndExit(41)
			}

			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			go func() {
				for {
					s := <-sigs

					if cmd.Process == nil {
						cleanAndExit(42)
					}

					fmt.Println("About to send singnal", s.String())
					cmd.Process.Signal(s)
				}
			}()

			fmt.Println("Start error", cmd.Start())
			fmt.Println("Wait error", cmd.Wait())

			if err, ok := err.(*exec.ExitError); ok {
				if status, ok := err.Sys().(syscall.WaitStatus); ok {
					fmt.Printf("Exit Status: %d\n", status.ExitStatus())
					cleanAndExit(status.ExitStatus())
				}
			}
		}()

		// Await finishing of command
		<-done
	} else {
		fmt.Println("Could not open basefile for parsing", err)
	}

	cleanAndExit(0)
}

func cleanAndExit(exitCode int) {
	files, err := filepath.Glob("*.pid")
	if err != nil {
		fmt.Println(err)
	} else {
		for _, f := range files {
			if err := os.Remove(f); err != nil {
				fmt.Println(err)
			}
		}
	}

	os.Exit(exitCode)
}
