package home

import "sync"

//
// FanOut provides a way to fan out elements from a channel to multiple channels.
// Note that if no one is reading from the channels, it won't create back pressure.
// Further, if one reader gets stuck, the whole fan out blocks eventually.
type FanOut[T any] struct {
	in    <-chan T
	done  chan struct{}
	mOuts sync.RWMutex
	outs  []chan T
}

func NewFanOut[T any](in <-chan T) *FanOut[T] {
	f := &FanOut[T]{
		in:   in,
		done: make(chan struct{}),
	}
	go f.run()

	return f
}

func (f *FanOut[T]) run() {
msgLoop:
	for {
		select {
		case <-f.done:
			break msgLoop
		case ev, ok := <-f.in:
			if !ok {
				break msgLoop
			}
			f.mOuts.RLock()
			for _, out := range f.outs {
				select {
				case out <- ev:
				case <-f.done:
				}
			}
			f.mOuts.RUnlock()
		}
	}

	f.mOuts.RLock()
	defer f.mOuts.RUnlock()
	for _, out := range f.outs {
		close(out)
	}
}

func (f *FanOut[T]) Out() (<-chan T, func()) {
	f.mOuts.Lock()
	defer f.mOuts.Unlock()
	outChan := make(chan T, cap(f.in))
	f.outs = append(f.outs, outChan)
	return outChan, func() {
		defer close(outChan)
		f.mOuts.Lock()
		defer f.mOuts.Unlock()
		for i, c := range f.outs {
			if c == outChan {
				f.outs = append(f.outs[:i], f.outs[i+1:]...)
				break
			}
		}
	}
}

func (f *FanOut[T]) Close() {
	close(f.done)
}

type FanOutChan[T any] struct {
	fan *FanOut[T]
	In  chan T
}

func NewFanOutChan[T any](cap int) *FanOutChan[T] {
	in := make(chan T, cap)
	return &FanOutChan[T]{
		In:  in,
		fan: NewFanOut[T](in),
	}
}

func (f *FanOutChan[T]) Close() {
	f.fan.Close()
}

func (f *FanOutChan[T]) Out() (<-chan T, func()) {
	return f.fan.Out()
}
