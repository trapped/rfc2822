package rfc2822

import (
	"os"
	"testing"
)

const (
	DATA = "testdata0.txt"
)

var testHeaders = map[string]string{
	"Header0": "Value0",
	"Header1": "Value1 Value1",
	"Header2": "Value2\n Value2\n Value2",
	"Header3": "Value3 Value3\n Value3\n Value3",
}

func TestManyHeaders(t *testing.T) {
	var (
		msg *Message
		err error
		//headers []Header
	)

	if msg, err = ReadFile(DATA); err != nil {
		t.Error(err)
	}

	/*if headers, err = msg.GetHeaders("header0"); err == nil {
		if headers[0].Value != "Unexpected" {
			t.Errorf("expected \"Unexpected\", got \"%s\"", headers[0].Value)
		}
		if headers[1].Value != "Value0" {
			t.Errorf("expected \"Value0\", got \"%s\"", headers[0].Value)
		}
	} else {
		t.Error(err)
	}*/

	t.Log("Headers:\r\n", msg.Headers)
	t.Log("Body:\r\n", msg.Body)

	file, err := os.OpenFile("/home/giorgio/go/src/github.com/trapped/rfc2822/recomposed.txt", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		t.Log("Error when recomposing:", err)
	}

	t.Log("Adding multiline header")
	msg.AddMultiHeader("Example", []string{"line 0", "line 1", "line 2"})

	_, err = file.WriteString(msg.Text())
	if err != nil {
		t.Log("Error when recomposing:", err)
	}
	t.Log("Recomposed: written to file")
}

func TestParse(t *testing.T) {
	/*var (
		msg *Message
		err error
	)

	if msg, err = ReadFile(DATA); err != nil {
		t.Error(err)
	}	for key, testValue := range testHeaders {
		if value, err := msg.GetHeader(key); err != nil {
			t.Error(err)
		} else {
			if value != testValue {
				t.Errorf("%s returned %s, expected %s", key, value, testValue)
			}
		}
	}*/
}

// vi: ai sw=4 ts=4 tw=0 et
