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
)

func main() {
	name := cleanAppName(os.Args[0])
	fmt.Print(name, " vs 0.3 ")

	pathRequestsDataConfig := flag.String("requestsConfigFileAt", "./requests.txt", "path to file 'requests'")
	urlPrefix := flag.String("urlPrefix", "http://localhost:8080/uaa", "URL prefix, e.g. 'http://localhost:8080/uaa' used for local testing")
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

	requestDefinitions := mainSupport.LoadRequestsData(*pathRequestsDataConfig)
	users := httpRequest.GetUsers(*concurrency)

	resultsTracker := mainSupport.NewResultsTracker(len(users), len(requestDefinitions), *passes,
		*pathRequestsDataConfig, *urlPrefix, *spoofedIpHeader)

	resultsChannel := make(chan *mainSupport.CallResult)

	g, ctx := errgroup.WithContext(context.Background())

	initialOffset := len(requestDefinitions) - 1

	for _, user := range users {
		caller := mainSupport.NewCaller(user, requestDefinitions, *spoofedIpHeader, initialOffset, *passes, resultsChannel, ctx)
		initialOffset -= 3
		if initialOffset < 0 {
			initialOffset += len(requestDefinitions)
		}
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
