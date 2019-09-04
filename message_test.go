package main

import (
	"encoding/json"
	"testing"

	lorem "github.com/drhodes/golorem"
)

func TestDecodeJSON(t *testing.T) {

	data := `{
		"sender": "bounce@example.com",
		"from": {"name": "Jérôme", "address": "jerome@example.com"},
		"to": "john@hotmail.com",
		"cc": {"name": "Simon", "address": "simon@hotmail.com"},
		"bcc": [{"name": "Lee", "address": "Lee@hotmail.com"}, {"name": "Jim", "address": "jim@hotmail.com"}],
		"subject": "Test message",
		"text": "` + lorem.Paragraph(200, 200) + `"
		}`

	msg := new(message)
	err := json.Unmarshal([]byte(data), msg)
	if err != nil {
		t.Fatal(err)
	}
	//encode(os.Stdout, msg)
}
