package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	var limit int
	flag.IntVar(&limit, "parallel", 10, "an integer specifying how many concurrent requests are allowed")
	flag.Parse()

	targets := flag.Args()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	ch := make(chan string)

	if len(targets) > limit {
		maxBatches := len(targets) / limit
		targetsPerBatch := limit
		for b := 0; b <= maxBatches; b++ {
			batchTargets := batch(targets, b, targetsPerBatch)
			runBatch(batchTargets, len(batchTargets), client, ch)
		}
	} else {
		limit = len(targets)
		runBatch(targets, limit, client, ch)
	}

}

//runBatch calls the goroutine to make requests and prints out info on the channel
func runBatch(targets []string, limit int, cl *http.Client, ch chan string) {
	for i := 0; i < limit; i++ {
		go makeRequest(cl, targets[i], ch)
	}

	for range targets {
		fmt.Println(<-ch)
	}
}

//batch defines the start and end of a batch to run
func batch(targets []string, page, pageSize int) []string {
	start := pageSize * page
	end := pageSize * (page + 1)
	if page == len(targets)/pageSize {
		return targets[start:]
	}
	return targets[start:end]
}

//MakeRequest sends a get request to url
func makeRequest(cl *http.Client, inurl string, ch chan<- string) {

	u, err := url.Parse(inurl)

	if err != nil {
		log.Fatal(err)
	}

	if u.Scheme == "" {
		inurl = "http://" + inurl
	}

	req, err := http.NewRequest(http.MethodGet, inurl, nil)
	if err != nil {
		log.Printf("failed creating request to %s with error %v", inurl, err)
		return
	}

	res, err := cl.Do(req)
	if err != nil {
		log.Printf("failed making request to %s with error %v", inurl, err)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("")
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	ch <- fmt.Sprintf("%s %s", inurl, HashResponse(body))
}

//HashResponse returns MD5 encoded string
func HashResponse(res []byte) string {
	hash := md5.Sum(res)
	return hex.EncodeToString(hash[:])
}
