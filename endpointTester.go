package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"endpointTester/httpRequest"
	"endpointTester/mainSupport"
	"endpointTester/requestConfig"
)

func main() {
	name := cleanAppName(os.Args[0])
	fmt.Print(name, " vs 0.5 ")

	pathRequestsDataConfig := flag.String("requestsConfigFileAt", "./requests.txt", "path to file 'requests'")
	urlPrefix := flag.String("urlPrefix", "http://localhost:8080", "URL prefix, e.g. 'http://localhost:8080' used for local testing")
	spoofedIpHeader := flag.String("spoofedIpHeader", "X-Client-IP", "as each user has a fake JWT and IP address, these IP addresses facilitate cache &/ rate limiting by IP address or email in the JWT.  The options are: X-Client-IP, X-Real-IP, and X-Forwarded-For")
	concurrency := flag.Int("users", 3, "indicates the number of users (1-25) where each user is run in parallel")
	passes := flag.Int("passes", 10, "how many times the requests are processed for each user")

	flag.Parse()

	httpRequest.UrlPrefix = *urlPrefix

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("(at " + pwd + ")")

	requestDefinitions := requestConfig.LoadRequestsData(*pathRequestsDataConfig)
	users := httpRequest.GetUsers(*concurrency)

	resultsTracker := mainSupport.NewResultsTracker(len(users), len(requestDefinitions), *passes,
		*pathRequestsDataConfig, *urlPrefix, *spoofedIpHeader)

	resultsChannel := make(chan *mainSupport.CallResult)

	g, ctx := errgroup.WithContext(context.Background())

	requestOffsetSupplier := NewRequestOffsetSupplier(requestDefinitions)

	for _, user := range users {
		initialOffset := requestOffsetSupplier.NextIndex()
		caller := mainSupport.NewCaller(user, requestDefinitions, *spoofedIpHeader, initialOffset, *passes, resultsChannel, ctx)
		g.Go(caller.Call)
	}

	routines := len(users)
	for (err == nil) && (routines > 0) {
		select {
		case result := <-resultsChannel:
			if result != nil {
				resultsTracker.Add(result)
			} else {
				routines--
			}
		case <-ctx.Done():
			err = ctx.Err()
			routines = -1
		default:
			time.Sleep(2 * time.Millisecond)
		}
	}

	if err == nil {
		err = g.Wait()
	}
	resultsTracker.Report(err)
}

func cleanAppName(name string) string {
	name = "/" + name + " "
	name = name[strings.LastIndex(name, "/")+1:]
	for name[0] == '_' {
		name = name[1:]
	}
	return strings.TrimSpace(name)
}

type RequestOffsetSupplier struct {
	requestDefIndexes []int
	next              int
}

func (r *RequestOffsetSupplier) NextIndex() int {
	rv := r.requestDefIndexes[r.next]
	r.next++
	if len(r.requestDefIndexes) <= r.next {
		r.next = 0
	}
	return rv
}

func NewRequestOffsetSupplier(requestDefs []*requestConfig.RequestDefinition) *RequestOffsetSupplier {
	var orderedIndexes []int
	for i, requestDef := range requestDefs {
		if requestDef.IsSequenceStart() {
			orderedIndexes = append(orderedIndexes, i)
		}
	}

	var indexes []int
	for i := len(orderedIndexes) - 1; 0 <= i; i -= 2 {
		indexes = append(indexes, orderedIndexes[i])
	}
	for i := len(orderedIndexes) - 2; 0 <= i; i -= 2 {
		indexes = append(indexes, orderedIndexes[i])
	}

	return &RequestOffsetSupplier{requestDefIndexes: indexes, next: 0}
}
