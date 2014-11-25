package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/vbatts/flaming-happiness/common"
)

func main() {
	flag.Parse()

	if len(ignoreChannels) > 0 {
		common.SetIgnores(ignoreChannels)
	}

	if len(flag.Args()) != 1 {
		log.Fatalf("provide the zmq_notify publisher! like tcp://example.com:2428")
	}

	subscriber, err := NewSubscriber(flag.Args()[0])
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if err = subscriber.Connect(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	if !quiet {
		log.Printf("Connected to [%s]", flag.Args()[0])
	}
	if err = subscriber.SetSubscribe(""); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	for {
		msg, err := subscriber.RecvMessageBytes(0)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
			break
		}
		noti_msg := common.IrcNotify{}
		json.Unmarshal(msg[0], &noti_msg)
		go common.Display(noti_msg, linger, quiet)
	}
}

type Subscriber interface {
	Connect() error
	SetSubscribe(string) error
	RecvMessage(flag int) (msg []string, err error)
	RecvMessageBytes(flag int) (msg [][]byte, err error)
}

var (
	linger         int64 = 5
	quiet          bool  = false
	ignoreChannels string
)

func init() {
	flag.Int64Var(&linger, "linger",
		linger, "time to let the notification linger")
	flag.BoolVar(&quiet, "quiet",
		false, "less output")
	flag.StringVar(&ignoreChannels, "ignore",
		"", "comma seperated list of pattern of channels to ignore")
}
