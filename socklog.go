// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"net"
	"os"
)

// This log writer sends output to a socket
type SocketLogWriter chan *LogRecord

// This is the SocketLogWriter's output method
func (w SocketLogWriter) LogWrite(rec *LogRecord) {
	w <- rec
}

func (w SocketLogWriter) Close() {
	close(w)
}

func NewSocketLogWriter(proto, hostport string) SocketLogWriter {
	sock, err := net.Dial(proto, hostport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewSocketLogWriter(%q): %s\n", hostport, err)
		return nil
	}

	w := SocketLogWriter(make(chan *LogRecord, LogBufferLength))

	go func() {
		defer func() {
			if sock != nil && proto == "tcp" {
				sock.Close()
			}
		}()

		for rec := range w {
			_, err := sock.Write([]byte(rec.Message))
			if err != nil {
				fmt.Fprintf(os.Stderr, "SockLogWriteErr:%s, # %s\n", err.Error(), rec.Message)
				sock, err = net.Dial(proto, hostport)
			}
		}
	}()

	return w
}
