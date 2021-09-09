package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/berquerant/jsonhttp/handler"
	"github.com/berquerant/jsonhttp/internal/logger"
	"github.com/berquerant/jsonhttp/internal/util"
	"github.com/berquerant/jsonhttp/pb"
)

func New(value *pb.Server) *Server {
	return &Server{
		value:  value,
		logger: logger.New("[server] "),
		closeC: make(chan struct{}),
	}
}

// Server is a http server based on pb Value.
type Server struct {
	value  *pb.Server
	logger logger.Logger
	closeC chan struct{}
	server *http.Server
}

func (s *Server) serveMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/checkalive", handler.CheckAlive())
	for _, x := range s.value.GetHandlers() {
		h, ok := handler.HandlerFromAction(x.GetAction())
		if !ok {
			s.logger.Warn("cannot handle %s", util.JSON(x))
			continue
		}
		s.logger.Info("handle %s", util.JSON(x))
		func(x *pb.Handler) {
			mux.HandleFunc(x.GetPath(), func(w http.ResponseWriter, r *http.Request) {
				s.logger.Debug(`%s %d "%s %s" "%s" for %s`,
					r.RemoteAddr,
					r.ContentLength,
					r.Method,
					r.URL,
					r.UserAgent(),
					util.JSON(x),
				)
				if r.Method != x.GetMethodType().String() {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				h.ServeHTTP(w, r)
			})
		}(x)
	}
	return mux
}

// Start starts listening.
func (s *Server) Start() {
	s.logger.Info("listening on %d", s.value.GetPort())
	s.server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", s.value.GetPort()),
		Handler: s.serveMux(),
	}
	err := s.server.ListenAndServe()
	<-s.closeC
	s.logger.Info("shutdown %v", err)
}

// Close closes this server.
func (s *Server) Close(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("on shutdown %v", err)
	}
	close(s.closeC)
}
