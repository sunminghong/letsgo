package net

import (
	"os"
	"os/signal"
	//"syscall"
)

func (self *Server) signalSetup() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, os.Kill)
	<-ch
	//self.stop(nil, "SIGTERM")
}
