// This file is part of graze/golang-service
//
// Copyright (c) 2016 Nature Delivered Ltd. <https://www.graze.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.
//
// @license https://github.com/graze/golang-service-logging/blob/master/LICENSE
// @link    https://github.com/graze/golang-service-logging

package testing

import (
    "net"
    "testing"
    "bufio"
    "io"
    "sync"
    "time"
)

func CreateTCPServer(t *testing.T, la string, done chan<- string) (addr net.Addr, sock io.Closer, wg *sync.WaitGroup) {
    wg = new(sync.WaitGroup)

	l, e := net.Listen("tcp", la)
	if e != nil {
		t.Fatalf("CreateTCPServer failed: %v", e)
	}

    addr = l.Addr()
	sock = l
	wg.Add(1)

	go func() {
		defer wg.Done()
		readStream(l, done, wg)
	}()
    return
}

func readStream(l net.Listener, done chan<- string, wg *sync.WaitGroup) {
    for {
        var c net.Conn
        var err error
        if c, err = l.Accept(); err != nil {
            return
        }
        wg.Add(1)
        go func(c net.Conn) {
            defer wg.Done()
            c.SetReadDeadline(time.Now().Add(1 * time.Second))
            b := bufio.NewReader(c)
            for ct := 1; ct&7 != 0; ct++ {
                s, err := b.ReadString('\n')
                if err != nil {
                    break
                }
                done <- s
            }
            c.Close()
        }(c)
    }
}
