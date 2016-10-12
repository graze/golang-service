// Network helpers for when unit testing networks
//
// Uses a channel log the messages recieved by the server
//
// Usage:
//     done := make(chan string)
//     addr, sock, srvWg := CreateServer(t, tc.net, tc.la, done)
//     defer srvWg.Wait()
//     defer os.Remove(addr.String())
//     defer sock.Close()
//
//     s, err := net.Dial(tc.net, addr.String())
//     defer s.Close()
//     fmt.Fprintf(s, "test message\n")
//     if "test message\n" != <-done {
//         t.Error("message not recieved")
//     }
package nettest
