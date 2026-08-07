package lib

import (
	"sync"
	"testing"
	"time"
)

func TestCloseChan(t *testing.T) {
	ch := make(chan int, 100)

	var wg sync.WaitGroup
	wg.Go(func() {
		for i := range 5 {
			ch <- (i + 1) * 10
			time.Sleep(time.Second)
		}
		close(ch)
	})
	wg.Go(func() {
		for {
			v, ok := <-ch
			t.Log(v, ok)

			if !ok {
				break
			}
		}
	})

	wg.Wait()
}
