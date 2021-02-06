package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var (
	version string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wings",
	Long:  `All software has versions. This is wings`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s go version %s %s/%s %s\n", magenta, runtime.Version(), runtime.GOOS, runtime.GOARCH, reset)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
