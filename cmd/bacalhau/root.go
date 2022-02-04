package bacalhau

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonrpcPort int
var developmentMode bool

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(submitCmd)
	rootCmd.PersistentFlags().IntVar(
		&jsonrpcPort, "jsonrpc-port", 1234,
		`The port for the client and server to communicate on over localhost (via jsonrpc).`,
	)
	rootCmd.PersistentFlags().BoolVar(
		&developmentMode, "dev", false,
		`Development mode makes it easier to run multiple bacalhau nodes on the same machine.`,
	)
}

var rootCmd = &cobra.Command{
	Use:   "bacalhau",
	Short: "Compute over data",
	Long:  `Compute over data`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}