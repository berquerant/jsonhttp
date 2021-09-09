package main

import (
	"context"
	"flag"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/berquerant/jsonhttp/internal/logger"
	"github.com/berquerant/jsonhttp/pb"
	"github.com/berquerant/jsonhttp/server"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	config  = flag.String("c", "server.json", "config file")
	port    = flag.Int("p", 0, "port number")
	isDebug = flag.Bool("debug", false, "enable debug log")
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	logger.IsDebug = *isDebug
	f, err := os.Open(*config)
	panicOnError(err)
	defer f.Close()
	b, err := io.ReadAll(f)
	var value pb.Server
	panicOnError(protojson.Unmarshal(b, &value))
	if *port > 0 {
		// override port number
		value.Port = int32(*port)
	}

	s := server.New(&value)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Close(ctx)
	}()
	s.Start()
}
