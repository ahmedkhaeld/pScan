// Package scan provides a simple way to scan a network for hosts
package scan

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single TCP port
type PortState struct {
	Port int   // TCP port number
	Open state // indicates if the port is open or closed
}

type state bool // true if open, false if closed

type Results struct {
	Host       string
	PortStates []PortState
	NotFound   bool
}

// String returns a string representation of the state
func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

// Run scans a list of hosts for open ports
// and returns a slice of Results
func Run(hl *HostsList, ports []int) []Results {
	res := make([]Results, 0, len(hl.Hosts))

	for _, host := range hl.Hosts { // for each host
		r := Results{Host: host}                        // create a new result
		if _, err := net.LookupHost(host); err != nil { // verify the host exists
			r.NotFound = true    // if not, set the NotFound flag
			res = append(res, r) // and
			continue             // skip to the next host
		}
		for _, port := range ports { // for each port
			r.PortStates = append(r.PortStates, scanPort(host, port)) // scan the port
		}
		res = append(res, r) // add the result to the slice
	}
	return res // return the slice of results
}

// scanPort scans a single port on a host
// and returns the state of the port by provide the port number
func scanPort(host string, port int) PortState {
	p := PortState{Port: port}

	//verify the port state using net.DialTimeout

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		p.Open = false
		return p
	}
	conn.Close()
	p.Open = true
	return p

}
