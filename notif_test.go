package main

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestParse(t *testing.T) {
	tt := []struct {
		input  string
		expext notif
	}{
		{
			"#channel <user> message",
			notif{
				Channel: "#channel",
				User:    "user",
				Message: "message",
			},
		},
		{
			"#channel <user> ",
			notif{
				Channel: "#channel",
				User:    "user",
				Message: "",
			},
		},
		{
			"user message",
			notif{
				Channel: "Query",
				User:    "user",
				Message: "message",
			},
		},
		{
			"user  ",
			notif{
				Channel: "Query",
				User:    "user",
				Message: " ",
			},
		},
	}

	for _, test := range tt {
		notif, err := Parse(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expext, notif)
	}
}
