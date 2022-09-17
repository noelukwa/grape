package cmd

import (
	"fmt"
	"log"

	"github.com/noelukwa/grape/app"
	"github.com/noelukwa/grape/config"
	"github.com/spf13/cobra"
)

const (
	// DefaultConfigPath is the default path to the config file
	DefaultConfigPath = "grape.json"
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "use [on] to configure grape on the go without a config file.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("on called")
	},
}

var grapeCmd = &cobra.Command{
	Use:  "grape",
	Long: `🍇 grape is a tiny tool for watching files and running commands when they change during development.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "use [run] to run grape with a config file and switch between namespaces.",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := config.FromJson(cmd.Flag("config").Value.String())
		if err != nil {
			log.Fatal(err.Error())
		}

		namespace := args[0]
		if err := app.Run(config, namespace); err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("run terminated")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "use [init] to create a config file in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.NewDefault(); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func RootCmd() *cobra.Command {
	grapeCmd.AddCommand(runCmd)
	grapeCmd.AddCommand(onCmd)
	grapeCmd.AddCommand(initCmd)
	onCmd.Flags().StringSliceP("ext", "e", []string{"*.go"}, "comma separated list of file extensions to watch [ default: .go ]")
	onCmd.Flags().StringSliceP("exclude", "x", []string{"vendor"}, "comma separated list of directories to exclude from watching")
	onCmd.Flags().StringP("run", "r", "", "command to run when a file is changed")
	runCmd.Flags().StringP("config", "c", DefaultConfigPath, "path to config file")
	return grapeCmd
}
