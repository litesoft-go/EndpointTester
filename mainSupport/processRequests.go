package mainSupport

import (
	"context"
	"endpointTester/httpRequest"
	"endpointTester/requestConfig"
	"fmt"
	"strings"
	"time"
)

func NewCaller(user *httpRequest.UserData, requestDefs []*requestConfig.RequestDefinition,
	spoofedIpHeader string, initialOffset, passes int, messages chan<- *CallResult, ctx context.Context) *Caller {

	var requests []*requestMapping
	for _, def := range requestDefs {
		requests = append(requests, create(user, spoofedIpHeader, def))
	}
	return &Caller{
		requests:       requests,
		initialOffset:  initialOffset,
		passes:         passes,
		resultsChannel: messages,
		ctx:            ctx,
	}
}

type Caller struct {
	requests       []*requestMapping
	initialOffset  int
	passes         int
	resultsChannel chan<- *CallResult
	ctx            context.Context
}

func (r *Caller) Call() error {
	requestOffset, passes := r.initialOffset, r.passes
	var err error
	for (err == nil) && (passes > 0) {
		err = r.call(r.requests[requestOffset])
		requestOffset = r.incrementAndConstrainOffset(requestOffset)
		for (err == nil) && (requestOffset != r.initialOffset) {
			err = r.call(r.requests[requestOffset])
			requestOffset = r.incrementAndConstrainOffset(requestOffset)
		}
		passes--
	}
	if err != nil {
		r.resultsChannel <- &CallResult{Error: err}
	}
	r.resultsChannel <- nil // indicate done!
	return err
}

func (r *Caller) incrementAndConstrainOffset(offset int) int {
	offset++
	if offset < len(r.requests) {
		return offset
	}
	return 0
}

func create(user *httpRequest.UserData, spoofedIpHeader string, requestDef *requestConfig.RequestDefinition) *requestMapping {
	callerId := user.IP()
	callerIdExt := "     "
	headers := []*httpRequest.Header{httpRequest.NewHeader(spoofedIpHeader, callerId)}
	authorization := requestDef.Authorization()
	if authorization != "NONE" {
		value := authorization
		if strings.HasPrefix(value, "Bearer ") {
			value = "Bearer " + user.JWT()
			callerIdExt = "(JWT)"
		}
		headers = append(headers, httpRequest.NewHeader(requestConfig.AuthorizationHeaderName, value))
	}
	headers = append(headers, requestDef.Headers()...)
	return &requestMapping{
		userID:     callerId + callerIdExt + ":",
		definition: requestDef,
		headers:    headers,
	}
}

type requestMapping struct {
	userID     string
	definition *requestConfig.RequestDefinition
	headers    []*httpRequest.Header
}

func (r *Caller) call(request *requestMapping) error {
	result, err := callWithFields(
		request.userID,
		request.definition.Summary(),
		request.definition.Method(),
		request.definition.Uri(),
		request.definition.Body(),
		request.headers...)
	if err == nil {
		result.StatusEvaluator = request.definition
		r.resultsChannel <- result
	}
	return err
}

func callWithFields(userID, pathClue string, method *httpRequest.Method, uri, reqBody string, headers ...*httpRequest.Header) (result *CallResult, err error) {
	start := time.Now()
	var status int
	var limitedRespBody string
	status, limitedRespBody, err = method.SendRequest(uri, headers, reqBody)
	result = &CallResult{UserID: userID, PathClue: pathClue, Duration: time.Since(start), Status: status, ResponseBodyClue: limitedRespBody}
	if err != nil {
		err = fmt.Errorf("error (%w) occured processing: %v", err, result)
	}
	return
}
