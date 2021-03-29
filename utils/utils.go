package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net"
	"time"
)

var (
	chars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// PrivateCIDRs ...
	PrivateCIDRs = []string{
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	// PrivateIPNets ...
	PrivateIPNets []*net.IPNet
)

func init() {

	rand.Seed(time.Now().UnixNano())

	for _, CIDR := range PrivateCIDRs {
		_, in, err := net.ParseCIDR(CIDR)
		if err != nil {
			continue
		}
		PrivateIPNets = append(PrivateIPNets, in)
	}
}

// Md5sum ...
func Md5sum(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s)) // nolint: errcheck
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetUnixMesc ...
func GetUnixMesc() int64 {
	return time.Now().UnixNano() / 1e6
}

// GenRandStr ...
func GenRandStr(length int) string {

	maxrb := 255 - (256 % 62)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return ""
		}

		for _, rb := range r {

			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}

			b[i] = chars[c%62]
			i++

			if i == length {
				return string(b)
			}
		}
	}
}
