package cmd

import (
	"dtool/prober"
	"dtool/utils"

	"fmt"

	"github.com/spf13/cobra"
)

var record_input string
var record_output string
var record_type string
var record_domain string
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "get specific record response",
	Long:  "get specific record response",
	Run:   record_probe,
}

func record_probe(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		if utils.IsValidIP(args[0]) {
			result, err := prober.SVCBProbeOnce(args[0], record_domain)
			if err == nil {
				if output_str, err := prober.OutputHandler(result); err == nil {
					fmt.Println(output_str)
				}
			} else {
				fmt.Println(err)
			}
		}
	} else {
		prober.CreateTask(prober.SVCBProbe, record_domain, record_input, record_output, 500)
	}
}

func init() {
	recordCmd.Flags().StringVarP(&record_input, "input", "i", "", "input file(optional)")
	recordCmd.Flags().StringVarP(&record_output, "output", "o", "", "output file(optional)")
	recordCmd.MarkFlagsRequiredTogether("input", "output")
	recordCmd.Flags().StringVarP(&record_type, "type", "t", "A", "request record type")
	recordCmd.Flags().StringVarP(&record_domain, "domain", "d", "example.com", "requested domain")
	rootCmd.AddCommand(recordCmd)
}
