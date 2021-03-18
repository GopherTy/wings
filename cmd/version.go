package cmd

import (
	"fmt"
	"runtime"

	"github.com/gopherty/wings/pkg/colors"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wings",
	Long:  `All software has versions. This is wings`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s go version %s %s/%s %s\n", colors.Magenta, runtime.Version(), runtime.GOOS, runtime.GOARCH, colors.Reset)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
