package main

import (
	"bytes"
	"encoding/json"
	"flag"
	zmq "github.com/pebbe/zmq2"
	"github.com/vbatts/flaming-happiness/common"
	"log"
)

func main() {
	flag.Parse()

  if len(ignoreChannels) > 0 {
    common.SetIgnores(ignoreChannels)
  }

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
		noti_msg := common.IrcNotify{}
		json.Unmarshal(bytes.NewBufferString(msg[0]).Bytes(), &noti_msg)
		go common.Display(noti_msg, linger, quiet)
	}
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
