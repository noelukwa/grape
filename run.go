package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

var (
	RunNotice = `🍇 now watching for changes ✨`

	StopNotice = `🍇 stopped watching for changes, cleaning up... ✨`
)

func run(ns *Namespace) *exec.Cmd {

	chunks := strings.Split(ns.Run, " ")

	cmd := exec.Command(chunks[0], chunks[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()
	fmt.Println(infoText(RunNotice))
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
					fmt.Println(delText(event.Name))
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

	for _, targets := range ns.Watch.Include {
		go walk(targets, watcher.Add, ns.Watch.Exclude)
	}
	<-exit

	return nil
}

func walk(watchTarget string, fn func(string) error, ignore []string) {
	pathsToWatch, err := fs.Glob(os.DirFS("."), watchTarget)
	if err != nil && err != fs.ErrNotExist {
		log.Fatal(err.Error())
	}

	for _, path := range pathsToWatch {
		if err := fn(path); err != nil {
			log.Fatal(err.Error())
		}
	}
}
