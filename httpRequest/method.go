package httpRequest

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var UrlPrefix string

type Method struct {
	name        string
	expectBody  bool
	sendRequest func(uri string, headers []*Header, body string) (status int, respBody string, err error)
}

func (r *Method) Name() string {
	if r == nil {
		return "noMethod"
	}
	return r.name
}

func (r *Method) BodyNormallyExpected() bool {
	if r == nil {
		return false
	}
	return r.expectBody
}

func (r *Method) SendRequest(uri string, headers []*Header, body string) (status int, respBody string, err error) {
	return r.sendRequest(uri, headers, body)
}

var (
	MethodGET    = &Method{name: "GET", sendRequest: getIgnoreBody}
	MethodPUT    = &Method{name: "PUT", expectBody: true, sendRequest: Put}
	MethodPOST   = &Method{name: "POST", expectBody: true, sendRequest: Post}
	MethodPATCH  = &Method{name: "PATCH", expectBody: true, sendRequest: Patch}
	MethodDELETE = &Method{name: "DELETE", sendRequest: deleteIgnoreBody}
)

var (
	methods = []*Method{MethodGET, MethodPUT, MethodPOST, MethodPATCH, MethodDELETE}
)

func NormalizeAndValidateMethod(httpMethod string) (*Method, error) {
	normalizedMethod := strings.ToUpper(strings.TrimSpace(httpMethod))
	for _, m := range methods {
		if m.name == normalizedMethod {
			return m, nil
		}
	}
	sb := strings.Builder{}
	sb.WriteString("no such method '")
	sb.WriteString(normalizedMethod)
	sb.WriteString("', valid options are")
	prefix := ": "
	for _, m := range methods {
		sb.WriteString(prefix)
		sb.WriteString(m.name)
		prefix = ", "
	}
	return nil, errors.New(sb.String())
}

//goland:noinspection GoUnusedParameter
func getIgnoreBody(uri string, headers []*Header, ignoredBody string) (status int, respBody string, err error) {
	return Get(uri, headers)
}

//goland:noinspection GoUnusedParameter
func deleteIgnoreBody(uri string, headers []*Header, ignoredBody string) (status int, respBody string, err error) {
	return Delete(uri, headers)
}

func Get(uri string, headers []*Header) (status int, respBody string, err error) {
	return handleNoBodyRequest("GET", uri, headers)
}

func Post(uri string, headers []*Header, reqBody string) (status int, respBody string, err error) {
	return handleBodyRequest("POST", uri, headers, reqBody)
}

func Put(uri string, headers []*Header, reqBody string) (status int, respBody string, err error) {
	return handleBodyRequest("PUT", uri, headers, reqBody)
}

func Patch(uri string, headers []*Header, reqBody string) (status int, respBody string, err error) {
	return handleBodyRequest("PATCH", uri, headers, reqBody)
}

func Delete(uri string, headers []*Header) (status int, respBody string, err error) {
	return handleNoBodyRequest("DELETE", uri, headers)
}

func handleNoBodyRequest(method, uri string, headers []*Header) (status int, respBody string, err error) {
	return handleRequest(method, uri, headers, nil)
}

func handleBodyRequest(method, uri string, headers []*Header, reqBody string) (status int, respBody string, err error) {
	return handleRequest(method, uri, headers, ioutil.NopCloser(strings.NewReader(reqBody)))
}

func handleRequest(method, uri string, headers []*Header, body io.ReadCloser) (status int, respBody string, err error) {
	var request *http.Request
	request, err = http.NewRequest(method, UrlPrefix+uri, body)
	if err != nil {
		return
	}

	for _, header := range headers {
		request.Header.Add(header.name, header.value)
	}

	var response *http.Response
	response, err = client.Do(request)
	if response == nil {
		status = -1
		return // presumably with an error!
	}
	status = response.StatusCode
	contentLength := response.ContentLength
	if contentLength < 0 {
		return
	}
	readCloser := response.Body
	if (contentLength > 0) && (err == nil) {

		respBody = extractResponseBodyOrFragment(readCloser, response.Header.Get("Content-Type"))
	}
	_ = readCloser.Close()
	return
}

func extractResponseBodyOrFragment(bodyReadCloser io.ReadCloser, contentType string) string {
	defer func() {
		_, _ = io.Copy(ioutil.Discard, bodyReadCloser)
	}()

	if contentType == "text/css" {
		return "'css'"
	}
	if contentType == "image/png" {
		return "'png'"
	}
	collector := make([]byte, 0, 1024)
	for len(collector) < 1024 {
		buffer := make([]byte, 1024-len(collector))
		bytesRead, err := bodyReadCloser.Read(buffer)
		if bytesRead > 0 {
			collector = append(collector, buffer[:bytesRead]...)
		}
		if err != nil {
			if err != io.EOF {
				return "WTH: " + err.Error()
			}
			break
		}
	}
	if len(collector) < 1024 {
		return string(collector)
	}
	return findDisplayablePortion(collector)
}

func findDisplayablePortion(buffer []byte) string {
	for i := len(buffer) - 1; 0 <= i; i-- {
		aByte := buffer[i]
		if aByte < 128 { // no hi bit, so assume Ascii
			return string(buffer[:i+1]) + "..."
		}
		if utf8MultiByteStart <= aByte {
			ok, result := checkUTF8(buffer, i)
			if ok {
				return result + "..."
			}
		}
	}
	return "......"
}

func checkUTF8(buffer []byte, utf8StartOffset int) (ok bool, result string) {
	aByte := buffer[utf8StartOffset]
	if (aByte & bytes2Mask) == bytes2Value {
		return checkUTF8plus(buffer, utf8StartOffset, 1)
	}
	if (aByte & bytes3Mask) == bytes3Value {
		return checkUTF8plus(buffer, utf8StartOffset, 2)
	}
	if (aByte & bytes4Mask) == bytes4Value {
		return checkUTF8plus(buffer, utf8StartOffset, 3)
	}
	return // !ok
}

var (
	utf8MultiByteStart byte = 128 + 64

	extensionByteValue byte = 128
	extensionByteMask       = extensionByteValue + 64

	bytes2Value = extensionByteMask
	bytes2Mask  = bytes2Value + 32
	bytes3Value = bytes2Mask
	bytes3Mask  = bytes3Value + 16
	bytes4Value = bytes3Mask
	bytes4Mask  = bytes4Value + 8
)

func checkUTF8plus(buffer []byte, utf8StartOffset int, utf8extensionBytes int) (ok bool, result string) {
	offset := utf8StartOffset + utf8extensionBytes
	if len(buffer) <= offset {
		return // !ok
	}
	for ; offset > utf8StartOffset; offset-- {
		aByte := buffer[offset]
		if (aByte & extensionByteMask) != extensionByteValue {
			return // !ok
		}
	}
	return true, string(buffer[:utf8StartOffset+utf8extensionBytes+1])
}

var noRedirectChecker = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

var timeout = time.Second * 10

// ignore expired SSL certificates
var transCfg = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

// do NOT follow Redirects
var client = &http.Client{Transport: transCfg, CheckRedirect: noRedirectChecker, Timeout: timeout}
