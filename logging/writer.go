package logging

import (
	"bytes"
	"log"
)

type Writer struct {
	iter   int
	Prefix string

	buffer bytes.Buffer
}

func (w *Writer) output() {
	for {
		line, err := w.buffer.ReadString('\n')
		if err != nil {
			if line != "" {
				w.buffer.Reset()
				_, _ = w.buffer.Write([]byte(line)) // n is the length of p, err is always nil
			}
			break
		}

		log.Print(w.Prefix, line[0:len(line)-1])
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	n, _ := w.buffer.Write(p) // err is always nil
	w.output()
	return n, nil
}

func (w *Writer) Flush() error {
	w.output()
	if w.buffer.Len() > 0 {
		buf := w.buffer.Bytes()
		log.Print(w.Prefix, string(buf))
		w.buffer.Reset()
	}
	return nil
}
