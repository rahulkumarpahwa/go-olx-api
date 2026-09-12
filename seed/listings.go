package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rahulkumarpahwa/go-olx-api/internal/types"
)

// Fixed IDs shared by every seeded row so that all listings reference the same
// user_id and category_id (the requirement) and the foreign keys always resolve.
var (
	seedUserID     = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	seedCategoryID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

// dummyListing is a single row to insert into the listings table.
// price is stored in paise (BIGINT, values in paise).
type dummyListing struct {
	title       string
	description string
	price       int64
	city        string
	status      types.ListingStatus
}

var dummyListings = []dummyListing{
	{"iPhone 13 128GB - Midnight", "Excellent condition, includes box and charger.", 4200000, "Bengaluru", types.ListingStatusActive},
	{"Samsung Galaxy S22 5G", "New condition, only 2 months old.", 3500000, "Delhi", types.ListingStatusActive},
	{"Vintage Wooden Dining Table", "Solid teak wood, 6 seater.", 8500000, "Mumbai", types.ListingStatusActive},
	{"Casio G-Shock Watch", "Unisex, water resistant.", 250000, "Pune", types.ListingStatusActive},
	{"Mountain Bike 26 inch", "Gearless, ideal for beginners.", 950000, "Hyderabad", types.ListingStatusInactive},
	{"MacBook Air M2 13 inch", "16GB RAM, 256GB SSD, mint condition.", 10500000, "Chennai", types.ListingStatusActive},
	{"Sony WH-1000XM4 Headphones", "Noise cancelling, barely used.", 1500000, "Kolkata", types.ListingStatusActive},
	{"Office Chair Ergonomic", "Mesh back, adjustable height.", 450000, "Jaipur", types.ListingStatusActive},
	{"Canon EOS R10 Camera Body", "Includes 18-45mm kit lens.", 6800000, "Ahmedabad", types.ListingStatusActive},
	{"Four-Poster Bed with Mattress", "Queen size, premium foam mattress.", 1200000, "Surat", types.ListingStatusInactive},
}

// SeedListings seeds the database with 10 dummy listing entries. Every entry
// shares the same user_id and the same category_id.
func SeedListings(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ensure the shared user and category exist so the foreign keys on
	// listings.user_id and listings.category_id resolve.
	if err := ensureSeedUser(ctx, db); err != nil {
		return err
	}
	if err := ensureSeedCategory(ctx, db); err != nil {
		return err
	}

	const query = `
		INSERT INTO listings (title, description, price, city, status, user_id, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // no-op after successful commit

	for _, l := range dummyListings {
		if _, err := tx.ExecContext(ctx, query, l.title, l.description, l.price, l.city, l.status, seedUserID, seedCategoryID); err != nil {
			return fmt.Errorf("seed listing %q: %w", l.title, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// ensureSeedUser inserts the single shared user when it does not already exist.
func ensureSeedUser(ctx context.Context, db *sql.DB) error {
	const query = `
		INSERT INTO users (id, name, email, password)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING`

	if _, err := db.ExecContext(ctx, query, seedUserID, "John Doe", "john.doe@example.com", "password123"); err != nil {
		return fmt.Errorf("seed user: %w", err)
	}
	return nil
}

// ensureSeedCategory inserts the single shared category when it does not exist.
func ensureSeedCategory(ctx context.Context, db *sql.DB) error {
	const query = `
		INSERT INTO category (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING`

	if _, err := db.ExecContext(ctx, query, seedCategoryID, "Electronics"); err != nil {
		return fmt.Errorf("seed category: %w", err)
	}
	return nil
}
