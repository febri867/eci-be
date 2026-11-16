package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	historyModel "eci-backend/internal/history/model"
)

type ElasticsearchRepository struct {
	client  *http.Client
	baseURL string
	index   string
}

type HistoryFilter struct {
	Action    string
	ActorID   string
	StartDate string
	EndDate   string
}

type ElasticsearchSearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source historyModel.OutboxEvent `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`

	Aggregations ElasticsearchAggregationResponse `json:"aggregations"`
}

type ElasticsearchAggregationResponse struct {
	ActionCounts struct {
		Buckets []struct {
			Key      string `json:"key"`
			DocCount int64  `json:"doc_count"`
		} `json:"buckets"`
	} `json:"action_counts"`
}

func NewElasticsearchRepository(baseURL, index string) *ElasticsearchRepository {
	return &ElasticsearchRepository{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
		index:   index,
	}
}

func (r *ElasticsearchRepository) FindAll(
	ctx context.Context,
	page int,
	limit int,
	filter HistoryFilter,
) ([]historyModel.OutboxEvent, int64, historyModel.HistorySummary, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	must := make([]any, 0)

	if filter.Action != "" {
		must = append(must, map[string]any{
			"term": map[string]any{
				"action.keyword": filter.Action,
			},
		})
	}

	if filter.ActorID != "" {
		must = append(must, map[string]any{
			"term": map[string]any{
				"actorId.keyword": filter.ActorID,
			},
		})
	}

	dateRange := map[string]any{}

	if filter.StartDate != "" {
		dateRange["gte"] = filter.StartDate
	}

	if filter.EndDate != "" {
		dateRange["lte"] = filter.EndDate
	}

	if len(dateRange) > 0 {
		must = append(must, map[string]any{
			"range": map[string]any{
				"createdAt": dateRange,
			},
		})
	}

	searchBody := map[string]any{
		"from": (page - 1) * limit,
		"size": limit,
		"sort": []any{
			map[string]any{
				"createdAt": map[string]any{
					"order": "desc",
				},
			},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"must": must,
			},
		},
		"aggs": map[string]any{
			"action_counts": map[string]any{
				"terms": map[string]any{
					"field": "action.keyword",
					"size":  10,
				},
			},
		},
	}

	body, err := json.Marshal(searchBody)
	if err != nil {
		return nil, 0, historyModel.HistorySummary{}, err
	}

	url := fmt.Sprintf("%s/%s/_search", r.baseURL, r.index)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, 0, historyModel.HistorySummary{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := r.client.Do(req)
	if err != nil {
		return nil, 0, historyModel.HistorySummary{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, 0, historyModel.HistorySummary{}, fmt.Errorf("elasticsearch returned status %d", res.StatusCode)
	}

	var result ElasticsearchSearchResponse

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, historyModel.HistorySummary{}, err
	}

	events := make([]historyModel.OutboxEvent, 0, len(result.Hits.Hits))

	for _, hit := range result.Hits.Hits {
		events = append(events, hit.Source)
	}

	summary := historyModel.HistorySummary{
		Total: result.Hits.Total.Value,
	}

	for _, bucket := range result.Aggregations.ActionCounts.Buckets {
		switch bucket.Key {
		case "CREATE_EMPLOYEE":
			summary.Created = bucket.DocCount
		case "UPDATE_EMPLOYEE":
			summary.Updated = bucket.DocCount
		case "DELETE_EMPLOYEE":
			summary.Deleted = bucket.DocCount
		}
	}

	return events, result.Hits.Total.Value, summary, nil
}
