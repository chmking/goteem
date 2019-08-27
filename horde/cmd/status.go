package cmd

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/chmking/horde/protobuf/public"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var tmpl = `
Manager:
  State: {{ .State }}
{{- if .Orders }}
  Orders:
    ID:    {{ .Orders.Id }}
    Count: {{ .Orders.Count }}
    Rate:  {{ .Orders.Rate }}
{{ end }}
`

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

		output, err := template.New("output").Parse(tmpl)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		if err := output.Execute(os.Stdout, resp); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	},
}
