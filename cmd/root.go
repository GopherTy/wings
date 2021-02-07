package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wings",
	Short: "Wings is a blog website",
	Long: `Wings symbolizes freedom. 
It is a personal dynamic blog site built by golang and angular.`,
}

// Execute .
func Execute() error {
	return rootCmd.Execute()
}
