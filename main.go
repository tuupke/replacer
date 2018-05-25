package main

import (
	"flag"
	"os"
	"io/ioutil"
	"os/exec"
	"os/signal"
	"syscall"
	"fmt"
	"strings"
	"regexp"
)

func main() {
	flag.Parse()
	if fileBytes, err := ioutil.ReadFile(flag.Arg(0)); err == nil {

		env := make(map[string]string)
		for _, v := range os.Environ() {
			split := strings.Split(v, "=")

			if len(split) < 2 {
				continue
			}

			env[split[0]] = strings.Join(split[1:], "=")
		}

		parsed := regexp.MustCompile("\\${[^}]+}").ReplaceAllFunc(fileBytes, func(bytes []byte) []byte {
			// Replace with lookup in environment file
			return []byte(env[string(bytes[2:len(bytes) - 1])])
		})

		fmt.Println("Generated configfile\n", string(parsed))

		fmt.Println("With environment:")
		for k, v := range env {
			fmt.Println("\t", k, "\t", v)
		}

		// Write the generated file to the desired location
		fmt.Println("Writing configfile error", ioutil.WriteFile(flag.Arg(1), []byte(parsed), os.ModePerm))

		done := make(chan bool, 1)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGABRT, syscall.SIGALRM, syscall.SIGBUS, syscall.SIGCHLD, syscall.SIGCONT, syscall.SIGFPE, syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT, syscall.SIGIO, syscall.SIGIOT, syscall.SIGKILL, syscall.SIGPIPE, syscall.SIGPROF, syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGSTOP, syscall.SIGSYS, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGTSTP, syscall.SIGTTIN, syscall.SIGTTOU, syscall.SIGURG, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGVTALRM, syscall.SIGWINCH, syscall.SIGXCPU, syscall.SIGXFSZ)

		go func() {
			defer func() {done <- true}()

			if len(flag.Args()) <= 2 {
				return
			}

			rest := flag.Args()[2:]

			cmd := exec.Command(rest[0], rest[1:]...)
			if cmd == nil {
				fmt.Println("Error starting command", rest[0])
				os.Exit(41)
			}

			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			go func() {
				for true {
					s := <- sigs

					if cmd.Process == nil {
						os.Exit(42)
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
					os.Exit(status.ExitStatus())
				}
			}
		}()

		// Await finishing of command
		<- done
	} else {
		fmt.Println("Could not open basefile for parsing", err)
	}
}
