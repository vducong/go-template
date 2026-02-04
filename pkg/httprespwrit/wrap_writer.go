package httprespwrit

import (
	"io"
	"net/http"
)

type WrapResponseWriter struct {
	http.ResponseWriter
	tee         io.Writer
	bytes       int
	code        int
	wroteHeader bool
	discard     bool
}

func New(w http.ResponseWriter) *WrapResponseWriter {
	return &WrapResponseWriter{
		ResponseWriter: w,
	}
}

func (ww *WrapResponseWriter) WriteHeader(code int) {
	if code >= 100 && code <= 199 && code != http.StatusSwitchingProtocols {
		if !ww.discard {
			ww.ResponseWriter.WriteHeader(code)
		}
	} else if !ww.wroteHeader {
		ww.code = code
		ww.wroteHeader = true

		if !ww.discard {
			ww.ResponseWriter.WriteHeader(code)
		}
	}
}

func (ww *WrapResponseWriter) Write(b []byte) (n int, err error) {
	ww.maybeWriteHeader()

	if !ww.discard {
		n, err = ww.ResponseWriter.Write(b)
		if ww.tee != nil {
			_, errTee := ww.tee.Write(b)
			if errTee != nil {
				err = errTee
			}
		}
	} else if ww.tee != nil {
		n, err = ww.tee.Write(b)
	} else {
		n, err = io.Discard.Write(b)
	}

	ww.bytes += n
	return n, err
}

func (ww *WrapResponseWriter) maybeWriteHeader() {
	if !ww.wroteHeader {
		ww.WriteHeader(http.StatusOK)
	}
}

func (ww *WrapResponseWriter) Status() int {
	return ww.code
}

func (ww *WrapResponseWriter) BytesWritten() int {
	return ww.bytes
}

func (ww *WrapResponseWriter) Tee(w io.Writer) {
	ww.tee = w
}

func (ww *WrapResponseWriter) Unpack() http.ResponseWriter {
	return ww.ResponseWriter
}

func (ww *WrapResponseWriter) Discard() {
	ww.discard = true
}
