package elasticsearch

import (
	"context"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
)

type Elasticsearch struct {
	client      *elasticsearch.Client
	typedClient *elasticsearch.TypedClient
}

func NewElasticsearch(username, password string, addresses []string) (*Elasticsearch, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	typedClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}

	return &Elasticsearch{
		client:      client,
		typedClient: typedClient,
	}, nil
}

func (e *Elasticsearch) Search(ctx context.Context, indexName string, req *search.Request) (*search.Response, error) {
	res, err := e.typedClient.Search().Index(indexName).Request(req).Do(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Elasticsearch) SearchLowLevel(ctx context.Context, indexName string, query string) (*esapi.Response, error) {
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(indexName),
		e.client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *Elasticsearch) Index(ctx context.Context, indexName, indexId string, request any) (*index.Response, error) {
	response, err := e.typedClient.Index(indexName).Id(indexId).Request(request).Do(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (e *Elasticsearch) GetTypedClient() *elasticsearch.TypedClient {
	return e.typedClient
}

func (e *Elasticsearch) GetClient() *elasticsearch.Client {
	return e.client
}
