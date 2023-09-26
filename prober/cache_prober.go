package prober

import (
	"dtool/utils"
	_ "fmt"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const query_num = 20

type CacheStruct struct {
	target string
	dict   map[int]map[string]bool
}

func RecursiveCacheProbe(ip string, sld string) CacheStruct {
	data := CacheStruct{ip, make(map[int]map[string]bool)}
	stop := 0
	time_now := strconv.FormatInt(time.Now().Unix(), 10)
	for i := 0; i < query_num; i++ {
		if stop >= 3 {
			break
		}
		subdomain := strings.Join([]string{strings.Replace(ip, ".", "-", -1), "fwd", strconv.Itoa(i), time_now}, "-")
		domain := dns.Fqdn(subdomain + "." + sld)
		res, err := utils.SendQuery(ip, domain)
		if err != nil {
			//fmt.Printf("Error sending query: %s\n", err)
			stop += 1
			continue
		}
		cache_id, rdns, err := utils.ParseCNAMEChain(res)
		if err != nil {
			//fmt.Printf("Error parsing response: %s\n", err)
			stop += 1
			continue
		}
		if data.dict[cache_id] == nil {
			data.dict[cache_id] = make(map[string]bool)
		}
		data.dict[cache_id][rdns] = true
		stop = 0
	}
	return data
}

func RecursiveCacheTest(ip string, num int) {
	res := make(map[string]map[int][]string)
	temp := make(map[int][]string)
	data := RecursiveCacheProbe(ip, "echodns.xyz")
	if len(data.dict) > 0 {
		for cache_id := range data.dict {
			for rdns := range data.dict[cache_id] {
				temp[cache_id] = append(temp[cache_id], rdns)
			}
		}
	}
	res[ip] = temp
	utils.OutputJSON(res, "-", "    ")
}

func ClientCacheProbe(ip string) {

}
