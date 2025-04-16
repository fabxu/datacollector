package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "datacollector-service",
		Short: "A demo for boilerplate project",
	}
	attachVersionCommand(rootCmd)
	attachInitCommand(rootCmd)
	attachRunCommand(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
