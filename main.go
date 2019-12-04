package main

import (
	"bufio"
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
	ln, _ := net.Listen("tcp", ":8081")

	// iterate, handling one private directive at a time
	for {
		// accept connection on port
		conn, _ := ln.Accept()

		// read until newline
		op, _ := bufio.NewReader(conn).ReadString('\n')
		doPrivateOp(op)
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

	defer sock.Close()

	// loop handling all client requests
	for {
		buf := make([]byte, buffSize)
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
	// TODO: parse op code and perform action
	// possible public operations:
	// - send to room
	// - disconnect

	// free up a slot in the semaphore for the next packet
	<-sem
}

// performs an operation from the private ingest
func doPrivateOp(op string) {
	// TODO: parse op code and perform action
	// possible private operations:
	// - request room
	// - close room
	// - send to room
	// - add user
	// - evict user
}
