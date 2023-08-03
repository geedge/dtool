package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "get server version with version.bind",
	Long:  "get server version with version.bind chaos txt request",
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {

}

func init() {
	rootCmd.AddCommand(versionCmd)
}
