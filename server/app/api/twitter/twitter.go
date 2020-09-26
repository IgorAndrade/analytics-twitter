package twitter

import (
	"context"
	"fmt"
	"time"

	"github.com/IgorAndrade/analytics-twitter/server/app/api"
	"github.com/IgorAndrade/analytics-twitter/server/app/config"
	"github.com/IgorAndrade/analytics-twitter/server/internal/model"
	"github.com/IgorAndrade/analytics-twitter/server/internal/usecase"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/sarulabs/di"
)

type TwitterListener interface {
	Listen(context.Context, chan model.Post, ...ListenOptions) (func(), error)
}

type TwitterWorker struct {
	client *twitter.Client
	ctx    context.Context
	cancel context.CancelFunc
	poster usecase.Poster
}

func NewTwitterWorker(ctx context.Context, cancel context.CancelFunc, ctn di.Container) (api.Server, error) {
	var p usecase.Poster
	if err := ctn.Fill(usecase.TWITTER, &p); err != nil {
		cancel()
		return nil, err
	}

	cfg := ctn.Get(config.NAME).(*config.Config)
	config := oauth1.NewConfig(cfg.Twitter.ConsumerKey, cfg.Twitter.ConsumerSecret)
	token := oauth1.NewToken(cfg.Twitter.Token, cfg.Twitter.TokenSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	httpClient.Timeout = time.Minute

	// Twitter client
	client := twitter.NewClient(httpClient)
	return &TwitterWorker{
		client: client,
		cancel: cancel,
		ctx:    ctx,
		poster: p,
	}, nil
}

type ListenOptions func(*twitter.StreamFilterParams)

func (t *TwitterWorker) Listen(ctx context.Context, ch chan model.Post, fnc ...ListenOptions) error {
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{},
		Language:      []string{"pt", "en"},
		StallWarnings: twitter.Bool(false),
	}

	for _, f := range fnc {
		f(filterParams)
	}

	stream, err := t.client.Streams.Filter(filterParams)
	if err != nil {
		return err
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		ch <- adapter(tweet)
	}

loop:
	for {
		select {
		case msg := <-stream.Messages:
			demux.Handle(msg)
		case <-ctx.Done():
			stream.Stop()
			break loop
		}
	}
	return nil
}

func WithTrack(itens []string) ListenOptions {
	return func(param *twitter.StreamFilterParams) {
		param.Track = itens
	}
}

func WithLanguage(Langs []string) ListenOptions {
	return func(param *twitter.StreamFilterParams) {
		param.Language = Langs
	}
}

func (t *TwitterWorker) Start() error {
	fmt.Println("Starting TwitterWorker")
	defer t.cancel()

	t.listenTimeline()
	return nil
}

func handlerTweet(send func(int64, model.Post)) func(*twitter.Tweet) {
	return func(tweet *twitter.Tweet) {
		send(tweet.ID, adapter(tweet))
	}
}

func (t TwitterWorker) Stop() error {
	fmt.Println("Stopping TwitterWorker")
	t.cancel()
	return nil
}

func (tw TwitterWorker) listenTimeline() {
	ticker := time.NewTicker(time.Second)
	var id int64 = 0
loop:
	for {
		select {
		case <-ticker.C:
			{
				param := &twitter.HomeTimelineParams{
					Count: 1,
				}
				if id > 0 {
					param.SinceID = id
				}
				tweets, _, _ := tw.client.Timelines.HomeTimeline(param)
				for _, tweet := range tweets {
					if id == tweet.ID {
						continue
					}
					tw.poster.Save(adapter(&tweet))
					id = tweet.ID
				}
			}
		case <-tw.ctx.Done():
			{
				ticker.Stop()
				break loop
			}

		}
	}
}
