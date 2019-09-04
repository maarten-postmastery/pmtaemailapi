package main

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"
)

type message struct {
	Sender string

	From    *address // mandatory
	ReplyTo *address

	To, Cc, Bcc *addressList

	Subject string

	HTML string
	Text string
}

type address struct {
	mail.Address
}

func (a *address) UnmarshalJSON(data []byte) error {
	var addr mail.Address // parent type to avoid recursion
	if err := json.Unmarshal(data, &addr); err == nil {
		a.Address = addr
		return nil
	}
	var email string
	if err := json.Unmarshal(data, &email); err == nil {
		a.Address = mail.Address{Address: email}
		return nil
	}
	return fmt.Errorf("cannot decode %q as address", string(data))
}

func (a *address) String() string {
	return a.Address.String()
}

type addressList []*address

func (l *addressList) UnmarshalJSON(data []byte) error {
	var list []*address // parent type to avoid recursion
	if err := json.Unmarshal(data, &list); err == nil {
		*l = addressList(list)
		return nil
	}
	var addr address
	if err := json.Unmarshal(data, &addr); err == nil {
		*l = addressList{&addr}
		return nil
	}
	return fmt.Errorf("cannot decode %q as address list", string(data))
}

func (l addressList) String() string {
	s := make([]string, len(l))
	for i, a := range l {
		s[i] = a.String()
	}
	return strings.Join(s, ", ")
}

func (l addressList) recipients() []string {
	s := make([]string, len(l))
	for i, a := range l {
		s[i] = a.Address.Address
	}
	return s
}
