// Package netcounts parses the /proc/net files to generate a Value
// structure based summary of network packet transfers.
package netcounts

import (
	"net"
	"os"
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
	Device map[string]*Snap
}

// Update refreshes the content of v.
func (v *Value) Update() error {
	b, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return err
	}
	i := func(n string) int64 {
		x, _ := strconv.ParseInt(n, 10, 64)
		return x
	}
	v.When = time.Now()
	for _, ifc := range strings.Split(string(b), "\n") {
		vs := strings.Fields(ifc)
		if len(vs) != 17 || vs[1] == "|bytes" {
			continue
		}
		dev := vs[0][:len(vs[0])-1]
		ifc, err := net.InterfaceByName(dev)
		ip4 := "unknown"
		ip6 := "unknown"
		if err == nil {
			as, err := ifc.Addrs()
			if err == nil {
				for _, a := range as {
					n, ok := a.(*net.IPNet)
					if !ok {
						continue
					}
					if n.IP.To4() != nil {
						ip4 = n.IP.String()
						continue
					}
					ip6 = n.IP.String()
				}
			}
		}
		v.Device[dev] = &Snap{
			IP:        ip4,
			IP6:       ip6,
			RxBytes:   i(vs[1]),
			RxPackets: i(vs[2]),
			TxBytes:   i(vs[9]),
			TxPackets: i(vs[10]),
		}
	}
	return nil
}

// NewValue returns a newly created value with a recent sample.
func NewValue() (*Value, error) {
	v := &Value{
		Device: make(map[string]*Snap),
	}
	return v, v.Update()
}
