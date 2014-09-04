package subscription

import (
	"io"
)

type TeeReaderCloser struct {
	r io.ReadCloser
	w io.WriteCloser
}

var _ io.ReadCloser = new(TeeReaderCloser)

func (t *TeeReaderCloser) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

func (t *TeeReaderCloser) Close() error {
	t.r.Close()
	t.w.Close()
	return nil
}
