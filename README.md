flaming-happiness
=================

Golang + ZeroMQ + libnotify tool to connect to a zmq_notify.rb weechat script

This is ideal for situatishes like screen + weechat is running on a faroff
remote host, so a notification plugin for weechat would have no good access to
DISPLAY on the local host that you are connecting from.
Further, since it is pub/sub, there is no limit to only having a single client 
being notified.


Install
-------

if you have ZeroMQ v3 installed, run:

	go get github.com/vbatts/flaming-happiness/noti

for ZeroMQ v2

	go get github.com/vbatts/flaming-happiness/noti2


Running
-------

First you'll need to add the zmq_notify.rb file to ~/.weechat/ruby/
and load it from weechat (you can symlink it in the ./autoload/ directory
to have this script loaded when weechat launches)

	/ruby load zmq_notify.rb

Then run `noti` against your site

	noti tcp://example.com:2428

Be sure you've allowed for any firewalling between noti and the zmq_notify

Thanks
------

To folks that had already supplied the pieces of a simpler enabler like this
and to github for the wonky repository name recommendations.

