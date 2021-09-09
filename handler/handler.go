package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/berquerant/jsonhttp/internal/util"
	"github.com/berquerant/jsonhttp/pb"
)

type Handler func(ResultWriter, *http.Request) error

func (Handler) getResponseStatus(givenStatus int, responseErr error) int {
	if responseErr == nil {
		return givenStatus
	}
	err, ok := errors.As(responseErr)
	if !ok {
		return http.StatusInternalServerError
	}
	switch err.Code() {
	case errors.UnknownError, errors.InvalidSettings:
		return http.StatusInternalServerError
	case errors.InvalidArgument, errors.OutOfRange, errors.NotFound, errors.InvalidValue, errors.Handler, errors.TypeCast:
		return http.StatusBadRequest
	default:
		return givenStatus
	}
}

func (Handler) makeErrorResponseBody(err error) []byte {
	b, _ := json.Marshal(map[string]string{
		"error":  err.Error(),
		"result": "error",
	})
	return b
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// read body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil || len(requestBody) == 0 {
		requestBody = []byte(`{}`) // use empty object
	}
	var (
		nr = r.WithContext(NewContext(requestBody).WithContext(r.Context()))
		nw = NewResultWriter()
	)
	if err != nil {
		FromContext(nr.Context()).Log().Warn("cannot read body %v", err)
	}
	FromContext(nr.Context()).Log().Debug(`"%s %s" "%s" %s`, nr.Method, nr.URL, requestBody, util.JSON(nr.Header))
	// set fundamental headers
	nw.Headers().Set("Content-Type", "application/json")
	nw.Headers().Set("X-Request-Id", FromContext(nr.Context()).ID())
	// process handler
	responseErr := h(nw, nr)
	c := FromContext(nr.Context())
	// measure duration taken to process handler
	elapsed := time.Since(c.Since())
	if err, ok := errors.As(responseErr); ok {
		c.Log().Error("responseErr trace: %s", err.Trace())
	}
	// set headers
	for k, v := range nw.Headers().AsMap() {
		w.Header().Add(k, v)
	}
	// prepare status code and body
	status := h.getResponseStatus(nw.Status().Get(), responseErr)
	body := func() []byte {
		if responseErr != nil {
			return h.makeErrorResponseBody(responseErr)
		}
		b, err := json.Marshal(nw.Body().AsMap())
		if err != nil {
			status = http.StatusInternalServerError
			c.Log().Error("failed to marshal body %v", err)
			return h.makeErrorResponseBody(err)
		}
		return b
	}()
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	// set status code
	w.WriteHeader(status)
	// set body
	if _, err := w.Write(body); err != nil {
		c.Log().Error(`failed to write body %v %s`, err, body)
	}
	// stdout log
	l := fmt.Sprintf(`%s %d %d %d %d "%s %s" "%s" %v`,
		nr.RemoteAddr,
		elapsed.Milliseconds(),
		status,
		nr.ContentLength,
		len(body),
		nr.Method,
		nr.URL,
		nr.UserAgent(),
		responseErr,
	)
	if responseErr != nil || status == http.StatusInternalServerError {
		c.Log().Error(l)
	} else {
		c.Log().Info(l)
	}
}

func HandlerFromAction(h *pb.Action) (Handler, bool) {
	switch h.GetAction().(type) {
	case *pb.Action_Return_:
		return ReturnHandler(h.GetReturn()), true
	case *pb.Action_Gateway_:
		return GatewayHandler(h.GetGateway()), true
	}
	return nil, false
}
