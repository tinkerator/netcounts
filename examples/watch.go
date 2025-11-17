// Program watch pays attention to the rate of network activity. By
// default, it ignores the "lo" device, but use --ignore="" to reveal
// those stats.
//
// This is a simple cli to demonstrate using the
// [zappem.net/pub/net/netcounts] package.
//
// [zappem.net/pub/net/netcounts]: https://zappem.net/pub/net/netcounts/
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"zappem.net/pub/net/netcounts"
)

var (
	debug  = flag.Bool("debug", false, "log extra details")
	poll   = flag.Duration("poll", 5*time.Second, "polling interval")
	ignore = flag.String("ignore", "lo", "comma separated list of devices to ignore")
)

func main() {
	flag.Parse()

	skip := make(map[string]bool)
	for _, d := range strings.Split(*ignore, ",") {
		if d == "" {
			continue
		}
		skip[d] = true
	}
	var vs [2]*netcounts.Value
	for i := 0; ; i++ {
		var err error
		if i < 2 {
			vs[i], err = netcounts.NewValue()
		} else {
			err = vs[i&1].Update()
		}
		if err != nil {
			log.Fatalf("data problem: %v", err)
		}
		v := vs[i&1]
		if *debug {
			log.Printf("%7d: %v", i, v.When)
		}
		for k, s := range v.Device {
			if skip[k] {
				continue
			}
			var rP, rB, tP, tB string
			if i > 0 {
				if os, ok := vs[1-(i&1)].Device[k]; ok {
					rP = fmt.Sprint(s.RxPackets - os.RxPackets)
					rB = fmt.Sprint(s.RxBytes - os.RxBytes)
					tP = fmt.Sprint(s.TxPackets - os.TxPackets)
					tB = fmt.Sprint(s.TxBytes - os.TxBytes)
				}
			}
			log.Printf("  %q[%s, %s]: RX: %d/%d (%8s/%-8s); TX: %d/%d (%8s/%-8s)", k, s.IP, s.IP6, s.RxPackets, s.RxBytes, rP, rB, s.TxPackets, s.TxBytes, tP, tB)
		}
		time.Sleep(*poll)
	}
}
