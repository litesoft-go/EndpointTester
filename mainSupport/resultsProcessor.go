package mainSupport

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (r *ResultsTracker) Add(result *CallResult) {
	r.results = append(r.results, result)
	r.calls++
	r.requestCounter++
	if r.requestCounter == r.requestsPerUser {
		r.requestCounter = 0
		print(".")
	}
}

func NewResultsTracker(users, requestDefinitions, passes int,
	pathRequestsDataConfig, urlPrefix, spoofedIpHeader string) *ResultsTracker {
	requestsPerUser := requestDefinitions * passes
	expectedCalls := users * requestsPerUser
	fmt.Println("Requests file: ", pathRequestsDataConfig)
	fmt.Println("URL Prefix: ", urlPrefix)
	fmt.Println("IP Header (spoofing): ", spoofedIpHeader)
	fmt.Println()
	fmt.Println(users, "Users")
	fmt.Println(passes, "Passes per User")
	fmt.Println(requestDefinitions, "Requests per Pass")
	fmt.Println("expected calls:", expectedCalls)
	return &ResultsTracker{results: make([]*CallResult, 0, expectedCalls), requestsPerUser: requestsPerUser}
}

func (r *CallResult) String() string {
	return r.UserID + " " + r.PathClue + " " + r.Duration.String() + " (" + strconv.Itoa(r.Status) + ") " + r.ResponseBodyClue
}

type CallResult struct {
	UserID           string
	PathClue         string
	Duration         time.Duration
	Status           int
	ResponseBodyClue string
	Error            error
}

type ResultsTracker struct {
	results         []*CallResult
	requestsPerUser int
	requestCounter  int
	calls           int
}

func (r *ResultsTracker) Report(concurrencyError error) {
	var lastErr error
	if r.calls != 0 {
		var totalDuration time.Duration
		m := make(map[string]*PathTracker)
		for _, cr := range r.results {
			if cr.Error != nil {
				lastErr = cr.Error
			} else {
				totalDuration += cr.Duration
				tracker, present := m[cr.PathClue]
				if !present {
					m[cr.PathClue] = &PathTracker{calls: 1, callsByStatus: map[int]int{cr.Status: 1},
						smallestDuration:    cr.Duration,
						LongestDuration:     cr.Duration,
						accumulatedDuration: cr.Duration,
					}
				} else {
					tracker.update(cr)
				}
			}
		}
		if lastErr == nil {
			avgDuration := averageDuration(totalDuration, r.calls)
			reqsPerSec := int64(time.Second) / int64(avgDuration)
			fmt.Println("\n", r.calls, " calls at (avg) ", avgDuration, " -> ", totalDuration, " or ", reqsPerSec, "r/s")
			fmt.Println("Path Data:")
			keys := sortKeyStrings(m)
			for _, pathClue := range keys {
				fmt.Println("   ", pathClue)
				m[pathClue].print("  ")
			}
		}
	}
	errorMsg := chkErrors(lastErr, concurrencyError)
	if errorMsg != "" {
		_, _ = os.Stderr.WriteString(errorMsg + "\n")
		os.Exit(1)
	}
}

func chkErrors(lastErr, concurrencyError error) (errorMsg string) {
	lErrMsg := checkNormalizeErrorAndPrefix("", lastErr)
	cErrMsg := checkNormalizeErrorAndPrefix(lErrMsg, concurrencyError)
	return lErrMsg + cErrMsg
}

func checkNormalizeErrorAndPrefix(prevText string, err error) (result string) {
	if err != nil {
		msg := strings.TrimSpace(err.Error())
		if msg != "context canceled" {
			if prevText == "" {
				result = "\n***** Error: " + msg
			} else {
				result = " AND " + msg
			}
		}
	}
	return
}

type PathTracker struct {
	calls               int
	callsByStatus       map[int]int
	smallestDuration    time.Duration
	LongestDuration     time.Duration
	accumulatedDuration time.Duration
}

func (r *PathTracker) update(cr *CallResult) {
	r.calls++
	r.callsByStatus[cr.Status]++
	duration := cr.Duration
	r.accumulatedDuration += duration
	if duration < r.smallestDuration {
		r.smallestDuration = duration
	}
	if r.LongestDuration < duration {
		r.LongestDuration = duration
	}
}

func (r *PathTracker) print(indent string) {
	fmt.Println(indent, indent, r.calls, "calls at (avg) ",
		averageDuration(r.accumulatedDuration, r.calls), " -> ", r.accumulatedDuration,
		" range ", r.smallestDuration, " <> ", r.LongestDuration)
	statuses := sortKeyInts(r.callsByStatus)
	for _, status := range statuses {
		fmt.Println(indent, indent, indent, status, ": ", r.callsByStatus[status])
	}
}

func averageDuration(accumulatedDuration time.Duration, divideBy int) time.Duration {
	return time.Duration(int64(accumulatedDuration) / int64(divideBy))
}

func sortKeyStrings(mapping map[string]*PathTracker) []string {
	keys := make([]string, 0, len(mapping))

	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortKeyInts(mapping map[int]int) []int {
	keys := make([]int, 0, len(mapping))

	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
