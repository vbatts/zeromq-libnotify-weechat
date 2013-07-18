package main

import (
	"bytes"
	"encoding/json"
	notify "github.com/mqu/go-notify"
	zmq "github.com/pebbe/zmq3"
	"log"
	"os"
	"os/signal"
	"time"
)

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

const (
	DELAY = 3000
)

func main() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	if len(os.Args) == 2 {
		subscriber.Connect(os.Args[1])
    log.Printf("Connected to [%s]", os.Args[1])
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
			log.Printf("%#v", noti_msg)
			hello := notify.NotificationNew(noti_msg.Server + "," + noti_msg.Channel,
      noti_msg.Message,
      "")

			if hello == nil {
				log.Println("ERROR: Unable to create a new notification")
				return
			}
			notify.NotificationSetTimeout(hello, DELAY)

			if e := notify.NotificationShow(hello); e != nil && len(e.Message()) > 0 {
				log.Printf("ERROR: %s", e.Message())
				return
			}
			time.Sleep(DELAY * time.Second)
			notify.NotificationClose(hello)
		}()
	}
}
