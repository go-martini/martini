package main

import (
	"github.com/judwhite/go-svc/svc"
	"github.com/liyue201/martini"
	"log"
	"sync"
	"syscall"
)

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
	log.Println("martini was stopped gracefully")
}

type program struct {
	m         *martini.ClassicMartini
	waitGroup WaitGroupWrapper
}

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

func (p *program) Init(env svc.Environment) error {
	m := martini.Classic()

	m.Get("/", func() string {
		return "Hello world!"
	})

	p.m = m
	return nil
}

func (p *program) Start() error {

	p.waitGroup.Wrap(func() {
		p.m.RunOnAddr(":80")
	})

	p.waitGroup.Wrap(func() {
		p.m.RunOnAddrTLS(":443", "server.crt", "server.key")
	})

	return nil
}

func (p *program) Stop() error {
	p.m.Stop()
	p.waitGroup.Wait()
	return nil
}
