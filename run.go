package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func run(ns *Namespace) *exec.Cmd {

	chunks := strings.Split(ns.Run, " ")

	cmd := exec.Command(chunks[0], chunks[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()

	return cmd
}

func kill(cmd *exec.Cmd) {
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15)
	}
	cmd.Wait()

}

func Run(config *Config, namespace string) error {

	quit := make(chan os.Signal, 1)
	exit := make(chan struct{}, 1)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	ns := config.GetNameSpace(namespace)

	cmd := run(ns)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Write || event.Op == fsnotify.Create {
					fmt.Println("±", event.Name)
					kill(cmd)
					cmd = run(ns)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-quit
		kill(cmd)
		exit <- struct{}{}
	}()

	for _, path := range ns.Watch.Include {
		go walk(path, watcher.Add, ns.Watch.Exclude)
	}

	fmt.Println("grape: watching", ns.Watch.Include)
	<-exit

	return nil
}

func walk(path string, fn func(string) error, ignore []string) error {

	return nil
}
