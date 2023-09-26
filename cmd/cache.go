package cmd

import (
	"dtool/prober"
	"dtool/utils"

	"github.com/spf13/cobra"
)

var query_cnt int
var inputfile string
var outputfile string
var goroutine_num int
var controlled_domain string
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
		prober.CreateTask(prober.RecursiveCacheProbe, controlled_domain, inputfile, outputfile, goroutine_num)
	}
}

func init() {
	cacheCmd.Flags().StringVarP(&controlled_domain, "domain", "d", "echodns.xyz", "controlled domain")
	cacheCmd.Flags().StringVarP(&inputfile, "input", "i", "", "input file(optional)")
	cacheCmd.Flags().StringVarP(&outputfile, "output", "o", "", "output file(optional)")
	cacheCmd.MarkFlagsRequiredTogether("input", "output")
	cacheCmd.Flags().IntVarP(&query_cnt, "num", "n", 20, "number of queries in one test")
	cacheCmd.Flags().IntVarP(&goroutine_num, "concurrency", "t", 150, "number of goroutine")
	rootCmd.AddCommand(cacheCmd)
}
