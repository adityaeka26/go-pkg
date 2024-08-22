package elasticsearch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/pkg/errors"
)

type Elasticsearch struct {
	client *elasticsearch.TypedClient
}

func NewElasticsearch(username, password string, addresses []string) (*Elasticsearch, error) {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Elasticsearch{
		client: client,
	}, nil
}

func (e *Elasticsearch) Search(ctx context.Context, indexName string, req *search.Request) (*search.Response, error) {
	res, err := e.client.Search().Index(indexName).Request(req).Do(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Elasticsearch) Index(ctx context.Context, indexName, indexId string, request any) (*index.Response, error) {
	response, err := e.client.Index(indexName).Id(indexId).Request(request).Do(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (e *Elasticsearch) GetClient() *elasticsearch.TypedClient {
	return e.client
}
