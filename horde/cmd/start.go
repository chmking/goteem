package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/chmking/horde/protobuf/public"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start COUNT RATE",
	Short: "Request the server to start a test",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		count, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			log.Fatal("Count is invalid")
		}

		rate, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			log.Fatal("Rate is invalid")
		}

		conn, err := grpc.Dial("127.0.0.1:8089", grpc.WithInsecure())
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		defer conn.Close()

		client := public.NewManagerClient(conn)
		ctx := context.Background()

		_, err = client.Start(ctx, &public.StartRequest{
			Users: int32(count),
			Rate:  rate,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
