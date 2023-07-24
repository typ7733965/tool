package tool

import (
	"net"
	"sync"
)

var lock = sync.RWMutex{}
var ip string = ""

func GetIp() string {
	lock.RLock()
	defer lock.RUnlock()
	if ip != "" {
		return ip
	}
	ip, _ = GetLocalIPV4()
	return ip
}
func GetLocalIPV4() (ipv4 string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.To4().String(), nil
			}
		}
	}
	return "", err
}
