package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IgorAndrade/analytics-twitter/server/app/api"
	"github.com/IgorAndrade/analytics-twitter/server/app/api/elasticsearch"
	rest "github.com/IgorAndrade/analytics-twitter/server/app/api/rest/webserver"
	"github.com/IgorAndrade/analytics-twitter/server/app/api/twitter"
	"github.com/IgorAndrade/analytics-twitter/server/app/config"
	"github.com/IgorAndrade/analytics-twitter/server/app/db/mongo"
	"github.com/IgorAndrade/analytics-twitter/server/internal/service"
	"github.com/IgorAndrade/analytics-twitter/server/internal/usecase"
	"github.com/sarulabs/di"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctn := initContainer()
	sub, _ := ctn.SubContainer()
	defer ctn.Delete()

	ctx, done := context.WithCancel(context.Background())
	defer done()
	g, gctx := errgroup.WithContext(ctx)

	s := rest.NewServer(gctx, done)
	t, err := twitter.NewTwitterWorker(gctx, done, sub)
	if err != nil {
		log.Fatal(err)
	}
	serv := api.List{s, t}
	serv.StartAll(g)
	g.Go(waitSignalChannel(gctx, serv))

	err = g.Wait()
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

func waitSignalChannel(gctx context.Context, serv api.List) func() error {
	return func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, os.Kill, syscall.SIGSEGV)
		defer serv.StopAll()

		select {
		case sig := <-signalChannel:
			fmt.Printf("Received signal: %s\n", sig)
		case <-gctx.Done():
			fmt.Printf("closing signal goroutine\n")
			return gctx.Err()
		}

		return nil
	}
}

func initContainer() di.Container {
	b := config.NewBuilder(
		config.Define,
		mongo.Define,
		service.Define,
		elasticsearch.Define,
		usecase.Define,
	)
	return config.Build(b)
}
