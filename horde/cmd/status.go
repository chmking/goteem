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
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints the current Horde status",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("127.0.0.1:8089", grpc.WithInsecure())
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		defer conn.Close()

		client := public.NewManagerClient(conn)
		ctx := context.Background()

		resp, err := client.Status(ctx, &public.StatusRequest{})
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		fmt.Println("Manager:")
		fmt.Printf("  Status: %s\n", resp.Status.String())
		fmt.Printf("  Agents: %d\n", len(resp.Agents))

		if len(resp.Agents) == 0 {
			return
		}

		fmt.Println("")
		fmt.Println("Agents:")
		for _, agent := range resp.Agents {
			fmt.Printf("  ID: %s\n")
			fmt.Printf("    Status: %s\n", agent.Status.String())
		}
	},
}
