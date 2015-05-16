package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xconstruct/go-pushbullet"
)

var help = `Usage:
    pistol command [arguments]


Commands:
    login -k key : register pushbullet api key
    run -f fnotify.txt : start watching fnotify file`

var (
	key     = flag.String("k", "", "register pushbullet api key")
	fnotify = flag.String("f", "", "fnotify file to watch")
)

func main() {
	flag.Parse()

	if *key != "" {
		login(*key)
	}
	if *fnotify != "" {
		run(*fnotify)
	}
}

func login(key string) error {
	cfg := config{
		Apikey:  key,
		Devices: chooseDevice(key),
	}
	err := writeCfg(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("pistolPB configured")
	return err
}

func run(fnotify string) {
	cfg, err := readCfg()
	if err != nil {
		log.Fatalln("error reading config:", err)
	}

	f, err := os.Open(fnotify)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	f.Seek(0, os.SEEK_END)

	pistol, err := newPistol(f, cfg)
	if err != nil {
		log.Fatalln("error conntion to pushbuller:", err)
	}

	// handle interruption
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		os.Exit(0)
	}()

	fmt.Printf("watching %s\n", f.Name())
	pistol.Watch()
}

func chooseDevice(key string) []*pushbullet.Device {
	cl := pushbullet.New(key)
	devs, err := cl.Devices()
	if err != nil {
		log.Fatalln("error retreiving devices:", err)
	}

	for i := 0; i < len(devs); i++ {
		if devs[i].Active && devs[i].Pushable {
			fmt.Printf("%d : %s\n", i, devs[i].Nickname)
		}
	}
	fmt.Print("choose on which device send notification (ex: 0,1,3) : ")
	var choises string
	_, err = fmt.Scanf("%s", &choises)
	if err != nil {
		log.Fatalln(err)
	}

	num := strings.Split(choises, ",")
	idens := []*pushbullet.Device{}
	for _, n := range num {
		i, err := strconv.Atoi(n)
		if err != nil {
			log.Println(err)
			continue
		}
		idens = append(idens, devs[i])
	}
	return idens
}

func writeCfg(cfg config) error {
	cfgfile := filepath.Join(os.Getenv("HOME"), ".pistolPB2")
	f, err := os.OpenFile(cfgfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	return json.NewEncoder(f).Encode(cfg)
}

func readCfg() (*config, error) {
	cfgfile := filepath.Join(os.Getenv("HOME"), ".pistolPB2")
	f, err := os.Open(cfgfile)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	cfg := &config{}
	err = json.NewDecoder(f).Decode(cfg)
	return cfg, err
}
