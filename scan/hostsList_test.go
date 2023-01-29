package scan_test

import (
	"errors"
	"github.com/ahmedkhaeld/pScan/scan"
	"os"
	"testing"
)

func TestHostsList_Add(t *testing.T) {

	tests := []struct {
		name   string
		host   string
		expLen int
		expErr error
	}{
		{"AddNewHost", "host2", 2, nil},
		{"AddExistingHost", "host1", 1, scan.ErrExists},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			//initialize the list with a host
			if err := hl.Add("host1"); err != nil {
				t.Fatalf("failed to initialize the list: %v", err)
			}

			err := hl.Add(tc.host)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tc.expErr)
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %v, got %v", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}

			//check the length of the list
			if len(hl.Hosts) != tc.expLen {
				t.Errorf("expected length %d, got %d", tc.expLen, len(hl.Hosts))
			}

			//check the host is in the list
			if hl.Hosts[1] != tc.host {
				t.Errorf("expected %s, got %s", tc.host, hl.Hosts[1])
			}
		})
	}
}

func TestHostsList_Remove(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		expLen int
		expErr error
	}{
		{"RemoveExistingHost", "host1", 1, nil},
		{"RemoveNonExistingHost", "host3", 1, scan.ErrNotExists},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			hl := &scan.HostsList{}
			for _, host := range []string{"host1", "host2"} {
				if err := hl.Add(host); err != nil {
					t.Fatalf("failed to initialize the list: %v", err)
				}
			}

			err := hl.Remove(tc.host)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tc.expErr)
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected error %v, got %v", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}

			//check the length of the list
			if len(hl.Hosts) != tc.expLen {
				t.Errorf("expected length %d, got %d", tc.expLen, len(hl.Hosts))
			}

			//check the host is not in the list
			for _, host := range hl.Hosts {
				if host == tc.host {
					t.Errorf("expected %s to be removed, but it is still in the list", tc.host)
				}
			}
		})
	}

}

func TestHostsList_LoadSave(t *testing.T) {
	hl1 := &scan.HostsList{}
	hl2 := &scan.HostsList{}

	hostName := "host1"
	if err := hl1.Add(hostName); err != nil {
		t.Fatalf("failed to initialize the list: %v", err)
	}

	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tf.Name())

	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("failed to save the list: %v", err)
	}

	if err := hl2.Load(tf.Name()); err != nil {
		t.Fatalf("failed to load the list: %v", err)
	}

	//make sure the host is in the list
	if hl2.Hosts[0] != hostName {
		t.Errorf("expected %s, got %s", hostName, hl2.Hosts[0])
	}
	if len(hl2.Hosts) != 1 {
		t.Errorf("expected length %d, got %d", 1, len(hl2.Hosts))
	}
}

func TestHostsList_NoFile(t *testing.T) {
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if err := os.Remove(tf.Name()); err != nil {
		t.Fatalf("failed to remove temp file: %v", err)
	}

	hl := &scan.HostsList{}
	if err := hl.Load(tf.Name()); err != nil {
		t.Errorf("expected error, got nil")
	}
}
