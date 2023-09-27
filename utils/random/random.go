package random

import (
	"time"
)

var (
	runes  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	Tunnel chan string
)

func init() {
	Tunnel = make(chan string, 1)
	go func() {
		for {
			b := make([]rune, 20)
			for i := range b {
				index := time.Now().UnixNano() % int64(len(runes))
				b[i] = runes[index]
			}
			Tunnel <- string(b)
		}
	}()
}
