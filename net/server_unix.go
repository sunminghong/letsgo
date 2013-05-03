package net

import (
	"os"
	"os/signal"
	"syscall"
)

func (self *TCPServer) signalSetup() {
	ch := make(chan os.Signal, 1)
        signal.Notify(ch, os.Kill)
	<-ch
	self.stop(nil, "SIGTERM")
}
