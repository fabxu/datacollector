package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	BuildTS   = "None"
	GitHash   = "None"
	GitBranch = "None"
)

func attachVersionCommand(rootCmd *cobra.Command) {
	versionCmd := &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Git Branch:       ", GitBranch)
			fmt.Println("Git Commit:       ", GitHash)
			fmt.Println("Build Time (UTC): ", BuildTS)
		},
	}
	rootCmd.AddCommand(versionCmd)
}
