package pb

import (
	"net/http"

	"github.com/berquerant/jsonhttp/internal/errors"
)

// HeaderBuilder extracts a http header.
type HeaderBuilder interface {
	Build(header *http.Header) (*Value, error)
}

func NewHeaderBuilder(header *Value_Header) HeaderBuilder {
	return &headerBuilder{
		header: header,
	}
}

type headerBuilder struct {
	header *Value_Header
}

func (s *headerBuilder) Build(header *http.Header) (*Value, error) {
	if header == nil {
		return nil, errors.New(errors.InvalidSettings, "header is nil")
	}
	if v := header.Get(s.header.GetKey()); v != "" {
		return NewS(v), nil
	}
	return nil, errors.Newf(errors.NotFound, "%s is not in headers", s.header.GetKey())
}
