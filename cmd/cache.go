package cmd

import (
	"dtool/prober"
	"dtool/scheduler"
	"dtool/utils"

	"github.com/spf13/cobra"
)

var query_cnt int
var inputfile string
var outputfile string
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "cache related test",
	Long:  "cache related test",
	Run:   cache_test,
}

func cache_test(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		if utils.IsValidIP(args[0]) {
			prober.RecursiveCacheTest(args[0], query_cnt)
		}
	} else {
		scheduler.CreateTask(prober.RecursiveCacheProbe, inputfile, outputfile, 10)
	}
}

func init() {
	cacheCmd.Flags().StringVarP(&inputfile, "input", "i", "", "input file(optional)")
	cacheCmd.Flags().StringVarP(&outputfile, "output", "o", "", "output file(optional)")
	cacheCmd.MarkFlagsRequiredTogether("input", "output")
	cacheCmd.Flags().IntVarP(&query_cnt, "num", "n", 20, "number of queries in one test")
	rootCmd.AddCommand(cacheCmd)
}
