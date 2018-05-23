package main

import (
	"flag"
	"os"
	"io/ioutil"
	"os/exec"
	"os/signal"
	"syscall"
	"fmt"
)

func main() {
	flag.Parse()
	if fileBytes, err := ioutil.ReadFile(flag.Arg(0)); err == nil {
		parsed := os.ExpandEnv(string(fileBytes))
		fmt.Println("Generated configfile\n", parsed)
		fmt.Println("With environment:")
		for _, v := range os.Environ() {
			fmt.Println("\t", v)
		}

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
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			go func() {
				for true {
					s := <- sigs
					cmd.Process.Signal(s)
				}
			}()

			fmt.Println("Start error", cmd.Start())
			fmt.Println("Wait error", cmd.Wait())
		}()

		// Await finishing of command
		<- done
	} else {
		fmt.Println("Could not open basefile for parsing", err)
	}
}
