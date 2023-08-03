package cmd

import (
	"github.com/spf13/cobra"
)

var input string
var output string
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "query the ip and nameserver information",
	Long:  "query the ip and nameserver information",
	Run:   getDomainInfo,
}

func getDomainInfo(cmd *cobra.Command, args []string) {

}

func init() {
	domainCmd.Flags().StringVarP(&input, "input", "i", "", "")
	domainCmd.Flags().StringVarP(&output, "output", "o", "", "")
	domainCmd.MarkFlagsRequiredTogether("input", "output")
	rootCmd.AddCommand(domainCmd)
}
