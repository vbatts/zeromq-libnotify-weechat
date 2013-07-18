package main

import (
	"bytes"
	"encoding/json"
	"flag"
	notify "github.com/mqu/go-notify"
	zmq "github.com/pebbe/zmq3"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	flag.Parse()

	subscriber, _ := zmq.NewSocket(zmq.SUB)
	if len(flag.Args()) == 1 {
		subscriber.Connect(flag.Args()[0])
		if !quiet {
			log.Printf("Connected to [%s]", flag.Args()[0])
		}
	} else {
		log.Fatalf("provide the zmq_notify publisher! like tcp://example.com:2428")
	}
	subscriber.SetSubscribe("")

	for {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			break
		}
		go func() {
			noti_msg := IrcNotify{}
			json.Unmarshal(bytes.NewBufferString(msg[0]).Bytes(), &noti_msg)
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
			time.Sleep(time.Duration(delay) * time.Second)
			notify.NotificationClose(hello)
		}()
	}
}

var (
	delay          int64 = 5
	quiet          bool  = false
	ignoreChannels string
)

func ShouldIgnore(channel string) bool {
	for _, ignore := range strings.Split(ignoreChannels, ",") {
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

	flag.Int64Var(&delay, "delay",
		delay, "time to let the notification linger")
	flag.BoolVar(&quiet, "quiet",
		false, "less output")
	flag.StringVar(&ignoreChannels, "ignore",
		"", "comma seperated list of pattern of channels to ignore")

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
