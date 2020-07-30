package twitter

import (
	"fmt"

	"github.com/IgorAndrade/analytics-twitter/server/internal/model"
	"github.com/dghubble/go-twitter/twitter"
)

func adapter(tweet *twitter.Tweet) {
	post := model.Post{}
	post.Author = tweet.User.Name
	post.Text = tweet.Text
	post.Location = tweet.User.Location
	post.Lang = tweet.Lang
	if tweet.RetweetedStatus != nil && tweet.RetweetedStatus.ExtendedTweet != nil {
		post.Location = tweet.RetweetedStatus.User.Location
		rt := tweet.RetweetedStatus.ExtendedTweet
		post.Text = rt.FullText
		hastag := make([]string, len(rt.Entities.Hashtags))
		for i, h := range rt.Entities.Hashtags {
			hastag[i] = h.Text
		}
		post.Hastags = hastag
	}
	// else {

	// 	b, _ := json.Marshal(tweet)
	// 	fmt.Println(tweet.ID, string(b))
	// }
	fmt.Println(tweet.ID, post)

}
