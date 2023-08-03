package cmd

import (
	"dtool/prober"
	"dtool/utils"
	"errors"

	"github.com/spf13/cobra"
)

var filename string
var output_file string
var thread_num int
var upstreamCmd = &cobra.Command{
	Use: "upstream",
	//Aliases: []string{"up", "stream"},
	Short: "probe upstream recursive resolvers",
	Long: `get the upstream recursive resolvers of the input resolvers
input target can be added as an argument or as a file

	-f	input file with limited ip addresses (limit=50)
	-o  output file default json type
	-t	number of goroutine`,
	//Args: cobra.ExactArgs(1),
	Run: upstream,
}

func upstream(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		panic(errors.New("too many arguments!"))
	} else if len(args) == 1 {
		if utils.IsValidIP(args[0]) {
			prober.Get_upstream_ip(args[0])
		} else {
			panic(errors.New("invalid ip address"))
		}
	} else if len(args) == 0 {
		prober.Get_upstream_file(filename, output_file, thread_num)
	}
}

func init() {
	upstreamCmd.Flags().StringVarP(&filename, "file", "f", "", "input file(optional)")
	upstreamCmd.Flags().StringVarP(&output_file, "output", "o", "", "output file(optional)")
	upstreamCmd.MarkFlagsRequiredTogether("file", "output")
	upstreamCmd.Flags().IntVarP(&thread_num, "threads", "t", 10, "number of concurrent threads")
	rootCmd.AddCommand(upstreamCmd)
}
