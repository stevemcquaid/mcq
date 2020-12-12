package colorwriter

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/fatih/color"
)

// PrefixWriter is an implementation of io.Writer which places a prefix before
// every line.
//
// Instances of PrefixWriter are not safe to use concurrently from multiple
// goroutines.
type ColorWriter struct {
	writer io.Writer
	indent []byte
	buffer []byte
	offset int
	color  *color.Color
}

// NewPrefixWriter constructs a PrefixWriter which outputs to w and prefixes
// every line with s.
func NewPrefixWriter(w io.Writer, c *color.Color) *ColorWriter {
	return &ColorWriter{
		color:  c,
		writer: w,
		indent: copyStringToBytes(""),
		buffer: make([]byte, 0, 256),
	}
}

// Base returns the underlying writer that w outputs to.
func (w *ColorWriter) Base() io.Writer {
	return w.writer
}

// Buffered returns a byte slice of the data currently buffered in the writer.
func (w *ColorWriter) Buffered() []byte {
	return w.buffer[w.offset:]
}

// Write writes b to w, satisfies the io.Writer interface.
func (w *ColorWriter) Write(b []byte) (int, error) {
	var c int
	var n int
	var err error

	forEachLine(b, func(line []byte) bool {
		// Always buffer so the input slice doesn't escape and WriteString won't
		// copy the string (it saves a dynamic memory allocation on every call
		// to WriteString).
		w.buffer = append(w.buffer, line...)

		if chunk := w.Buffered(); isLine(chunk) {
			c, err = w.writeLine(chunk)
			w.discard(c)
		}

		n += len(line)
		return err == nil
	})

	return n, err
}

// WriteString writes s to w.
func (w *ColorWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

// Flush forces all buffered data to be flushed to the underlying writer.
func (w *ColorWriter) Flush() error {
	n, err := w.write(w.buffer)
	w.discard(n)
	return err
}

// Width satisfies the fmt.State interface.
func (w *ColorWriter) Width() (int, bool) {
	f, ok := Base(w).(fmt.State)
	if ok {
		return f.Width()
	}
	return 0, false
}

// Precision satisfies the fmt.State interface.
func (w *ColorWriter) Precision() (int, bool) {
	f, ok := Base(w).(fmt.State)
	if ok {
		return f.Precision()
	}
	return 0, false
}

// Flag satisfies the fmt.State interface.
func (w *ColorWriter) Flag(c int) bool {
	f, ok := Base(w).(fmt.State)
	if ok {
		return f.Flag(c)
	}
	return false
}

func (w *ColorWriter) writeLine(b []byte) (int, error) {
	return w.write(b)
}

func (w *ColorWriter) write(b []byte) (int, error) {
	return w.color.Fprint(w.writer, BytesToString(b))
	// return w.writer.Write(b)
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{Data: bh.Data, Len: bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func (w *ColorWriter) discard(n int) {
	if n > 0 {
		w.offset += n

		switch {
		case w.offset == len(w.buffer):
			w.buffer = w.buffer[:0]
			w.offset = 0

		case w.offset > (cap(w.buffer) / 2):
			copy(w.buffer, w.buffer[w.offset:])
			w.buffer = w.buffer[:len(w.buffer)-w.offset]
			w.offset = 0
		}
	}
}

func copyStringToBytes(s string) []byte {
	b := make([]byte, len(s))
	copy(b, s)
	return b
}

func forEachLine(b []byte, do func([]byte) bool) {
	for len(b) != 0 {
		i := bytes.IndexByte(b, '\n')

		if i < 0 {
			i = len(b)
		} else {
			i++ // include the newline character
		}

		if !do(b[:i]) {
			break
		}

		b = b[i:]
	}
}

func isLine(b []byte) bool {
	return len(b) != 0 && b[len(b)-1] == '\n'
}

var (
	_ fmt.State = (*ColorWriter)(nil)
)

// Base returns the direct base of w, which may be w itself if it had no base
// writer.
func Base(w io.Writer) io.Writer {
	if d, ok := w.(decorator); ok {
		return coalesceWriters(d.Base(), w)
	}
	return w
}

// Root returns the root writer of w, which is found by going up the list of
// base writers.
//
// The node is usually the writer where the content ends up being written.
func Root(w io.Writer) io.Writer {
	switch x := w.(type) {
	case tree:
		return coalesceWriters(x.Root(), w)
	case node:
		return coalesceWriters(Root(x.Parent()), w)
	case decorator:
		return coalesceWriters(Root(x.Base()), w)
	default:
		return w
	}
}

// Parent returns the parent writer of w, which is usually a writer of a similar
// type on tree-like writer structures.
func Parent(w io.Writer) io.Writer {
	switch x := w.(type) {
	case node:
		return coalesceWriters(x.Parent(), w)
	case decorator:
		return coalesceWriters(Parent(x.Base()), w)
	default:
		return x
	}
}

type decorator interface {
	Base() io.Writer
}

type node interface {
	Parent() io.Writer
}

type tree interface {
	Root() io.Writer
}

func coalesceWriters(writers ...io.Writer) io.Writer {
	for _, w := range writers {
		if w != nil {
			return w
		}
	}
	return nil
}
