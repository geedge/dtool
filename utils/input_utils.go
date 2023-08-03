package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func RetrieveLines(pool chan string, filename string) {
	cnt := 0
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("reading file ...")
	reader := bufio.NewReader(f)
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		s = strings.Trim(s, "\n")
		pool <- s
		cnt++
		if cnt%10 == 0 {
			fmt.Println(cnt)
		}
	}
	close(pool)
}
