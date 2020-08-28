package twitter

import (
	"context"
	"fmt"
	"time"

	"github.com/IgorAndrade/analytics-twitter/server/app/api"
	"github.com/IgorAndrade/analytics-twitter/server/app/config"
	"github.com/IgorAndrade/analytics-twitter/server/internal/model"
	"github.com/IgorAndrade/analytics-twitter/server/internal/repository"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterWorker struct {
	client *twitter.Client
	stream *twitter.Stream
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTwitterWorker(ctx context.Context, cancel context.CancelFunc) api.Server {
	cfg := config.Container.Get(config.CONFIG).(*config.Config)
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
	}
}

func (t *TwitterWorker) Start() error {
	fmt.Println("Starting TwitterWorker")
	defer t.cancel()

	var r repository.Elasticsearch
	if err := config.Container.Fill(repository.ELASTICSEARCH, &r); err != nil {
		return err
	}

	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"globo", "sbt"},
		Language:      []string{"pt"},
		StallWarnings: twitter.Bool(false),
	}
	stream, err := t.client.Streams.Filter(filterParams)
	if err != nil {
		return err
	}

	t.stream = stream

	demux := twitter.NewSwitchDemux()
	demux.Tweet = handlerTweet(r.Post)

	t.listemTimeline(r)
	//	demux.HandleChan(stream.Messages)
	return nil
}

func handlerTweet(send func(int64, model.Post) error) func(*twitter.Tweet) {
	return func(tweet *twitter.Tweet) {
		send(tweet.ID, adapter(tweet))
	}
}

func (t TwitterWorker) Stop() error {
	fmt.Println("Stopping TwitterWorker")
	t.cancel()
	if t.stream != nil {
		t.stream.Stop()
	}
	return nil
}

func (t TwitterWorker) listemTimeline(r repository.Elasticsearch) {
	ticker := time.NewTicker(5 * time.Second)
	time.AfterFunc(1*time.Hour, func() {
		t.cancel()
	})
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
				tweets, _, _ := t.client.Timelines.HomeTimeline(param)
				for _, t := range tweets {
					if id == t.ID {
						continue
					}
					r.Post(t.ID, adapter(&t))
					id = t.ID
				}
			}
		case <-t.ctx.Done():
			{
				ticker.Stop()
				break loop
			}

		}
	}
}
