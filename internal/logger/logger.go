package logger

import (
	"fmt"
	"log"
)

var (
	IsDebug = false
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

type Logger interface {
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

func New(prefix string) Logger {
	return &logger{
		prefix: prefix,
	}
}

type logger struct {
	prefix string
}

func (s *logger) write(v string) {
	log.Printf("%s%s", s.prefix, v)
}

func (s *logger) Info(format string, v ...interface{}) {
	s.write(fmt.Sprintf("I %s", fmt.Sprintf(format, v...)))
}

func (s *logger) Warn(format string, v ...interface{}) {
	s.write(fmt.Sprintf("W %s", fmt.Sprintf(format, v...)))
}

func (s *logger) Error(format string, v ...interface{}) {
	s.write(fmt.Sprintf("E %s", fmt.Sprintf(format, v...)))
}

func (s *logger) Debug(format string, v ...interface{}) {
	if !IsDebug {
		return
	}
	s.write(fmt.Sprintf("D %s", fmt.Sprintf(format, v...)))
}
