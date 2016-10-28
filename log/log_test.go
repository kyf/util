package log

import (
	//"log"
	"sync"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	return
	p := "/var/log/6ryim_test/6ryim_test.log"
	w, err := NewWriter(p)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	var wg sync.WaitGroup
	for i := 0; i <= 500; i++ {
		go func() {
			for k := 0; k < 10000; k++ {
				_, err = w.Write([]byte("this is a log\n"))
				if err != nil {
					t.Error(err)
				}
			}
			wg.Done()
		}()
		wg.Add(1)
	}

	wg.Wait()
}

func TestLogger(t *testing.T) {
	logger, err := NewLogger("/var/log/6ryim_test/", "[6ry_test]", 0) //log.LstdFlags|log.Llongfile)
	if err != nil {
		t.Errorf("NewLogger", err)
	}
	defer logger.Close()
	for {
		select {
		case <-time.After(time.Second * 1):
			logger.Print("after 1 second ...")
			logger.Printf("hello %s", "after 5 second ...")
			logger.Error("after 1 second ...")
			logger.Errorf("hello %s", "after 5 second ...")
			logger.Fatal("after 60 second , game over...")
		}
	}
}
