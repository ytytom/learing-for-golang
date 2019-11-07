package main

import (
	"fmt"
	"time"
)

// Here's the worker, of which we'll run several
// concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding
// results on `results`. We'll sleep a second per job to
// simulate an expensive task.
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
} //接受三个参数，一个int id 一个channel 接受job 过来的值。一个channel 存放 results的值

func main() {

	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan int, 100)    //定义个100缓存的job
	results := make(chan int, 100) //定义一个100缓存的result

	// This starts up 3 workers, initially blocked
	// because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	//同时起3个 worker 此时并没有执行，因为job中没有值
	//然后关闭该通道，以指示这是所有的工作
	// Here we send 5 `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for j := 1; j <= 6; j++ {
		jobs <- j
	}
	//往job里传5个值
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 1; a <= 5; a++ {
		<-results
	}
	//最后释放掉result里的值
}
