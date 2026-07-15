package json

import (
	"testing"

	"solod.dev/so/io"
	"solod.dev/so/mem"
)

// reader is an io.Reader over a byte slice that hands back at most chunk bytes
// per Read, so a test can drive the decoder's refill and compaction paths.
type reader struct {
	data  []byte
	pos   int
	chunk int
}

func (r *reader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	end := min(r.pos+r.chunk, len(r.data))
	n := copy(p, r.data[r.pos:end])
	r.pos += n
	return n, nil
}

// writer is an io.Writer collecting everything written to it.
type writer struct{ b []byte }

func (w *writer) Write(p []byte) (int, error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

// tokens decodes doc and returns a flat, printable form of its token stream.
// A nil reader decodes doc in place; otherwise it is streamed through one.
func tokens(doc []byte, stream bool, bufSize, chunk int) ([]string, error) {
	var d Decoder
	if stream {
		d = NewReaderWith(mem.System, &reader{data: doc, chunk: chunk}, ReaderOptions{BufSize: bufSize})
	} else {
		d = NewDecoder(mem.System, doc)
	}
	defer d.Free()

	var out []string
	for d.Next() {
		switch d.Kind() {
		case KindString:
			out = append(out, "s:"+d.Str())
		case KindNumber:
			out = append(out, "n:"+string(d.tok))
		case KindBool:
			if d.Bool() {
				out = append(out, "b:true")
			} else {
				out = append(out, "b:false")
			}
		case KindNull:
			out = append(out, "null")
		case KindObjBeg:
			out = append(out, "{")
		case KindObjEnd:
			out = append(out, "}")
		case KindArrBeg:
			out = append(out, "[")
		case KindArrEnd:
			out = append(out, "]")
		}
	}
	return out, d.Err()
}

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestStreamBufferSizing(t *testing.T) {
	// The sizes streamBuffer picks decide how much the decoder reads at a time
	// and how large a token it can hold. The token limit is observable, and the
	// ./test suite covers it; the sizes themselves are not, so they are checked
	// here at the source.
	tests := []struct {
		bufSize, maxTok  int
		wantMin, wantMax int
	}{
		// The buffer starts at bufSize and grows to maxTok+1, which is what a
		// maxTok-byte token needs: the scanner ends the token on the byte after
		// it, and that byte has to fit alongside it.
		{16, 100, 16, 101},
		{64, 4096, 64, 4097},

		// bufSize already holds any token the limit allows, so the buffer starts
		// at its maximum and the growth loop never runs.
		{101, 100, 101, 101},

		// bufSize is never lowered to fit maxTok: a small token limit must not
		// shrink the reads. The buffer size wins, and the effective token limit
		// becomes bufSize-1 (4095 here). Same rule as bufio.Scanner.Buffer.
		{4096, 100, 4096, 4096},

		// Both are raised to minBufSize (16), and a 16-byte token needs 17 bytes.
		{1, 1, 17, 17},
	}
	for _, tt := range tests {
		b := streamBuffer(mem.System, &reader{}, tt.bufSize, tt.maxTok)
		if b.minSize != tt.wantMin || b.maxSize != tt.wantMax {
			t.Errorf("streamBuffer(bufSize=%d, maxTok=%d): minSize = %d, maxSize = %d; want %d and %d",
				tt.bufSize, tt.maxTok, b.minSize, b.maxSize, tt.wantMin, tt.wantMax)
		}
	}
}
