package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/xconstruct/go-pushbullet"

	"io"
)

type Pistol struct {
	notifile    io.Reader
	client      *pushbullet.Client
	Devices     []*pushbullet.Device
	lastQueries map[string]time.Time
}

func newPistol(r io.Reader, cfg *config) (*Pistol, error) {
	p := &Pistol{
		notifile:    r,
		client:      pushbullet.New(cfg.Apikey),
		Devices:     cfg.Devices,
		lastQueries: map[string]time.Time{},
	}

	return p, nil
}

func (p *Pistol) Watch() {
	for notif := range readnotif(p.notifile) {

		// don't spam notification when queries come from same person
		if p.sameLessThan(notif, time.Minute*15) {
			continue
		}

		p.broadcast(notif)
		log.Printf("send %+v\n", notif)
	}
}

// sameLessThan tests if the notif is a query that come from user who sends a previous
// message less than duration d ago
func (p *Pistol) sameLessThan(notif notif, d time.Duration) bool {
	if notif.Channel == queryType {
		last, present := p.lastQueries[notif.User]
		now := time.Now()
		p.lastQueries[notif.User] = now
		if present && now.Sub(last) < d {
			return true
		}
	}
	return false
}

func (p *Pistol) broadcast(notif notif) {
	for _, d := range p.Devices {
		fmt.Println(d.Nickname)
		msg := fmt.Sprintf("Channel:%s\nFrom : %s\n%s", notif.Channel, notif.User, notif.Message)
		p.client.PushNote(d.Iden, fmt.Sprintf("IRC: %s", notif.Channel), msg)
	}
}

func readnotif(r io.Reader) chan notif {
	c := make(chan notif)
	br := bufio.NewReader(r)

	go func() {
		tick := time.NewTicker(time.Second)
		var notif = notif{}
		for range tick.C {
			l, err := br.ReadString('\n')
			if err != nil {
				// log.Println(err)
				continue
			}
			notif, err = Parse(l)
			if err != nil {
				log.Printf("error parsing line (%v) : %v\n", l, err)
				continue
			}
			c <- notif
		}
		close(c)
	}()

	return c
}

type notif struct {
	Channel string
	User    string
	Message string
}

var regexMsg = regexp.MustCompile("#(.+) <(.+)> (.+)")
var regexQuery = regexp.MustCompile("(.+) (.+)")

const queryType = "Query"

func Parse(s string) (notif, error) {
	n := notif{}
	if s == "" {
		return n, fmt.Errorf("parded ligne is empty")
	}

	if strings.HasPrefix(s, "#") {
		// send from a channel
		i := strings.Index(s, " ")
		n.Channel = s[:i]
		s = s[i+2:]
		i = strings.Index(s, ">")
		n.User = s[:i]
		n.Message = s[i+2:]
	} else {
		n.Channel = queryType
		i := strings.Index(s, " ")
		n.User = s[:i]
		n.Message = s[i+1:]
	}
	return n, nil
}

type config struct {
	Apikey  string
	Devices []*pushbullet.Device
}
