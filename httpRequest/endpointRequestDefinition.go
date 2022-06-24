package httpRequest

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

const AuthorizationHeaderName = "Authorization"

type Header struct {
	name  string
	value string
}

func NewHeader(name, value string) *Header {
	return &Header{name: name, value: value}
}

func (r *Header) String() string {
	if r == nil {
		return "<none>"
	}
	return r.name + ": " + r.value
}

type RequestDefinition struct {
	method        *Method
	uri           string
	authorization string
	headers       []*Header
	body          string
	summary       string
}

func Parse(index int, doc string) (*RequestDefinition, error) {
	pd := &parseData{request: &RequestDefinition{}, lines: strings.Split(doc, "\n")}
	pd.consumeLeadingLines()
	if pd.allLinesConsumed() {
		return nil, nil
	}
	pd.consumeMethodAndUriLine()
	pd.consumeHeaderLines()
	pd.consumeEmptyLine()
	pd.consumeBodyLines()
	pd.consumeTrailingEmptyLines()
	if pd.err != nil {
		return nil, fmt.Errorf("doc[%d]: %w", index, pd.err)
	}
	pd.buildSummary()
	return pd.request, nil
}

func (r *RequestDefinition) Method() *Method {
	return r.method
}

func (r *RequestDefinition) Uri() string {
	return r.uri
}

func (r *RequestDefinition) Authorization() string {
	return r.authorization
}

func (r *RequestDefinition) Headers() []*Header {
	return r.headers
}

func (r *RequestDefinition) Body() string {
	return r.body
}

func (r *RequestDefinition) Summary() string {
	return r.summary
}

func (r *RequestDefinition) ShouldBeBody() bool {
	return r.method.BodyNormallyExpected() // method is nil safe
}

func (r *RequestDefinition) String() string {
	if r == nil {
		return "<none>\n"
	}
	sb := &strings.Builder{}
	_, _ = sb.WriteString(r.method.Name()) // method is nil safe
	_ = sb.WriteByte(' ')
	_, _ = sb.WriteString(r.uri)
	_ = sb.WriteByte('\n')
	strAddHeader(sb, "Authorization", r.authorization)
	for _, header := range r.headers {
		strAddHeader(sb, header.name, header.value)
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
		_, _ = sb.WriteString(": ")
		_, _ = sb.WriteString(value)
		_ = sb.WriteByte('\n')
	}
}

func (r *RequestDefinition) headerValue(name string) string {
	for _, header := range r.headers {
		if header.name == name {
			return header.value
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
		r.headers = append(r.headers, NewHeader(name, value))
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

func (r *parseData) currentLine() (exists bool, line string) {
	if !r.allLinesConsumed() {
		return true, strings.TrimRightFunc(r.lines[r.index], unicode.IsSpace)
	}
	return false, ""
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

func (r *parseData) nextLine() (exists bool, line string) {
	r.currentLineSuccessfullyProcessed()
	return r.currentLine()
}

func (r *parseData) consumeLeadingLines() {
	exists, line := r.currentLine()
	for exists && (r.err == nil) {
		if (line == "") || strings.HasPrefix(line, "#") {
			exists, line = r.nextLine()
		} else {
			return
		}
	}
}

func (r *parseData) consumeMethodAndUriLine() {
	if r.err != nil {
		return
	}
	exists, line := r.currentLine()
	if exists {
		at := strings.Index(line, ` `)
		if 3 <= at {
			r.request.method, r.err = NormalizeAndValidateMethod(line[:at])
			if r.err == nil {
				restOfLine := line[at+1:]
				at = strings.Index(restOfLine, `/`)
				if at == 0 {
					r.request.uri = restOfLine
				} else {
					r.uniquer = restOfLine[:at]
					r.request.uri = restOfLine[at:]
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
		exists, line := r.currentLine()
		if !exists || line == "" { // blank line -> end of Headers section
			return
		}
		at := strings.Index(line, `: `)
		if at < 1 {
			r.err = errors.New("expected a 'header' line, but got: " + line)
		} else {
			name := line[:at]
			value := line[at+2:]
			r.currentLineProcessed(r.request.addHeader(name, value))
		}
	}
}

func (r *parseData) consumeEmptyLine() {
	exists, line := r.currentLine()
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
	exists, line := r.currentLine()
	for exists && (r.err == nil) {
		if line == "" {
			exists, line = r.nextLine()
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
