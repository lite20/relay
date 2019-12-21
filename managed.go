package main

import (
	"net"
)

var pools map[int][]net.Addr

func poolCreate() int {
	id := len(pools)
	pools[id] = make([]net.Addr, 0, 1)

	return id
}

// assumes pool id is a valid id (exists)
func poolDestroy(poolID int) {
	delete(pools, poolID)
}

// assumes pool id is a valid id (exists)
// assumes address is a valid address
func poolUserAdd(poolID int, addr net.Addr) {
	pools[poolID] = append(pools[poolID], addr)
}

// assumes pool id is a valid id (exists)
func poolUserRemove(poolID int, addr net.Addr) {
	pool := pools[poolID]
	poolEnd := len(pool) - 1

	// check if address is in pool
	for i, pAddr := range pools[poolID] {
		if addr == pAddr {
			// put the value at the end in the place of the value we want replaced
			pools[poolID][i] = pools[poolID][poolEnd]

			// return a slice shrunk by one (removing the last value)
			pools[poolID] = pools[poolID][:poolEnd]
		}
	}
}

// assumes valid pool
func poolWriteTo(poolID int, sock net.PacketConn, buf []byte) {
	for _, addr := range pools[poolID] {
		sock.WriteTo(buf, addr)
	}
}
