package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

// IsPrivateIP ...
func IsPrivateIP(IP net.IP) bool {

	for _, in := range PrivateIPNets {
		if in.Contains(IP) {
			return true
		}
	}

	return false
}

// GetIPAddrs ...
func GetIPAddrs() (private, public []string) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, address := range addrs {
		in, ok := address.(*net.IPNet)
		if !ok || in.IP.IsLoopback() || in.IP.To4() == nil {
			continue
		}

		if IsPrivateIP(in.IP) {
			private = append(private, in.IP.String())
		} else {
			public = append(public, in.IP.String())
		}
	}

	return
}

// ParseIPFromUint32 ...
func ParseIPFromUint32(n uint32) string {

	return fmt.Sprintf("%d.%d.%d.%d", byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

// ParseIPToUint32 ...
func ParseIPToUint32(ip string) uint32 {

	if v := net.ParseIP(ip); v != nil {
		return binary.BigEndian.Uint32(v.To4())
	}

	return 0
}

// ParseCIDRToStartEnd ...
func ParseCIDRToStartEnd(CIDR string) (uint32, uint32) {
	items := strings.Split(strings.TrimSpace(CIDR), "/")

	mask, _ := strconv.ParseUint(items[1], 10, 32)
	start := ParseIPToUint32(items[0]) & (math.MaxUint32 << (32 - mask))
	end := start + (1<<(32-mask) - 1)

	return start, end
}

// ParseIPRangeToCIDRs ...
func ParseIPRangeToCIDRs(start, end string) (CIDRs []string) {

	s := ParseIPToUint32(start)
	e := ParseIPToUint32(end)

	for s <= e {

		var i int64
		for ; i < 32; i++ {
			if s+(1<<i)-1 > e {
				i--
				break
			}

			if s&(1<<i) == (1 << i) {
				break
			}
		}

		CIDRs = append(CIDRs, fmt.Sprintf("%s/%d", ParseIPFromUint32(uint32(s)), 32-i))

		s = (1 << i) + s
	}

	return
}
