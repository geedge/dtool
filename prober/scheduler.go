package prober

import (
	"bufio"
	"dtool/utils"
	"fmt"
	"os"
	"sync"
)

//type ProbeTask func(string) interface{}

func output_process(output chan interface{}, file string, wg *sync.WaitGroup) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	for {
		if data, ok := <-output; ok {
			str, err := OutputHandler(data)
			if err != nil {
				fmt.Printf("Error generating output: %s\n", err)
				continue
			}
			_, err = writer.WriteString(str + "\n")
			if err != nil {
				fmt.Printf("Error writing file: %s\n", err)
			}
		} else {
			break
		}
	}
	writer.Flush()
	wg.Done()
}

func concurrent_execution[T any](fn func(string, string) T, domain string, input chan string, output chan interface{}, wg *sync.WaitGroup) {
	for {
		if ip, ok := <-input; ok {
			data := fn(ip, domain)
			output <- data
		} else {
			break
		}
	}
	wg.Done()
}

func CreateTask[T any](fn func(string, string) T, domain string, input_file string, output_file string, concurrent_num int) {
	input_pool := make(chan string, 500)
	output_pool := make(chan interface{}, 500)
	var probe_tasks sync.WaitGroup
	var store_tasks sync.WaitGroup

	go utils.RetrieveLines(input_pool, input_file)
	probe_tasks.Add(concurrent_num)
	for i := 0; i < concurrent_num; i++ {
		go concurrent_execution(fn, domain, input_pool, output_pool, &probe_tasks)
	}
	store_tasks.Add(1)
	go output_process(output_pool, output_file, &store_tasks)
	probe_tasks.Wait()
	close(output_pool)
	store_tasks.Wait()
}
