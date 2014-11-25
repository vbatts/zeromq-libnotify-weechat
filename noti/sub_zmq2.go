// +build zmq2

package main

import "github.com/pebbe/zmq2"

func NewSubscriber(endpoint string) (Subscriber, error) {
	sub, err := zmq2.NewSocket(zmq2.SUB)
	if err != nil {
		return nil, err
	}
	return v2Subscriber{Socket: sub}, nil
}

type v2Subscriber struct {
	Socket   *zmq2.Socket
	endpoint string
}

func (v2b v2Subscriber) Connect() error {
	return v2b.Socket.Connect(v2b.endpoint)
}
func (v2b v2Subscriber) SetSubscribe(topic string) error {
	return v2b.Socket.SetSubscribe(topic)
}
func (v2b v2Subscriber) RecvMessage(flag int) (msg []string, err error) {
	return v2b.Socket.RecvMessage(zmq2.Flag(flag))
}
func (v2b v2Subscriber) RecvMessageBytes(flag int) (msg [][]byte, err error) {
	return v2b.Socket.RecvMessageBytes(zmq2.Flag(flag))
}
