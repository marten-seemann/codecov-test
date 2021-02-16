package main

import (
	// "crypto/rand"
	"log"
	"net"
	"time"

	"github.com/libp2p/go-yamux"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}
	go func() {
		if err := runServer(ln); err != nil {
			panic(err)
		}
	}()

	if err := runClient(ln.Addr().String()); err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
}

func runServer(ln net.Listener) error {
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	sess, err := yamux.Server(conn, yamux.DefaultConfig())
	if err != nil {
		return err
	}
	log.Println("Accepted Yamux session")
	var counter int
	for {
		_, err := sess.AcceptStream()
		if err != nil {
			return err
		}
		log.Printf("\tAccepted stream %d\n", counter)
		time.Sleep(time.Millisecond)
		counter++
	}
	// str, err := sess.AcceptStream()
	// if err != nil {
	// 	return err
	// }

	// for {
	// 	b := make([]byte, 6)
	// 	n, err := str.Read(b)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	b = b[:n]
	// 	// log.Printf("\t\tRead: %#x. Num streams: %d\n", b, sess.NumStreams())
	// }
	return nil
}

func runClient(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	conf := yamux.DefaultConfig()
	conf.AcceptBacklog = 1000
	sess, err := yamux.Client(conn, conf)
	if err != nil {
		return err
	}
	log.Println("Established Yamux session.")
	for i := 0; i < 1000; i++ {
		str, err := sess.OpenStream()
		if err != nil {
			return err
		}
		log.Printf("\tOpened stream %d\n", i)
		go func(i int) {
			b := make([]byte, 6)
			str.Write(b)
			// str.Close()
			// _, err = str.Read(b)
			// log.Printf("Stream %d write err: %s\n", i, err)
		}(i)
	}
	time.Sleep(time.Second)
	_, err = sess.OpenStream()
	if err != nil {
		return err
	}
	log.Printf("\tOpened stream %d\n", 1001)
	return nil
}
