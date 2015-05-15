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
	notifile io.Reader
	Client   *pushbullet.Client
	Devices  []*pushbullet.Device
}

func newPistol(r io.Reader, apiKey string) (*Pistol, error) {
	p := &Pistol{
		notifile: r,
		Client:   pushbullet.New(apiKey),
	}

	devices, err := p.Client.Devices()
	if err != nil {
		return nil, err
	}
	for _, d := range devices {
		if d.Active {
			p.Devices = append(p.Devices, d)
		}
	}

	return p, nil
}

func (p *Pistol) Watch() {
	for notif := range readnotif(p.notifile) {
		p.broadcast(notif)
		log.Printf("send %+v\n", notif)
	}
}

func (p *Pistol) broadcast(notif notif) {
	for _, d := range p.Devices {
		msg := fmt.Sprintf("Channe:%s\nFrom : %s\n%s", notif.Channel, notif.User, notif.Message)
		p.Client.PushNote(d.Iden, fmt.Sprintf("IRC: %s", notif.Channel), msg)
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
				// log.Println(err)
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
		n.Channel = "Query"
		i := strings.Index(s, " ")
		n.User = s[:i]
		n.Message = s[i+1:]
	}
	return n, nil
}
