package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

const (
	// DefaultConfigPath is the default path to the config file
	DefaultConfigPath = "grape.json"
)

var rootCmd = &cobra.Command{
	Use: "grape",
	Long: `grape is a process manager for go projects but it could be configured to work with other runtimes as needed.
		Run [ grape ] without any flags to use default config file in root directory`,
	Run: func(cmd *cobra.Command, args []string) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		// Start listening for events.
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op == fsnotify.Write {
						fmt.Printf("%s triggered change event\n", event.Name)
					}
					if event.Op == fsnotify.Create {
						fmt.Printf("%s triggered create event\n", event.Name)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()

		// Add a path.
		err = transverse("/ignore", []string{".txt"}, watcher.Add)
		if err != nil {
			log.Fatal(errors.New("failed to add path"))
		}

		// Block main goroutine forever.
		<-make(chan struct{})
	},
}

var runCmd = &cobra.Command{
	Use:  "run",
	Long: `run is a command that runs the grape process manager`,
	Run: func(cmd *cobra.Command, args []string) {

		config, err := ConfigFromJson(cmd.Flag("config").Value.String())

		if err != nil {
			log.Fatal(err)
		}

		for _, namespace := range config.Namespaces {
			fmt.Printf("namespace: %s\n", namespace.Tag)
			fmt.Printf("%s watch exclude: %s\n", namespace.Tag, namespace.Watch.Exclude)
			fmt.Printf("%s watch include: %s\n", namespace.Tag, namespace.Watch.Include)
			fmt.Printf("%s command: %s\n", namespace.Tag, namespace.Runner.Command)

		}

		// ns := args[0]
		// config, err := cmd.Flags().GetString("config")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// configPath, err := filepath.Abs(config)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Printf("config file is %s and namespace is %s and the path is %s \n ", config, ns, configPath)
	},
}

func transverse(dir string, exts []string, fn func(string) error) error {
	cwd, _ := os.Getwd()
	// fmt.Printf("the full cwd is %s\n", cwd)

	relpath, e := filepath.Abs(dir)
	if e != nil {
		fmt.Println(e.Error())
	}
	fullPath := filepath.Join(cwd, relpath)

	// fmt.Printf("the full path is %s\n", fullPath)
	return filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
		e = fn(fullPath)
		if e != nil {
			fmt.Printf("err: %s\n", e.Error())
		}
		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }

		// fmt.Println(exts)
		// fmt.Printf("walkDir : %s", dir)
		// if d.IsDir() {
		// 	return transverse(dir, exts, fn)
		// }
		for _, extension := range exts {
			if extension == filepath.Ext(path) {
				// fmt.Printf("add file %s with extension %s to be watched\n", path, extension)
				if e := fn(path); e != nil {
					fmt.Println(e)
				}
			}
		}
		return err
	})
}
func cmd() *cobra.Command {
	rootCmd.AddCommand(runCmd)
	rootCmd.PersistentFlags().StringP("config", "c", DefaultConfigPath, "path to config file")
	return rootCmd
}
