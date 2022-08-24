package tools

import (
	"math/rand"
	"net"
	"time"
)

// RandomIPv4 不完全合法，包含了 0/24 127/8 等地址，可能造成被WAF拦截
func RandomIPv4() string {
	rand.Seed(time.Now().Unix())
	return net.IPv4(byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255))).String()
}

// RandomIPv6 不确定合法程度，没有细看相关私有地址等内容，可能造成被WAF拦截
func RandomIPv6() string {
	rand.Seed(time.Now().Unix())
	ip := make([]byte, 16)
	for i := 0; i < 16; i++ {
		ip[i] = byte(rand.Intn(256))
	}
	return net.IP(ip).To16().String()
}
