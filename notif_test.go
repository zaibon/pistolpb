package main

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tt := []struct {
		input  string
		expect notif
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
		assert.Equal(t, test.expect, notif)
	}
}

func TestSameLessThan(t *testing.T) {
	pistol := &Pistol{
		lastQueries: map[string]time.Time{},
	}
	tt := []struct {
		input  notif
		expect bool
	}{
		{
			input: notif{
				Channel: queryType,
				User:    "user1",
			},
			expect: false,
		},
		{
			input: notif{
				Channel: queryType,
				User:    "user1",
			},
			expect: true,
		},
		{
			input: notif{
				Channel: queryType,
				User:    "user2",
			},
			expect: false,
		},
		{
			input: notif{
				Channel: "normalChan",
				User:    "user1",
			},
			expect: false,
		},
	}

	for i := 0; i < len(tt); i++ {
		test := tt[i]
		now := time.Now()
		res := pistol.sameLessThan(test.input, time.Second*5)
		if !assert.Equal(t, res, test.expect) {
			t.Log("now :", now)
			t.Log("input :", test.input)
			t.Log("lastQueries :", pistol.lastQueries)
		}
	}

}
