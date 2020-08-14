package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IgorAndrade/analytics-twitter/server/app/api/elasticsearch"
	rest "github.com/IgorAndrade/analytics-twitter/server/app/api/rest/webserver"
	"github.com/IgorAndrade/analytics-twitter/server/app/api/twitter"
	"github.com/IgorAndrade/analytics-twitter/server/app/config"
	"github.com/IgorAndrade/analytics-twitter/server/app/db/mongo"
	"github.com/IgorAndrade/analytics-twitter/server/internal/service"
	"golang.org/x/sync/errgroup"
)

func main() {
	b := config.NewBuilder()
	config.Define(b)
	mongo.Define(b)
	service.Define(b)
	elasticsearch.Define(b)
	config.Build(b)
	defer config.Container.Delete()

	ctx, done := context.WithCancel(context.Background())
	defer done()
	g, gctx := errgroup.WithContext(ctx)
	s := rest.NewServer(gctx, done)
	t := twitter.NewTwitterWorker(gctx, done)
	g.Go(s.Start)
	g.Go(t.Start)
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, os.Kill)
		defer s.Stop()
		defer t.Stop()
		select {
		case sig := <-signalChannel:
			fmt.Printf("Received signal: %s\n", sig)
		case <-gctx.Done():
			fmt.Printf("closing signal goroutine\n")
			return gctx.Err()
		}

		return nil
	})

	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Print("context was canceled")
		} else {
			fmt.Printf("received error: %v", err)
		}
	} else {
		fmt.Println("finished clean")
	}
}
