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

type NumberToWordsService struct {
	logger ActivityLogger
}

func NewNumberToWordsService(logger ActivityLogger) *NumberToWordsService {
	return &NumberToWordsService{
		logger: logger,
	}
}

var numbers = []string{
	"", "Satu", "Dua", "Tiga", "Empat",
	"Lima", "Enam", "Tujuh", "Delapan",
	"Sembilan", "Sepuluh", "Sebelas",
}

func convertNumber(number int64) string {
	switch {
	case number < 12:
		return numbers[number]
	case number < 20:
		return numbers[number-10] + " Belas"
	case number < 100:
		return strings.TrimSpace(
			convertNumber(number/10) + " Puluh " +
				convertNumber(number%10),
		)
	case number < 200:
		return strings.TrimSpace("Seratus " + convertNumber(number-100))
	case number < 1000:
		return strings.TrimSpace(
			convertNumber(number/100) + " Ratus " +
				convertNumber(number%100),
		)
	case number < 2000:
		return strings.TrimSpace("Seribu " + convertNumber(number-1000))
	case number < 1_000_000:
		return strings.TrimSpace(
			convertNumber(number/1000) + " Ribu " +
				convertNumber(number%1000),
		)
	case number < 1_000_000_000:
		return strings.TrimSpace(
			convertNumber(number/1_000_000) + " Juta " +
				convertNumber(number%1_000_000),
		)
	case number < 1_000_000_000_000:
		return strings.TrimSpace(
			convertNumber(number/1_000_000_000) + " Miliar " +
				convertNumber(number%1_000_000_000),
		)
	case number < 1_000_000_000_000_000:
		return strings.TrimSpace(
			convertNumber(number/1_000_000_000_000) + " Triliun " +
				convertNumber(number%1_000_000_000_000),
		)
	default:
		return ""
	}
}

func (s *NumberToWordsService) Convert(
	ctx context.Context,
	actorID string,
	number int64,
) (string, error) {
	if number < 0 || number >= 1_000_000_000_000_000 {
		return "", fmt.Errorf("number is outside the supported range")
	}

	result := convertNumber(number)

	if number == 0 {
		result = "Nol"
	}

	result += " Rupiah"

	if s.logger != nil {
		err := s.logger.Record(ctx, historyModel.OutboxEvent{
			Action:   "GENERATE_NUMBER_TO_WORDS",
			ActorID:  actorID,
			EntityID: "",
			After: map[string]any{
				"number": number,
				"result": result,
			},
			CreatedAt: time.Now(),
		})

		if err != nil {
			fmt.Printf("failed to record number to words activity: %v\n", err)
		}
	}

	return result, nil
}
