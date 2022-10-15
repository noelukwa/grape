package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/noelukwa/grape/config"
)

func run(ns *config.Namespace) *exec.Cmd {

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

func Run(config *config.Config, namespace string) error {

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

	err = transverse(ns, watcher.Add)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("grape: watching", ns.Watch.Include)
	<-exit

	return nil
}

func transverse(ns *config.Namespace, fn func(string) error) error {

	paths_to_watch := make(chan []string, 1)

	go func(w *config.FWatcher) {
		for _, path := range w.Include {
			matches, err := filepath.Glob(path)
			if err != nil {
				log.Fatal(err)
			}
			paths_to_watch <- matches
		}
		close(paths_to_watch)
	}(&ns.Watch)

	for paths := range paths_to_watch {
		for _, path := range paths {
			if err := fn(path); err != nil {
				return err
			}
		}
	}

	return nil
}
