package relink

import (
	"fmt"
	"strings"

	"bytes"
	"io"
)

type Relink struct {
	domainMap map[string]string
	buf       *bytes.Buffer
	w         io.Writer
	insideUrl bool
}

var _ io.Writer = new(Relink)

func New(w io.Writer, domains map[string]string) *Relink {
	return &Relink{
		domainMap: domains,
		w:         w,
		buf:       bytes.NewBuffer(make([]byte, 0, 256)),
	}
}

func (r *Relink) Write(p []byte) (n int, err error) {
	n = len(p)
	if !r.insideUrl {
		for idx := bytes.IndexByte(p, 'h'); idx > -1; idx = bytes.IndexByte(p, 'h') {
			fmt.Printf("I: %s\nI: %s^\n", p, strings.Repeat(" ", idx))
			switch {
			case bytes.Equal(p[idx:idx+4] == []byte(`http`):
				
			}

			fmt.Printf("Writing to w: %q\n", p[:idx])
			if _, err := r.w.Write(p[:idx]); err != nil {
				return n, err
			}
			p = p[idx+1:]
		}
		// fmt.Printf("Writing to buf: %q\n", p[idx:])
		// r.buf.Write(p[idx:])
	}
	return
}

// strRe := `https?://([^/]+).?(` + strings.Join(s.TLDs, "|") + `)['"</]`
// c.Logf("Regexp: %s", strRe)
// re := regexp.MustCompile(strRe)

// r := io.TeeReader(buf, os.Stdout)

// loc := re.Find(buf.Bytes())
// c.Logf("%s - %v", test.In, loc)
// if loc != nil {
// 	c.Log(strings.Repeat(" ", loc[0]) + "^" + strings.Repeat("-", loc[1]-loc[0]-2) + "^")
// }
