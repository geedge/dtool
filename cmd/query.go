package cmd

import (
	"dtool/utils"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// input parameters
// var query_target string
var query_port string
var query_name string
var query_type uint16
var query_class uint16
var query_rate int
var query_packets int
var query_random_prefix bool

/*
DNS请求测试功能，可设置参数包括域名（随机前缀），资源记录type/class，请求速率，请求数量
*/
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "",
	Long:  "",
	Run:   queryTest,
}

func concurrent_request(ip, domain, port string, qtype, qclass uint16, result chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	res, err := utils.SendQuery(ip, domain, port, qtype, qclass)
	if err != nil {
		result <- (-1)
		//fmt.Println(err)
	} else {
		rcode := utils.ParseRcode(res)
		result <- rcode
	}
	wg.Done()
}

func response_stat(stat []int, result chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		if v, ok := <-result; ok {
			switch v {
			case 0:
				stat[0]++
			case 2:
				stat[1]++
			case 3:
				stat[2]++
			case -1:
				stat[3]++
			default:
				stat[4]++
			}
		} else {
			break
		}
	}
}

func queryTest(cmd *cobra.Command, args []string) {
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	result := make(chan int)
	stat := make([]int, 5)
	if len(args) == 1 {
		if utils.IsValidIP(args[0]) {
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			go response_stat(stat, result, &wg2)
			for i := 0; i < query_packets; i++ {
				var qname string
				if query_random_prefix {
					qname = strings.Join([]string{timestamp, strconv.Itoa(i)}, "-") + query_name
				} else {
					qname = query_name
				}
				go concurrent_request(args[0], qname, query_port, query_type, query_class, result, &wg1)
				time.Sleep(time.Second / time.Duration(query_rate))
			}
			wg1.Wait()
			close(result)
			wg2.Wait()
			fmt.Printf("total: %v noerr: %v servfail: %v nx: %v error: %v other: %v\n", query_packets, stat[0], stat[1], stat[2], stat[3], stat[4])
		}
	}
}

func init() {
	queryCmd.Flags().BoolVarP(&query_random_prefix, "random", "R", false, "random query name")
	queryCmd.Flags().StringVarP(&query_port, "port", "p", "53", "query target port")
	queryCmd.Flags().StringVarP(&query_name, "name", "d", "example.com", "query name in DNS requests")
	queryCmd.Flags().Uint16VarP(&query_type, "type", "t", 1, "target record type")
	queryCmd.Flags().Uint16VarP(&query_class, "class", "c", 1, "target record class")
	queryCmd.Flags().IntVarP(&query_rate, "rate", "r", 10, "request rate (packets per second)")
	queryCmd.Flags().IntVarP(&query_packets, "packets", "n", 100, "total request number")
	rootCmd.AddCommand(queryCmd)
}
