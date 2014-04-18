//package rfc2822 implements an RFC2822 parser. It does not (yet) support the entire standard.
package rfc2822

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Header struct {
	Key, Value string
}

//Message type encapsulates one RFC 2822 message.
type Message struct {
	Headers map[string][]Header
	Body    []string
}

//GetHeader retrieves an unstructured header value by its name, or an error
//if the requested header does not exist. If more than one header of a
//given name exists, the last one is returned.
func (msg *Message) GetHeader(header string) (string, error) {
	if val, good := msg.Headers[strings.ToLower(header)]; good {
		return val[len(val)-1].Value, nil
	}
	return "", errors.New("NOTFOUND")
}

//AddHeader adds a single-line header to the message.
func (msg *Message) AddHeader(name string, value string) {
	msg.Headers[strings.ToLower(name)] = append(msg.Headers[strings.ToLower(name)], Header{Key: name, Value: value})
}

//AddMultiHeader adds a multi-line header to the message. Indentation is automatic.
func (msg *Message) AddMultiHeader(name string, value []string) {
	val := ""
	for i := 0; i < len(value); i++ {
		if i != 0 {
			val += "\r\n        " //8 spaces or 2 tabs
		}
		val += value[i]
	}
	msg.Headers[strings.ToLower(name)] = append(msg.Headers[strings.ToLower(name)], Header{Key: name, Value: val})
}

//GetHeaders returns one or more unstructured headers for a name, or an error
//if the requested header does not exist. When more than one header exists
//for a name, they are returned in the order they appear in the message.
func (msg *Message) GetHeaders(header string) ([]Header, error) {
	if headers, good := msg.Headers[strings.ToLower(header)]; good {
		return headers, nil
	}
	return []Header{}, errors.New("NOTFOUND")
}

//GetBody retrieves a message body if it exists, or an error if not.
func (msg *Message) GetBody() (string, error) {
	if len(msg.Body) < 1 {
		return "", errors.New("NOTFOUND")
	}
	return strings.Join(msg.Body, " "), nil
}

func (msg *Message) HeadersText() string {
	result := ""
	for _, header := range msg.Headers {
		for _, head := range header {
			result += head.Key + ": "
			result += head.Value + "\r\n"
		}
	}
	return result
}

//Text returns a string representing the whole message.
func (msg *Message) Text() string {
	result := ""
	for _, header := range msg.Headers {
		for _, head := range header {
			result += head.Key + ": "
			result += head.Value + "\r\n"
		}
	}
	if result != "" {
		result += "\r\n"
	}
	result += strings.Join(msg.Body, "")
	return result
}

//ReadFile parses an RFC 2822 formatted input and returns a Message type.
func Read(reader io.Reader) (*Message, error) {
	//buff := bufio.NewReader(reader)
	headers := make(map[string][]Header)

	var (
		key, val, lowerKey string
		lineNo             int
		inContent          bool = false
		body               []string
	)

	//Remove junk (white lines) at the start of the message - empty messages are not contemplated
	bufd, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	cleaned := strings.TrimLeft(string(bufd), "\r\n\t\v\f\b\a ")

	buff := bufio.NewReader(strings.NewReader(cleaned))

	for {
		line, ioerr := buff.ReadString('\n')
		lineNo++

		if ioerr != nil {
			if ioerr != io.EOF {
				return nil, ioerr
			}
			if len(line) == 0 {
				break
			}
		}

		switch {

		case inContent:
			body = append(body, line)

		case strings.HasPrefix(line, "\n") || strings.HasPrefix(line, "\r\n"):
			inContent = true
			continue

		case strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t"): //a field-body continuation?
			if len(key) == 0 {
				return nil, errors.New(strconv.Itoa(lineNo) + ":no match for continuation")
			}
			// Concatenate onto previous value
			val = fmt.Sprintf("%s\r\n        %s", val, strings.TrimSpace(line))
			field := Header{key, val}
			headers[lowerKey][len(headers[lowerKey])-1] = field

		default:
			if i := strings.Index(line, ":"); i > 0 {
				key = strings.TrimSpace(line[0:i])
				lowerKey = strings.ToLower(key)
				val = strings.TrimSpace(line[i+1:])
				field := Header{key, val}

				if _, has := headers[lowerKey]; has {
					headers[lowerKey] = append(headers[lowerKey], field)
				} else {
					headers[lowerKey] = []Header{field}
				}
			} else {
				return nil, errors.New(strconv.Itoa(lineNo) + ":cannot parse field")
			}
		}
	}

	return &Message{headers, body}, nil
}

func ReadString(text string) (*Message, error) {
	return Read(strings.NewReader(text))
}

// ReadFile parses an RFC 2822 formatted file and returns a Message type.
func ReadFile(fname string) (*Message, error) {
	file, err := os.Open(fname)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	return Read(file)
}

// vi: ai sw=4 ts=4 tw=0 et
