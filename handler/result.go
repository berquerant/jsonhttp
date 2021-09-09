package handler

import "net/http"

type (
	ResultWriter interface {
		Headers() Headers
		Body() Body
		Status() Status
	}
	Headers interface {
		Get(key string) (string, bool)
		Set(key, value string)
		Del(key string)
		AsMap() map[string]string
	}
	Body interface {
		Get(key string) (interface{}, bool)
		Set(key string, value interface{})
		Del(key string)
		AsMap() map[string]interface{}
	}
	Status interface {
		Get() int
		Set(statusCode int)
	}
)

func NewResultWriter() ResultWriter {
	return &resultWriter{
		headers: NewHeaders(),
		body:    NewBody(),
		status:  NewStatus(),
	}
}

type resultWriter struct {
	headers Headers
	body    Body
	status  Status
}

func (s *resultWriter) Headers() Headers { return s.headers }
func (s *resultWriter) Body() Body       { return s.body }
func (s *resultWriter) Status() Status   { return s.status }

func NewHeaders() Headers {
	return &headers{
		v: map[string]string{},
	}
}

type headers struct {
	v map[string]string
}

func (s *headers) Get(key string) (string, bool) {
	x, ok := s.v[key]
	return x, ok
}
func (s *headers) Set(key, value string)    { s.v[key] = value }
func (s *headers) Del(key string)           { delete(s.v, key) }
func (s *headers) AsMap() map[string]string { return s.v }

func NewBody() Body {
	return &body{
		v: map[string]interface{}{},
	}
}

type body struct {
	v map[string]interface{}
}

func (s *body) Get(key string) (interface{}, bool) {
	x, ok := s.v[key]
	return x, ok
}
func (s *body) Set(key string, value interface{}) { s.v[key] = value }
func (s *body) Del(key string)                    { delete(s.v, key) }
func (s *body) AsMap() map[string]interface{}     { return s.v }

func NewStatus() Status {
	return &status{
		v: http.StatusOK,
	}
}

type status struct {
	v int
}

func (s *status) Get() int           { return s.v }
func (s *status) Set(statusCode int) { s.v = statusCode }
