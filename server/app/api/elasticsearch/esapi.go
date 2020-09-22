package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/IgorAndrade/analytics-twitter/server/app/config"
	"github.com/IgorAndrade/analytics-twitter/server/internal/model"
	"github.com/IgorAndrade/analytics-twitter/server/internal/repository"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mitchellh/mapstructure"
	"github.com/sarulabs/di"
	"github.com/yalp/jsonpath"
)

const INDEX string = "analytics-twitter"
const DOCUMENTTYPE string = "tweet"

type Elasticsearch struct {
	client *elasticsearch.Client
}

func Define(b *di.Builder) {
	b.Add(di.Def{
		Name:  repository.ELASTICSEARCH,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			c := ctn.Get(config.NAME).(*config.Config)
			return newServer(c.Elasticsearch)
		},
	})
}

func newServer(cfg config.Elasticsearch) (repository.Elasticsearch, error) {
	elsCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.Address,
		},
		Username: cfg.Username,
		Password: cfg.Password,
	}
	client, err := elasticsearch.NewClient(elsCfg)
	if err != nil {
		return nil, err
	}
	// Test connect
	res, err := client.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	log.Println(res)

	return &Elasticsearch{
		client: client,
	}, nil
}

func (s Elasticsearch) Post(id int64, m model.Post) error {
	req := esapi.IndexRequest{
		Index:        INDEX,
		DocumentType: DOCUMENTTYPE,
		DocumentID:   strconv.Itoa(int(id)),
		Body:         strings.NewReader(m.String()),
		Refresh:      "true",
	}

	res, err := req.Do(context.TODO(), s.client)
	if res != nil {
		res.Body.Close()
	}
	fmt.Println(id, m)
	return err
}

func (s Elasticsearch) Find(ctx context.Context, query map[string]string) ([]model.Post, error) {
	var posts []model.Post
	buf := new(bytes.Buffer)
	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"match": query,
		},
	}
	json.NewEncoder(buf).Encode(queryBody)
	es, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(INDEX),
		s.client.Search.WithBody(buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithPretty(),
	)
	if err != nil {
		return posts, err
	}
	var data interface{}
	if err := json.NewDecoder(es.Body).Decode(&data); err != nil {
		return posts, err
	}
	fmt.Println(err)
	raw, err := jsonpath.Read(data, "$.._source")
	if err != nil {
		return posts, err
	}

	mapstructure.Decode(raw, &posts)
	return posts, nil
}
