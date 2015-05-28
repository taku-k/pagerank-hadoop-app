package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type FindPerson struct {
	queue   chan string
	done    chan bool
	mu      *sync.Mutex
	count   int
	ranking map[float64]string
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// Check wheather the person is alive or dead by wikipage body
func isAlive(body []byte) bool {
	aPat := "生年月日|出生|生誕"
	dPat := "死去|死没|没年月日"
	isAlive, err := regexp.Match(aPat, body)
	errCheck(err)
	isDead, err := regexp.Match(dPat, body)
	errCheck(err)
	if isAlive && !isDead {
		return true
	}
	return false
}

func splitLine(line string) (float64, string) {
	split := strings.Split(line, "\t")
	rank, err := strconv.ParseFloat(split[0], 64)
	errCheck(err)
	query := split[1]
	return rank, query
}

func initFindPerson() *FindPerson {
	fp := new(FindPerson)
	fp.queue = make(chan string)
	fp.done = make(chan bool)
	fp.mu = new(sync.Mutex)
	fp.count = 0
	fp.ranking = make(map[float64]string)

	return fp
}

// Execute query
func executeQuery(fp *FindPerson, line string) {
	rank, query := splitLine(line)
	resp, err := http.Get("http://ja.wikipedia.org/wiki/" + query)
	errCheck(err)

	body, err := ioutil.ReadAll(resp.Body)
	errCheck(err)

	alive := isAlive(body)
	if !alive {
		return
	}

	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.count += 1
	fp.ranking[rank] = query
	if fp.count == 10 {
		fp.done <- true
	}
}

func fetcher(fp *FindPerson) {
	for {
		select {
		case line := <-fp.queue:
			go executeQuery(fp, line)
		}
	}
}

func main() {
	fp := initFindPerson()
	f, err := os.Open(os.Args[1])
	errCheck(err)
	defer f.Close()

	reader := bufio.NewReaderSize(f, 4096)
	go fetcher(fp)

	var line string
	for {
		select {
		case <-fp.done:
			break
		default:
			line, _ = reader.ReadString('\n')
			fp.queue <- line
		}
	}

	for key, value := range fp.ranking {
		fmt.Println(key, value)
	}
}
