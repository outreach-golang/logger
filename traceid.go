package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func NewTraceID() string {
	var tid string

	tid = fmt.Sprintf("%s%d%d%d", ipTo16(), currentTime(), incId(), processId())

	return tid
}

func ipTo16() string {
	var (
		ip  = LocalIP()
		ips = strings.Split(ip, ".")
		rs  string
	)

	for _, v := range ips {
		rs += fmt.Sprintf("%x", v)
	}

	return rs

}

func currentTime() int64 {
	return time.Now().Unix()
}

var (
	id   = 1000
	lock sync.Mutex
)

func incId() int {
	defer lock.Unlock()

	lock.Lock()

	if id == 9999 {
		id = 1000
		return id
	}

	id += 1

	return id
}

func processId() int {
	return os.Getppid()
}
