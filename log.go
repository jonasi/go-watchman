package watchman

import (
	"bufio"
	"io"
	"log"
)

var logEnabled = true

func logf(fmt string, args ...interface{}) {
	if logEnabled {
		log.Printf(fmt, args...)
	}
}

func logReader(prefix string, reader io.Reader) io.Reader {
	if !logEnabled {
		return reader
	}

	r, w := io.Pipe()
	rd := io.TeeReader(reader, w)

	go logLine(prefix, r)

	return rd
}

func logWriter(prefix string, writer io.Writer) io.Writer {
	if !logEnabled {
		return writer
	}

	r, w := io.Pipe()
	wr := io.MultiWriter(writer, w)

	go logLine(prefix, r)

	return wr
}

func logLine(prefix string, r io.Reader) {
	b := bufio.NewReader(r)

	for {
		l, _ := b.ReadString('\n')
		logf("%s: %s", prefix, l)
	}
}
