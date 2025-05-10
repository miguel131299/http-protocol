package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	stateInitialized = iota
	stateDone
)

type Request struct {
	RequestLine RequestLine
	state       int // internal parser state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	const initialBufSize = 8
	buf := make([]byte, initialBufSize)
	start := 0
	end := 0
	request := &Request{}

	for {
		// Grow buffer is needed
		if end == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf[start:end])
			end = end - start
			start = 0
			buf = newBuf
		}

		// Read into buffer[end:]
		n, err := reader.Read(buf[end:])
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading from reader: %w", err)
		}
		end += n

		// Parse from unprocessed region
		nParsed, parseErr := request.parse(buf[start:end])
		if parseErr != nil {
			return nil, fmt.Errorf("parse error: %w", parseErr)
		}
		start += nParsed

		// Shift remaining data to the start if needed
		if start == end {
			start = 0
			end = 0
		}

		if err == io.EOF {
			break
		}
	}

	if request.state != stateDone {
		return nil, fmt.Errorf("incomplete HTTP request line")
	}

	return request, nil
}

func isAllUppercaseLetters(s string) bool {
	if s == "" {
		return false // or true depending on how you treat empty strings
	}

	for _, r := range s {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func parseRequestLine(str string) (*RequestLine, int, error) {
	crlfIndex := strings.Index(str, "\r\n")
	if crlfIndex == -1 {
		return nil, 0, nil // Incomplete
	}

	line := str[:crlfIndex]
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("request line must have 3 parts")
	}

	method := parts[0]
	if !isAllUppercaseLetters(method) {
		return nil, 0, fmt.Errorf("request method must only contain uppercase letters")
	}

	versionStr := parts[2]
	if versionStr != "HTTP/1.1" {
		return nil, 0, fmt.Errorf("only HTTP/1.1 is supported")
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: parts[1],
		HttpVersion:   "1.1",
	}, crlfIndex + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == stateDone {
		return 0, nil // Already done, nothing to do
	}

	str := string(data)

	requestLine, bytesConsumed, err := parseRequestLine(str)
	if err != nil {
		return 0, err
	}
	if bytesConsumed == 0 {
		// Incomplete, need more data
		return 0, nil
	}

	// Successfully parsed request line
	r.RequestLine = *requestLine
	r.state = stateDone

	return bytesConsumed, nil
}
