package home

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFanOut(t *testing.T) {
	t.Run("no-reader", func(t *testing.T) {
		in := make(chan int, 1)
		fan := NewFanOut[int](in)
		defer fan.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		for i := 0; i < 100; i++ {
			select {
			case in <- i:
			case <-ctx.Done():
				t.Errorf("pushing to a fan-out where no one is listening should not block")
			}
		}
	})

	t.Run("one-reader", func(t *testing.T) {
		in := make(chan int, 0)
		fan := NewFanOut[int](in)
		defer fan.Close()
		reader, closer := fan.Out()

		var (
			outElem int
			read    = make(chan struct{})
		)
		go func() {
			defer close(read)
			outElem = <-reader
		}()
		go func() { in <- 23 }()

		<-read
		require.Len(t, fan.outs, 1)
		closer()
		require.Len(t, fan.outs, 0)
		require.EqualValues(t, outElem, 23)
	})

	t.Run("multi-reader", func(t *testing.T) {
		in := make(chan int, 0)
		fan := NewFanOut[int](in)
		defer fan.Close()
		reader1, closer1 := fan.Out()
		reader2, closer2 := fan.Out()

		var (
			outElem1 int
			outElem2 int
			read     = make(chan struct{})
		)
		go func() {
			defer close(read)
			outElem1 = <-reader1
			outElem2 = <-reader2
		}()
		go func() { in <- 23 }()

		<-read
		require.Len(t, fan.outs, 2)
		closer1()
		closer2()
		require.Len(t, fan.outs, 0)
		require.EqualValues(t, outElem1, 23)
		require.EqualValues(t, outElem2, 23)
	})
}
