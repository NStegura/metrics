package metricsapi

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func (c *compressWriter) Write(p []byte) (int, error) {
	count, err := c.zw.Write(p)
	if err != nil {
		return 0, fmt.Errorf("failed to write comperessed resp: %w", err)
	}
	return count, nil
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < http.StatusMultipleChoices {
		c.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	if err := c.zw.Close(); err != nil {
		return fmt.Errorf("failed to close writer %w", err)
	}
	return nil
}

type compressReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to init gzip reader %w", err)
	}

	return &compressReader{
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	count, err := c.zr.Read(p)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return count, io.EOF
		}
		return 0, fmt.Errorf("failed to read %w", err)
	}
	return count, nil
}

func (c *compressReader) Close() error {
	if err := c.zr.Close(); err != nil {
		return fmt.Errorf("failed to close gzip reader %w", err)
	}
	return nil
}

func (s *APIServer) gzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := &compressWriter{ResponseWriter: w, zw: gzip.NewWriter(w)}
			ow = cw
			defer func() {
				if err := cw.Close(); err != nil {
					s.logger.Error(err)
				}
			}()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				s.logger.Errorf("failed to init compress reader for body, err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func() {
				if err := cr.Close(); err != nil {
					s.logger.Error(err)
				}
			}()
		}
		h.ServeHTTP(ow, r)
	})
}
