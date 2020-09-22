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

type TwitterWorker struct {
	client *twitter.Client
	stream *twitter.Stream
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

func (t *TwitterWorker) Start() error {
	fmt.Println("Starting TwitterWorker")
	defer t.cancel()

	// filterParams := &twitter.StreamFilterParams{
	// 	Track:         []string{"Golang", "Java", "nodejs"},
	// 	Language:      []string{"pt", "en"},
	// 	StallWarnings: twitter.Bool(false),
	// }
	// stream, err := t.client.Streams.Filter(filterParams)
	// if err != nil {
	// 	return err
	// }

	// t.stream = stream

	// demux := twitter.NewSwitchDemux()
	// demux.Tweet = handlerTweet(t.poster.SavePost)

	t.listemTimeline()
	//demux.HandleChan(stream.Messages)
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
	if t.stream != nil {
		t.stream.Stop()
	}
	return nil
}

func (tw TwitterWorker) listemTimeline() {
	ticker := time.NewTicker(5 * time.Second)
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
					tw.poster.Save(tweet.ID, adapter(&tweet))
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
