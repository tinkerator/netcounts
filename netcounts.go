// Package netcounts executes `ifconfig` and parses the output to
// generate a Value structure based summary.
package netcounts

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Snap struct {
	IP, IP6            string
	RxPackets, RxBytes int64
	TxPackets, TxBytes int64
}

// Value holds a snapshot of the most recent network device counts.
type Value struct {
	When   time.Time
	Device map[string]Snap
}

// cmd executes a command and returns the stdout in a bytes.Buffer.
func cmd(cmd string, args ...string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	c := exec.Command(cmd, args...)
	c.Stdout = buf
	err := c.Run()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

var extractorRE = regexp.MustCompile(`^([^:]+):.+inet\s+([0-9\.]+)\s.+inet6\s+([a-f0-9:]+)\s.+RX\s+packets\s+([0-9]+)\s+bytes\s+([0-9]+)\s.+TX\s+packets\s+([0-9]+)\s+bytes\s+([0-9]+)\s`)

// File location of the ifconfig binary.
var IfconfigBinary = "/usr/sbin/ifconfig"

// Update refreshes the content of v.
func (v *Value) Update() error {
	b, err := cmd(IfconfigBinary)
	if err != nil {
		return err
	}
	i := func(n string) int64 {
		x, _ := strconv.ParseInt(n, 10, 64)
		return x
	}
	v.When = time.Now()
	for _, ifc := range strings.Split(b.String(), "\n\n") {
		if ifc == "" {
			break
		}
		ifc = strings.ReplaceAll(ifc, "\n", " ")
		vs := extractorRE.FindAllStringSubmatch(ifc, 1)
		if len(vs) != 1 || len(vs[0]) != 8 {
			continue
		}
		d := vs[0][1:]
		v.Device[d[0]] = Snap{
			IP:        d[1],
			IP6:       d[2],
			RxPackets: i(d[3]),
			RxBytes:   i(d[4]),
			TxPackets: i(d[5]),
			TxBytes:   i(d[6]),
		}
	}
	return nil
}

// NewValue returns a newly created value with a recent sample.
func NewValue() (*Value, error) {
	v := &Value{
		Device: make(map[string]Snap),
	}
	return v, v.Update()
}
