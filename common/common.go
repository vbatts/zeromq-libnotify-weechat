package common

import (
	notify "github.com/mqu/go-notify"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	ignoreList = []string{}
)

func SetIgnores(ignores string) {
	ignoreList = strings.Split(ignores, ",")
}

func Display(noti_msg IrcNotify, linger int64, quiet bool) {
	if !quiet {
		log.Printf("%#v", noti_msg)
	}
	if ShouldIgnore(noti_msg.Channel) {
		if !quiet {
			log.Printf("Ignoring: %s", noti_msg.Channel)
		}
		return
	}
	hello := notify.NotificationNew(noti_msg.Server+","+noti_msg.Channel,
		noti_msg.Message,
		"")

	if hello == nil {
		log.Println("ERROR: Unable to create a new notification")
		return
	}
	notify.NotificationSetTimeout(hello, 0)

	if e := notify.NotificationShow(hello); e != nil && len(e.Message()) > 0 {
		log.Printf("ERROR: %s", e.Message())
		return
	}
	time.Sleep(time.Duration(linger) * time.Second)
	notify.NotificationClose(hello)
}

func ShouldIgnore(channel string) bool {
	for _, ignore := range ignoreList {
		if strings.Contains(channel, ignore) {
			return true
		}
	}
	return false
}

type IrcNotify struct {
	Highlight bool     `json:"highlight"`
	Type      string   `json:"type"`
	Channel   string   `json:"channel"`
	Message   string   `json:"message"`
	Server    string   `json:"server"`
	Date      string   `json:"date"`
	Tags      []string `json:"tags"`
}

func init() {
	notify.Init("IRC-noti")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			notify.UnInit()
			os.Exit(0)
		}
	}()
}
