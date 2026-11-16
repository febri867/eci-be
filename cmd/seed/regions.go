package main

import (
	"context"
	"fmt"
	"time"

	"eci-backend/internal/region/model"
	"eci-backend/internal/region/repository"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type seedRegion struct {
	ID       string
	Name     string
	ParentID string
	Level    int
}

func seedRegions(ctx context.Context, db *mongo.Database) error {
	repo := repository.NewRegionRepository(db)
	now := time.Now()

	regions := []seedRegion{
		{"31", "DKI Jakarta", "", 1},
		{"32", "Jawa Barat", "", 1},
		{"33", "Jawa Tengah", "", 1},
		{"35", "Jawa Timur", "", 1},
		{"36", "Banten", "", 1},

		{"3171", "Jakarta Selatan", "31", 2},
		{"3172", "Jakarta Timur", "31", 2},
		{"3173", "Jakarta Pusat", "31", 2},
		{"3174", "Jakarta Barat", "31", 2},
		{"3175", "Jakarta Utara", "31", 2},

		{"3276", "Depok", "32", 2},
		{"3273", "Bandung", "32", 2},
		{"3271", "Bogor", "32", 2},
		{"3275", "Bekasi", "32", 2},

		{"3374", "Semarang", "33", 2},
		{"3578", "Surabaya", "35", 2},
		{"3671", "Tangerang", "36", 2},

		{"3171010", "Jagakarsa", "3171", 3},
		{"3171020", "Pasar Minggu", "3171", 3},
		{"3171030", "Cilandak", "3171", 3},
		{"3172010", "Matraman", "3172", 3},
		{"3172020", "Pulo Gadung", "3172", 3},

		{"3276010", "Sawangan", "3276", 3},
		{"3276020", "Pancoran Mas", "3276", 3},
		{"3276030", "Sukmajaya", "3276", 3},
		{"3276040", "Cimanggis", "3276", 3},
		{"3276050", "Beji", "3276", 3},
		{"3276060", "Limo", "3276", 3},

		{"3273010", "Bandung Kulon", "3273", 3},
		{"3273020", "Babakan Ciparay", "3273", 3},

		{"3374010", "Semarang Tengah", "3374", 3},
		{"3578010", "Tegalsari", "3578", 3},
		{"3671010", "Tangerang", "3671", 3},
	}

	models := make([]model.Region, 0, len(regions))

	for _, region := range regions {
		models = append(models, model.Region{
			ID:        region.ID,
			Name:      region.Name,
			ParentID:  region.ParentID,
			Level:     region.Level,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	if err := repo.UpsertMany(ctx, models); err != nil {
		return fmt.Errorf("seed regions: %w", err)
	}

	fmt.Printf("Successfully seeded %d regions\n", len(models))

	return nil
}
