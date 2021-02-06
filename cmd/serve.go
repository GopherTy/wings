package cmd

import (
	"log"
	"net"

	"github.com/gopherty/blog/module"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {

		l, err := net.Listen("tcp", ":10000")
		if err != nil {
			log.Fatalf("listen failed: %v\n", err)
		}
		srv := grpc.NewServer()
		if err = module.InitManager(srv); err != nil {
			log.Fatalf("Initialzation module manager failed: %v\n", err)
		}
		if err := srv.Serve(l); err != nil {
			log.Fatalf("Serve failed: %v\n", err)
		}
	},
}
