package errors

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=Code -output errors_generated.go
type Code int

const (
	UnknownError Code = iota
	InvalidArgument
	OutOfRange
	NotFound
	InvalidValue
	InvalidSettings
	Handler
	TypeCast
	TypeCoerce
	Jsonify
)

type Err struct {
	err  error
	msg  string
	code Code
}

func As(err error) (*Err, bool) {
	x, ok := err.(*Err)
	return x, ok
}

func New(code Code, msg string) *Err {
	return &Err{
		code: code,
		msg:  msg,
	}
}

func Newf(code Code, format string, v ...interface{}) *Err {
	return New(code, fmt.Sprintf(format, v...))
}

func Wrap(err error, code Code, msg string) *Err {
	return &Err{
		err:  err,
		code: code,
		msg:  msg,
	}
}

func Wrapf(err error, code Code, format string, v ...interface{}) *Err {
	return Wrap(err, code, fmt.Sprintf(format, v...))
}

func (s *Err) Code() Code { return s.code }

func (s *Err) Error() string {
	t := s
	for {
		e, ok := t.err.(*Err)
		if !ok {
			if t.msg != "" {
				return fmt.Sprintf("%s %s", t.code, t.msg)
			}
			return t.code.String()
		}
		t = e
	}
}

func (s *Err) Trace() string {
	var b strings.Builder
	t := s
	for {
		e, ok := t.err.(*Err)
		if !ok {
			b.WriteString(fmt.Sprintf("%s %s %v", t.code, t.msg, t.err))
			return b.String()
		}
		b.WriteString(fmt.Sprintf("%s %s\n", t.code, t.msg))
		t = e
	}
}
