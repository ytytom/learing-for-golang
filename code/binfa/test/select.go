package main

import "time"
import "fmt"

func main() {

    // For our example we'll select across two channels.
    c1 := make(chan string)
    c2 := make(chan string)

    // Each channel will receive a value after some amount
    // of time, to simulate e.g. blocking RPC operations
    // executing in concurrent goroutines.
    go func() {
        time.Sleep(time.Second * 3)
        c1 <- "one"
    }()
    go func() {
        time.Sleep(time.Second * 2)
        c2 <- "two"
    }()

    // We'll use `select` to await both of these values
    // simultaneously, printing each one as it arrives.
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-c1:
            fmt.Println("received", msg1)
        case msg2 := <-c2:
            fmt.Println("received", msg2)
        }
    }
}//原文出自【易百教程】，商业转载请联系作者获得授权，非商业请保留原文链接：https://www.yiibai.com/go/golang-select.html

