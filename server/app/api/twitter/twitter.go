package twitter

import (
	"context"
	"fmt"
	"time"

	"github.com/IgorAndrade/analytics-twitter/server/app/api"
	"github.com/IgorAndrade/analytics-twitter/server/app/config"
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
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"#", "canabis", "legalização"},
		Language:      []string{"pt", "en"},
		StallWarnings: twitter.Bool(false),
	}
	stream, err := t.client.Streams.Filter(filterParams)
	if err != nil {
		return err
	}

	t.stream = stream

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		//send(tweet.ID, adapter(tweet))
	}

	t.listemTimeline()
	demux.HandleChan(stream.Messages)
	return nil
}
func (t TwitterWorker) Stop() error {
	fmt.Println("Stopping TwitterWorker")
	t.cancel()
	t.stream.Stop()
	return nil
}

func (t TwitterWorker) listemTimeline() {
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
					//send(t.ID, adapter(&t))
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
