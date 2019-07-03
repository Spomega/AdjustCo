package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var limit = flag.Int("parallel", 10, "an integer specifying how many concurrent requests are allowed")

func main() {

	flag.Parse()
	ctx := context.Background()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	targets := make(chan string)

	if *limit < 1 {
		log.Fatal("parallel flag must be at least one")
	}

	go func() {

		defer close(targets)

		for _, target := range flag.Args() {
			targets <- target
		}
	}()

	var wg sync.WaitGroup

	results := make(chan string)

	wg.Add(*limit)

	for i := 0; i < *limit; i++ {

		go func() {
			defer wg.Done()

			for target := range targets {
				result, err := makeRequest(ctx, client, target)

				if err != nil {
					log.Printf("error from %q: %v", target, err)
					continue
				}

				results <- fmt.Sprintf("%s\t%s", target, result)
			}
		}()

	}

	go func() {
		wg.Wait()

		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

}

//MakeRequest sends a get request to url
func makeRequest(ctx context.Context, cl *http.Client, inurl string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	u, err := url.Parse(inurl)

	if err != nil {
		log.Fatal(err)
	}

	if u.Scheme == "" {
		inurl = "http://" + inurl
	}

	req, err := http.NewRequest(http.MethodGet, inurl, nil)
	if err != nil {
		return "", err
	}

	res, err := cl.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return "", errors.New(res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return HashResponse(body), nil
}

//HashResponse returns MD5 encoded string
func HashResponse(res []byte) string {
	hash := md5.Sum(res)
	return hex.EncodeToString(hash[:])
}
