package log

import (
	"io"
	"sync"
)

type LockOutput struct {
	lock sync.Mutex
	io.Writer
}

func (w *LockOutput) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.Writer.Write(p)
}

func NewLockOutput(w io.Writer) *LockOutput {
	return &LockOutput{
		Writer: w,
	}
}
