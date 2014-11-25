// +build !zmq2

package main

import "github.com/pebbe/zmq3"

func NewSubscriber(endpoint string) (Subscriber, error) {
	sub, err := zmq3.NewSocket(zmq3.SUB)
	if err != nil {
		return nil, err
	}
	return v3Subscriber{Socket: sub}, nil
}

type v3Subscriber struct {
	Socket   *zmq3.Socket
	endpoint string
}

func (v3b v3Subscriber) Connect() error {
	return v3b.Socket.Connect(v3b.endpoint)
}
func (v3b v3Subscriber) SetSubscribe(topic string) error {
	return v3b.Socket.SetSubscribe(topic)
}
func (v3b v3Subscriber) RecvMessage(flag int) (msg []string, err error) {
	return v3b.Socket.RecvMessage(zmq3.Flag(flag))
}
func (v3b v3Subscriber) RecvMessageBytes(flag int) (msg [][]byte, err error) {
	return v3b.Socket.RecvMessageBytes(zmq3.Flag(flag))
}
