// Package scan provide types and function to perform TCP port
// scan on a list of hosts
package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("host already exists")
	ErrNotExists = errors.New("host does not exist")
)

// HostsList is a list of hosts to run port scan
type HostsList struct {
	Hosts []string
}

// search searches for a host in sorted list of hosts
func (hl *HostsList) search(host string) (bool, int) {
	sort.Strings(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}
	return false, i
}

// Add adds a host to the list of hosts to scan
func (hl *HostsList) Add(host string) error {
	if exists, _ := hl.search(host); exists {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}

	hl.Hosts = append(hl.Hosts, host)
	return nil
}

// Remove deletes a host from the list of hosts to scan
func (hl *HostsList) Remove(host string) error {
	if exists, i := hl.search(host); exists {
		hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
		return nil
	}

	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// Load loads the list of hosts from a file
func (hl *HostsList) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // file does not exist
		}
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}
	return scanner.Err()
}

// Save saves the list of hosts to a file
func (hl *HostsList) Save(hostsFile string) error {
	output := ""
	for _, host := range hl.Hosts {
		output += fmt.Sprintln(host)
	}

	return os.WriteFile(hostsFile, []byte(output), 0644)
}
