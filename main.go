package main

import (
	"fmt"
	"log"
	"net"
)

const buffSize = 1024
const channelCount = 10

func main() {
	go ingestPublic()
	ingestPrivate()
}

// handles incoming private data
func ingestPrivate() {
	ln, err := net.Listen("tcp", ":3095")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()
	fmt.Println("TCP on port 3095")

	// iterate, handling one private directive at a time
	buf := make([]byte, buffSize)
	for {
		// accept connection on port
		conn, _ := ln.Accept()
		defer conn.Close()

		for {
			// read until connection close
			if _, err := conn.Read(buf); err != nil {
				break
			}

			doPrivateOp(buf)
		}
	}
}

// handles all incoming public data
func ingestPublic() {
	sem := make(chan struct{}, channelCount)

	// listen to incoming udp packets
	sock, err := net.ListenPacket("udp", ":23498")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP on port 23498")
	defer sock.Close()

	// loop handling all client requests
	buf := make([]byte, buffSize)
	for {
		n, addr, err := sock.ReadFrom(buf)
		if err != nil {
			continue
		}

		sem <- struct{}{} // increment the semaphore (or freeze if limit hit)
		go doPublicOp(sem, sock, addr, buf[:n])
	}
}

// performs an operation from the public ingest
func doPublicOp(sem chan struct{}, pc net.PacketConn, addr net.Addr, buf []byte) {
	switch op := buf[0]; op {
	case 0:
		fmt.Println("disconnect operation")
	case 1:
		fmt.Println("echo operation")
	default:
		fmt.Println("Unknown op")
	}

	// free up a slot in the semaphore for the next packet
	<-sem
}

// performs an operation from the private ingest
func doPrivateOp(buf []byte) {
	switch op := buf[0]; op {
	case 1:
		fmt.Println("request room operation")
	case 2:
		fmt.Println("close room operation")
	case 3:
		fmt.Println("echo to room operation")
	case 4:
		fmt.Println("add user to room operation")
	case 5:
		fmt.Println("remove user from room operation")
	default:
		fmt.Println("Unknown private op")
	}
}
