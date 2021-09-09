package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/berquerant/jsonhttp/internal/logger"
	"github.com/google/uuid"
)

type ctxKey string

const ctxKeyValue ctxKey = "ctxKeyValue"

// Context is additional information for request scope.
type Context interface {
	// ID returns the request id.
	ID() string
	// Since returns the time at which the request started.
	Since() time.Time
	Log() logger.Logger
	// Body returns the body of the request.
	Body() []byte
	WithContext(ctx context.Context) context.Context
}

func FromContext(ctx context.Context) Context {
	return ctx.Value(ctxKeyValue).(Context)
}

func NewContext(body []byte) Context {
	id := uuid.NewString()
	return &contextImpl{
		id:     id,
		since:  time.Now(),
		logger: logger.New(fmt.Sprintf("[%s] ", id)),
		body:   body,
	}
}

type contextImpl struct {
	id     string
	since  time.Time
	logger logger.Logger
	body   []byte
}

func (s *contextImpl) ID() string         { return s.id }
func (s *contextImpl) Since() time.Time   { return s.since }
func (s *contextImpl) Log() logger.Logger { return s.logger }
func (s *contextImpl) Body() []byte       { return s.body }
func (s *contextImpl) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyValue, s)
}
