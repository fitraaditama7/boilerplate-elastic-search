package elasticsearch

import (
	"boilerplate-elastic-search/internal/pkg/domain"
	"boilerplate-elastic-search/internal/pkg/storage"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var _ storage.PostStorer = PostStorage{}

type PostStorage struct {
	elastic ElasticSearch
	timeout time.Duration
}

func NewPostStorage(elastic ElasticSearch) (PostStorage, error) {
	return PostStorage{
		elastic: elastic,
		timeout: time.Second * 10,
	}, nil
}

func (p PostStorage) Insert(ctx context.Context, post storage.Post) error {
	body, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("insert: marshal: %w", err)
	}

	req := esapi.CreateRequest{
		Index:      p.elastic.alias,
		DocumentID: post.ID,
		Body:       bytes.NewReader(body),
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("insert: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		return domain.ErrConflict
	}

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}

func (p PostStorage) Update(ctx context.Context, post storage.Post) error {
	body, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("update: marshal: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      p.elastic.alias,
		DocumentID: post.ID,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("update: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return domain.ErrNotFound
	}

	if res.IsError() {
		return fmt.Errorf("update: response: %s", res.String())
	}

	return nil
}

func (p PostStorage) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      p.elastic.alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return fmt.Errorf("delete: request: %w", err)
	}

	if res.StatusCode == 404 {
		return domain.ErrNotFound
	}

	if res.IsError() {
		return fmt.Errorf("delete: response: %s", res.String())
	}
	return nil
}

func (p PostStorage) FindOne(ctx context.Context, id string) (storage.Post, error) {
	req := esapi.GetRequest{
		Index:      p.elastic.alias,
		DocumentID: id,
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := req.Do(ctx, p.elastic.client)
	if err != nil {
		return storage.Post{}, fmt.Errorf("delete: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return storage.Post{}, domain.ErrNotFound
	}

	if res.IsError() {
		return storage.Post{}, fmt.Errorf("find one: repsonse: %s", res.String())
	}

	var (
		post storage.Post
		body document
	)

	body.Source = &post

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return storage.Post{}, fmt.Errorf("find one: decode: %w", err)
	}

	return post, nil
}

func (p PostStorage) Search(ctx context.Context, keyword string) (interface{}, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"fields": []string{"title", "text", "tags"},
				"query":  keyword,
			},
			// "match": map[string]interface{}{
			// 	"title": keyword,
			// },
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("search: encode: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	res, err := p.elastic.client.Search(
		p.elastic.client.Search.WithContext(ctx),
		p.elastic.client.Search.WithIndex("post"),
		p.elastic.client.Search.WithBody(&buf),
		p.elastic.client.Search.WithTrackTotalHits(true),
		p.elastic.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("search: request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search: response: %s", res.String())
	}

	var (
		post storage.Post
		body document
		coba interface{}
	)

	body.Source = &post

	if err := json.NewDecoder(res.Body).Decode(&coba); err != nil {
		return nil, fmt.Errorf("search: decode: %w", err)
	}

	return coba, nil
}
