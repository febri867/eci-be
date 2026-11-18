package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	historyModel "eci-backend/internal/history/model"
	"eci-backend/internal/history/repository"
)

type ElasticsearchWorker struct {
	repository *repository.OutboxRepository
	client     *http.Client
	baseURL    string
	index      string
	apiKey     string
}

func NewElasticsearchWorker(
	repo *repository.OutboxRepository,
	baseURL string,
	index string,
	apiKey string,
) *ElasticsearchWorker {
	return &ElasticsearchWorker{
		repository: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
		index:   index,
		apiKey:  apiKey,
	}
}

func (w *ElasticsearchWorker) Run(ctx context.Context) {
	w.ProcessPending(ctx)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.ProcessPending(ctx)
		}
	}
}

func (w *ElasticsearchWorker) ProcessPending(ctx context.Context) {
	events, err := w.repository.FindPending(ctx, 50)
	if err != nil {
		fmt.Printf("failed to fetch pending events: %v\n", err)
		return
	}

	for _, event := range events {
		if ctx.Err() != nil {
			return
		}

		if err := w.sendToElasticsearch(ctx, event); err != nil {
			fmt.Printf("failed to send event %s: %v\n", event.ID.Hex(), err)
			_ = w.repository.MarkFailed(ctx, event.ID, err.Error())
			continue
		}

		if err := w.repository.MarkProcessed(ctx, event.ID); err != nil {
			fmt.Printf("failed to mark event %s as processed: %v\n", event.ID.Hex(), err)
		}
	}
}

func (w *ElasticsearchWorker) sendToElasticsearch(
	ctx context.Context,
	event historyModel.OutboxEvent,
) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/%s/_doc/%s",
		w.baseURL,
		w.index,
		event.ID.Hex(),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey "+w.apiKey)

	res, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("elasticsearch returned status %d", res.StatusCode)
	}

	return nil
}
