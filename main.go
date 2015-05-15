package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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
	err := writeCfg(key)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("API key saved")
	return err
}

func run(fnotify string) {
	apikey, err := readCfg()
	if err != nil {
		log.Fatalln("error reading config:", err)
	}

	f, err := os.Open(fnotify)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	f.Seek(0, os.SEEK_END)

	fmt.Println(apikey)
	pistol, err := newPistol(f, apikey)
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

func writeCfg(key string) error {
	cfgfile := filepath.Join(os.Getenv("HOME"), ".pistolPB")
	f, err := os.OpenFile(cfgfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = f.WriteString(key)
	return err
}

func readCfg() (string, error) {
	cfgfile := filepath.Join(os.Getenv("HOME"), ".pistolPB")
	f, err := os.Open(cfgfile)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	buff := &bytes.Buffer{}
	_, err = io.Copy(buff, f)
	return buff.String(), err
}
