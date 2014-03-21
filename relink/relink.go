package relink

import (
	"bytes"
	"github.com/300brand/logger"
	"io"
)

type Relink struct {
	domainMap map[string]string
	buf       *bytes.Buffer
	r         io.Reader
	filled    bool
}

var _ io.Reader = new(Relink)

var prefixes = [][]byte{
	[]byte(`https://`),
	[]byte(`http://`),
}

func New(r io.Reader, domains map[string]string) *Relink {
	return &Relink{
		domainMap: domains,
		r:         r,
		buf:       bytes.NewBuffer(make([]byte, 0, 32*1024)),
	}
}

func (r *Relink) FillBuffer() (err error) {
	_, err = io.Copy(r.buf, r.r)
	return
}

func (r *Relink) Relink() (err error) {
	doReset := false
	relinked := r.buf.Bytes()
	for from, to := range r.domainMap {
		bTo := append([]byte(`http://`), []byte(to)...)
		for _, prefix := range prefixes {
			bFrom := append(prefix, []byte(from)...)
			if bytes.Index(relinked, bFrom) > -1 {
				doReset = true
				relinked = bytes.Replace(relinked, bFrom, bTo, -1)
			}
		}
	}
	if doReset {
		logger.Debug.Printf("Doing reset and rewrite")
		r.buf.Reset()
		_, err = r.buf.Write(relinked)
	}
	return
}

func (r *Relink) Read(p []byte) (n int, err error) {
	if !r.filled {
		if err = r.FillBuffer(); err != nil {
			return
		}
		if err = r.Relink(); err != nil {
			return
		}
		r.filled = true
	}
	return r.buf.Read(p)
}
