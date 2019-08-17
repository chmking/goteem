package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chmking/horde/protobuf/public"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Request the server to start a test",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("127.0.0.1:8089", grpc.WithInsecure())
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		defer conn.Close()

		client := public.NewManagerClient(conn)
		ctx := context.Background()

		_, err = client.Start(ctx, &public.StartRequest{})
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	},
}
