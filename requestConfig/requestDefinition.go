package requestConfig

import (
	"endpointTester/httpRequest"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type RequestDefinition struct {
	sequenceValue        string
	sequenceContinuation bool
	expectedStatusCodes  []int
	method               *httpRequest.Method
	uri                  string
	authorization        string
	headers              []*httpRequest.Header
	body                 string
	summary              string
}

func Parse(index int, block string) (*RequestDefinition, error) {
	pd := &parseData{request: &RequestDefinition{}, lines: strings.Split(block, "\n")}
	pd.consumeSequenceIndicator()
	pd.consumeLeadingLines()
	if pd.allLinesConsumed() {
		return nil, nil
	}
	pd.consumeExpectLine()
	pd.consumeMethodAndUriLine()
	pd.consumeHeaderLines()
	pd.consumeEmptyLine()
	pd.consumeBodyLines()
	pd.consumeTrailingEmptyLines()
	if pd.err != nil {
		return nil, fmt.Errorf("block[%d]: %w", index, pd.err)
	}
	pd.buildSummary()
	return pd.request, nil
}

func (r *RequestDefinition) Method() *httpRequest.Method    { return r.method }
func (r *RequestDefinition) Uri() string                    { return r.uri }
func (r *RequestDefinition) Authorization() string          { return r.authorization }
func (r *RequestDefinition) Headers() []*httpRequest.Header { return r.headers }
func (r *RequestDefinition) Body() string                   { return r.body }
func (r *RequestDefinition) Summary() string                { return r.summary }

func (r *RequestDefinition) IsSequenceStart() bool { return !r.sequenceContinuation }

func (r *RequestDefinition) IsExpectedStatusCode(statusCode int) bool {
	if len(r.expectedStatusCodes) == 0 {
		return true
	}
	for _, expectedStatusCode := range r.expectedStatusCodes {
		if statusCode == expectedStatusCode {
			return true
		}
	}
	return false
}

func (r *RequestDefinition) ShouldBeBody() bool {
	return r.method.BodyNormallyExpected() // method is nil safe
}

func (r *RequestDefinition) String() string {
	if r == nil {
		return "<none>\n"
	}
	sb := &strings.Builder{}
	_, _ = sb.WriteString(BlockStart)
	if r.sequenceValue != "" {
		_, _ = sb.WriteString(SequenceLineStartsWith)
		_ = sb.WriteByte(' ')
		_, _ = sb.WriteString(r.sequenceValue)
	}
	_ = sb.WriteByte('\n')
	if len(r.expectedStatusCodes) != 0 {
		_, _ = sb.WriteString(ExpectLineStartsWith)
		for i, statusCode := range r.expectedStatusCodes {
			if i != 0 {
				_, _ = sb.WriteString(ExpectStatusCodeSeparator)
			}
			_ = sb.WriteByte(' ')
			_, _ = sb.WriteString(strconv.Itoa(statusCode))
		}
		_, _ = sb.WriteString(r.sequenceValue)
		_ = sb.WriteByte('\n')
	}
	_, _ = sb.WriteString(r.method.Name()) // method is nil safe
	_ = sb.WriteByte(' ')
	_, _ = sb.WriteString(r.uri)
	_ = sb.WriteByte('\n')
	strAddHeader(sb, AuthorizationHeaderName, r.authorization)
	for _, header := range r.headers {
		strAddHeader(sb, header.Name(), header.Value())
	}
	_ = sb.WriteByte('\n')
	if r.body != "" {
		_, _ = sb.WriteString(r.body)
		_ = sb.WriteByte('\n')
	}
	return sb.String()
}

func strAddHeader(sb *strings.Builder, name string, value string) {
	if value != "" {
		_, _ = sb.WriteString(name)
		_, _ = sb.WriteString(HeaderLineNameValueSeparator)
		_ = sb.WriteByte(' ')
		_, _ = sb.WriteString(value)
		_ = sb.WriteByte('\n')
	}
}

func (r *RequestDefinition) headerValue(name string) string {
	for _, header := range r.headers {
		if header.Name() == name {
			return header.Value()
		}
	}
	return ""
}

func (r *RequestDefinition) addHeaderNonEmpty(name, value string) (err error) {
	existingValue := ""
	if strings.EqualFold(AuthorizationHeaderName, name) {
		existingValue = r.authorization
		r.authorization = value
	} else {
		existingValue = r.headerValue(name)
		r.headers = append(r.headers, httpRequest.NewHeader(name, value))
	}
	if existingValue != "" {
		err = errors.New("header '" + name + "' already has a value of: " + existingValue)
	}
	return
}

func (r *RequestDefinition) addHeader(name, value string) error {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	if name == "" {
		return errors.New("no 'name' for header with value: '" + value + "'")
	}
	if value == "" {
		return errors.New("no 'value' for header with name: " + name)
	}
	return r.addHeaderNonEmpty(name, value)
}

func (r *RequestDefinition) addBody(body string) {
	r.body = body
}

func (r *RequestDefinition) addExpectedStatusCode(value int) {
	r.expectedStatusCodes = append(r.expectedStatusCodes, value)
}

type parseData struct {
	request *RequestDefinition
	lines   []string
	index   int
	uniquer string
	err     error
}

func (r *parseData) allLinesConsumed() bool {
	return len(r.lines) <= r.index
}

func (r *parseData) currentLineTrimmed() (exists bool, line string) {
	exists, line = r.currentLine()
	line = strings.TrimLeftFunc(line, unicode.IsSpace)
	return
}

func (r *parseData) nextLineTrimmed() (exists bool, line string) {
	r.currentLineSuccessfullyProcessed()
	return r.currentLineTrimmed()
}

func (r *parseData) currentLine() (exists bool, line string) {
	if !r.allLinesConsumed() {
		return true, strings.TrimRightFunc(r.lines[r.index], unicode.IsSpace)
	}
	return false, ""
}

func (r *parseData) nextLine() (exists bool, line string) {
	r.currentLineSuccessfullyProcessed()
	return r.currentLine()
}

func (r *parseData) currentLineSuccessfullyProcessed() {
	r.index++
}

func (r *parseData) currentLineProcessed(err error) {
	if err != nil {
		r.err = err
	} else {
		r.currentLineSuccessfullyProcessed()
	}
}

func (r *parseData) consumeLeadingLines() {
	exists, line := r.currentLineTrimmed()
	for exists && (r.err == nil) {
		if (line == "") || strings.HasPrefix(line, CommentLineStartsWith) {
			exists, line = r.nextLineTrimmed()
		} else {
			return
		}
	}
}

func (r *parseData) consumeSequenceIndicator() {
	if r.err != nil {
		return
	}
	exists, line := r.currentLineTrimmed()
	if exists && strings.HasPrefix(line, SequenceLineStartsWith) {
		r.request.sequenceValue = strings.TrimSpace(line[len(SequenceLineStartsWith):])
		restOfLine := strings.ToUpper(r.request.sequenceValue)
		r.request.sequenceContinuation = len(restOfLine) > 0 && restOfLine != SequenceLineToUpperSequenceStartMatch
		r.currentLineProcessed(r.err)
	}
}

func (r *parseData) consumeExpectLine() {
	if r.err != nil {
		return
	}
	exists, line := r.currentLineTrimmed()
	if exists && strings.HasPrefix(line, ExpectLineStartsWith) {
		codeStrings := strings.Split(line[len(ExpectLineStartsWith):], ExpectStatusCodeSeparator)
		for i, codeStr := range codeStrings {
			codeStr = strings.TrimSpace(codeStr)
			if codeStr != "" {
				value, err := strconv.Atoi(codeStr)
				r.request.addExpectedStatusCode(value)
				if err != nil {
					r.err = fmt.Errorf("parsing expected status code (entry[%d] of '%s') in line '%s' errored with: %w", i, codeStr, line, err)
					break
				}
			}
		}
		r.currentLineProcessed(r.err)
	}
}

func (r *parseData) consumeMethodAndUriLine() {
	if r.err != nil {
		return
	}
	exists, line := r.currentLineTrimmed()
	if exists {
		at := strings.Index(line, MethodAndUriLineMethodRestSeparator)
		if 3 <= at {
			r.request.method, r.err = httpRequest.NormalizeAndValidateMethod(line[:at])
			if r.err == nil {
				restOfLine := strings.TrimSpace(line[at+len(MethodAndUriLineMethodRestSeparator):])
				at = strings.Index(restOfLine, UriStart)
				if at == 0 {
					r.request.uri = restOfLine
				} else if at > 0 {
					r.uniquer = restOfLine[:at]
					r.request.uri = restOfLine[at:]
				} else {
					r.err = errors.New("expected 'method' and 'uri' line, but no 'uri' (starts with '" + UriStart + "') found, got: " + line)
				}
			}
			r.currentLineProcessed(r.err)
			return // with possible error
		}
	}
	if exists {
		r.err = errors.New("expected 'method' and 'uri' line, but got: " + line)
	} else {
		r.err = errors.New("no 'method' and 'uri' line")
	}
}

func (r *parseData) consumeHeaderLines() {
	for r.err == nil {
		exists, line := r.currentLineTrimmed()
		if !exists || line == "" { // blank line -> end of Headers section
			return
		}
		at := strings.Index(line, HeaderLineNameValueSeparator)
		if at < 1 {
			r.err = errors.New("expected a 'header' line, but got: " + line)
		} else {
			name := strings.TrimSpace(line[:at])
			value := strings.TrimSpace(line[at+1:])
			r.currentLineProcessed(r.request.addHeader(name, value))
		}
	}
}

func (r *parseData) consumeEmptyLine() {
	exists, line := r.currentLineTrimmed()
	if (r.err == nil) && exists {
		if line == "" { // blank line -> end of Headers section
			r.currentLineSuccessfullyProcessed()
		} else {
			r.err = errors.New("expected a blank line (ending the Headers section), but got: " + line)
		}
	}
}

func (r *parseData) consumeBodyLines() {
	if (r.err != nil) || !r.request.method.BodyNormallyExpected() {
		return
	}
	sb := strings.Builder{}
	exists, line := r.currentLine()
	for exists && (line != "") { // blank line -> end of Body section
		_, _ = sb.WriteString(line)
		_ = sb.WriteByte('\n')
		exists, line = r.nextLine()
	}
	body := sb.String()

	if body != "" {
		r.request.addBody(body)
	} else {
		r.err = errors.New("expected body, but none found")
	}
}

func (r *parseData) consumeTrailingEmptyLines() {
	exists, line := r.currentLineTrimmed()
	for exists && (r.err == nil) {
		if line == "" {
			exists, line = r.nextLineTrimmed()
		} else {
			r.err = errors.New("expected only empty trailing lines, but got: " + line)
		}
	}
}

// "GET /mfa-providers   . . . . . . . . : "
// "DELETE /mfa-providers/bae8a2a4-||... : "
// -1234567-101234567-201234567-30123456789
func (r *parseData) buildSummary() {
	workingLength := 37
	method := r.request.method.Name()
	uniquer := r.uniquer
	uri := r.request.uri
	sb := &strings.Builder{}
	_, _ = sb.WriteString(method)
	_ = sb.WriteByte(' ')
	_, _ = sb.WriteString(uniquer)
	maxUriToAdd := (workingLength - 2) - sb.Len()
	if len(uri) <= (maxUriToAdd - 1) { // Assume Ascii
		_, _ = sb.WriteString(uri)
		_ = sb.WriteByte(' ')
	} else {
		_, _ = sb.WriteString(uri[:maxUriToAdd-4])
		_, _ = sb.WriteString("||...")
	}
	_ = sb.WriteByte(' ')
	toAdd := workingLength - sb.Len()
	if (toAdd & 1) == 1 { // is Odd
		_ = sb.WriteByte(' ')
		toAdd--
	}
	for toAdd > 0 {
		_, _ = sb.WriteString(". ")
		toAdd -= 2
	}
	r.request.summary = sb.String()
}
