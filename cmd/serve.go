package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/gopherty/wings/common"
	"github.com/gopherty/wings/common/conf"
	"github.com/gopherty/wings/common/db"
	"github.com/gopherty/wings/common/logger"
	"github.com/gopherty/wings/module"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
		// init common module
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

		// run grpc serve
		cnf := conf.Instance()
		l, err := net.Listen("tcp", cnf.Server.Address)
		if err != nil {
			logger.Instance().Sugar().Fatalf("server listen failed. %s %v %s \n", red, err, reset)
		}
		srv := grpc.NewServer()
		if err = module.InitManager(srv); err != nil {
			logger.Instance().Sugar().Fatalf("init grpc module manager failed. %s %v  %s\n", red, err, reset)
		}
		if err := srv.Serve(l); err != nil {
			logger.Instance().Sugar().Fatalf("server serve failed. %s %v %s\n", red, err, reset)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
