package cmd

import (
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/gopherty/blog/common/conf"
	"github.com/gopherty/blog/common/logger"
	"github.com/gopherty/blog/module"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
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
