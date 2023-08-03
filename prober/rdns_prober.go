package prober

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"dtool/utils"

	"github.com/miekg/dns"
)

type RecursiveStruct struct {
	target string
	dict   map[string]bool
}

var dataset map[string][]string

func active_probe(n int, addr string) RecursiveStruct {
	target_ip := addr
	data := RecursiveStruct{target_ip, make(map[string]bool)}
	stop := 0
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	for i := 0; i < n; i++ {
		subdomain := strings.Join([]string{strings.Replace(target_ip, ".", "-", -1), "echo", strconv.Itoa(i), timestamp}, "-")
		domain := dns.Fqdn(subdomain + ".echodns.xyz")
		res, err := utils.SendQuery(addr, domain)
		if err != nil {
			stop += 1
			continue
		}
		rdns_ip, err := utils.ParseAResponse(res)
		if err == nil {
			data.dict[rdns_ip] = true
		} else {
			stop += 1
		}
		if stop == 3 {
			break
		}
	}
	return data
}

func upstream_prober(ip_pool chan string, data_pool chan RecursiveStruct, wg *sync.WaitGroup) {
	for {
		if s, ok := <-ip_pool; ok {
			data := active_probe(20, s)
			if data.dict != nil {
				data_pool <- data
			}
		} else {
			break
		}
	}
	wg.Done()
}

func create_probers(num int, ip_pool chan string, data_pool chan RecursiveStruct, wg *sync.WaitGroup) {
	for i := 0; i < num; i++ {
		wg.Add(1)
		go upstream_prober(ip_pool, data_pool, wg)
	}
}

func store_data(pool chan RecursiveStruct, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		var temp []string
		if data, ok := <-pool; ok {
			if len(data.dict) > 0 {
				for rdns := range data.dict {
					temp = append(temp, rdns)
				}
				dataset[data.target] = temp
			}
		} else {
			break
		}
	}
	wg.Done()
}

func Get_upstream_file(filename string, output string, prober_num int) {
	dataset = map[string][]string{}
	ip_pool := make(chan string, 500)
	data_pool := make(chan RecursiveStruct, 200)
	var probe_tasks sync.WaitGroup
	var store_task sync.WaitGroup

	go utils.RetrieveLines(ip_pool, filename)
	create_probers(prober_num, ip_pool, data_pool, &probe_tasks)
	go store_data(data_pool, &store_task)
	probe_tasks.Wait()
	close(data_pool)
	store_task.Wait()
	utils.OutputJSON(dataset, output, "")
}

func Get_upstream_ip(ip string) {
	dataset = make(map[string][]string)
	var temp []string
	data := active_probe(10, ip)
	for rdns := range data.dict {
		temp = append(temp, rdns)
	}
	dataset[data.target] = temp
	utils.OutputJSON(dataset, "-", "    ")
}
