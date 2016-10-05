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
    "io/ioutil"
    "sync"
    "time"
    "os"
)

func CreateServer(t *testing.T, n, la string, done chan<- string) (addr net.Addr, sock io.Closer, wg *sync.WaitGroup) {
    if n == "udp" || n == "tcp" {
        la = "127.0.0.1:0"
    } else {
        // unix and unixgram: choose an address if none given
        if la == "" {
            // use ioutil.TempFile to get a name that is unique
            f, err := ioutil.TempFile("", "servertest")
            if err != nil {
                t.Fatalf("TempFile: %v", err)
            }
            f.Close()
            la = f.Name()
        }
        os.Remove(la)
    }

    wg = new(sync.WaitGroup)
    if n == "udp" || n == "unixgram" {
        l, err := net.ListenPacket(n, la)
        if err != nil {
            t.Fatalf("CreateServer failed: %v", err)
        }
        addr = l.LocalAddr()
        sock = l
        wg.Add(1)
        go func() {
            defer wg.Done()
            readPackets(l, done)
        }()
    } else {
        l, err := net.Listen(n, la)
        if err != nil {
            t.Fatalf("CreateServer failed: %v", err)
        }

        addr = l.Addr()
        sock = l
        wg.Add(1)

        go func() {
            defer wg.Done()
            readStream(l, done, wg)
        }()
    }
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

func readPackets(c net.PacketConn, done chan<- string) {
    var buf [4096]byte
    var rcvd string
    ct := 0
    for {
        var n int
        var err error

        c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
        n, _, err = c.ReadFrom(buf[:])
        rcvd += string(buf[:n])
        if err != nil {
            if oe, ok := err.(*net.OpError); ok {
            if ct < 3 && oe.Temporary() {
                ct++
                continue
                }
            }
        break
        }
    }
    c.Close()
    done <- rcvd
}
