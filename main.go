package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	flag           int = 0
	wg             sync.WaitGroup
	statisticTimes = make(map[string]int)
	wLock          sync.RWMutex
)

func main() {
	var dirname string
	dirname = "./newsSource"
	// dirname = "/Users/ziyuehu/Documents/go/gotaskfilecount-master/newsSource"
	myfolder := dirname
	listFile(myfolder)
}

// 列出目录下的文件
func listFile(myfolder string) {
	files, _ := ioutil.ReadDir(myfolder)
	for _, file := range files {
		if file.IsDir() {
			listFile(myfolder + "/" + file.Name())
		} else {
			if file.Name() != "main.go" && file.Name() != "output.txt" && file.Name() != "output.csv" {
				wg.Add(1)
				go readFile(myfolder + "/" + file.Name())
			}
		}
	}

	wg.Wait()

	for word, counts := range statisticTimes {

		wg.Add(1)
		go write2csv(word, counts)
	}

	wg.Wait()
	fmt.Println("Write output.txt sucess!\n")
}

// 读文件--写入map
func readFile(fileName string) {
	// goroutine - 1
	wLock.Lock()
	defer wg.Done()
	defer wLock.Unlock()
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	wordsLength := strings.Fields(string(buf))

	for counts, word := range wordsLength {
		// 判断key是否存在，这个word是字符串，这个counts是统计的word的次数。
		word, ok := statisticTimes[word]
		if ok {
			word = word
			statisticTimes[wordsLength[counts]] = statisticTimes[wordsLength[counts]] + 1
		} else {
			statisticTimes[wordsLength[counts]] = 1
		}
	}
}

// 写入csv文件--读map
func write2csv(word string, counts int) {
	wLock.Lock()
	defer wg.Done()
	defer wLock.Unlock()
	file, err := os.OpenFile("output.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open file is failed, err: ", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	if flag == 0 {
		w.Write([]string{"word", "count"})
		flag = 1
		w.Flush()
	}

	w.Write([]string{word, strconv.Itoa(counts)})
	w.Flush()
}
