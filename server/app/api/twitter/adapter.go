package twitter

import (
	"fmt"
	"regexp"
	"time"

	"github.com/IgorAndrade/analytics-twitter/server/internal/model"
	"github.com/dghubble/go-twitter/twitter"
)

func adapter(tweet *twitter.Tweet) model.Post {
	post := model.Post{}
	date, _ := time.Parse(time.RubyDate, tweet.CreatedAt)
	post.CreatedAt = date.Local()
	post.Author = tweet.User.Name
	post.Text = tweet.Text
	post.Location = tweet.User.Location
	post.Lang = tweet.Lang
	post.Hastags = getHashtag(tweet.Text)
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
	//fmt.Println(tweet)
	return post
}

func getHashtag(str string) []string {
	hastag := make([]string, 0)
	var re = regexp.MustCompile(`#\S+`)
	for _, match := range re.FindAllString(str, -1) {
		hastag = append(hastag, match)
		fmt.Println("hashtad :", match)
	}
	return hastag
}
