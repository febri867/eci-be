package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	historyModel "eci-backend/internal/history/model"
)

type ActivityLogger interface {
	Record(ctx context.Context, event historyModel.OutboxEvent) error
}

type StarService struct {
	logger ActivityLogger
}

func NewStarService(logger ActivityLogger) *StarService {
	return &StarService{logger: logger}
}

func (s *StarService) Generate(
	ctx context.Context,
	actorID string,
	rows int,
	pattern string,
) (string, error) {
	if rows < 1 || rows > 50 {
		return "", fmt.Errorf("rows must be between 1 and 50")
	}

	var result strings.Builder

	switch pattern {
	case "triangle":
		for i := 1; i <= rows; i++ {
			result.WriteString(strings.Repeat("*", i))
			if i < rows {
				result.WriteString("\n")
			}
		}
	case "inverted":
		for i := 1; i <= rows; i++ {
			result.WriteString(strings.Repeat(" ", rows-i))
			result.WriteString(strings.Repeat("*", i))
			if i < rows {
				result.WriteString("\n")
			}
		}
	case "pyramid":
		for i := 1; i <= rows; i++ {
			result.WriteString(strings.Repeat(" ", rows-i))
			result.WriteString(strings.Repeat("* ", i))
			if i < rows {
				result.WriteString("\n")
			}
		}
	default:
		return "", fmt.Errorf("unsupported pattern")
	}

	output := result.String()

	if s.logger != nil {
		err := s.logger.Record(ctx, historyModel.OutboxEvent{
			Action:  "GENERATE_STAR_PATTERN",
			ActorID: actorID,
			After: map[string]any{
				"rows":    rows,
				"pattern": pattern,
				"result":  output,
			},
			CreatedAt: time.Now(),
		})

		if err != nil {
			fmt.Printf("failed to record star pattern activity: %v\n", err)
		}
	}

	return output, nil
}
