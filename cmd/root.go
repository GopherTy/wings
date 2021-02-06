package cmd

import (
	"log"

	"github.com/gopherty/blog/common"
	"github.com/gopherty/blog/common/conf"
	"github.com/gopherty/blog/common/db"
	"github.com/gopherty/blog/common/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wings",
	Short: "Wings is a blog website",
	Long: `Wings symbolizes freedom. 
It is a personal dynamic blog site built by golang and angular.`,
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	registers := []common.IRegister{
		conf.Register{},
		logger.Register{},
		db.Register{},
	}

	for _, reg := range registers {
		if err := reg.Regist(); err != nil {
			log.Fatalf("%s register failed. %s %v %s\n", reg.Name(), red, err, reset)
		}
	}
}

// Execute .
func Execute() error {
	return rootCmd.Execute()
}
