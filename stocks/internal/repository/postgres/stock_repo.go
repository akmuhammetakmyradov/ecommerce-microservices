package postgres

import (
	"context"
	"stocks/internal/constants"
	"stocks/internal/models"
	"stocks/internal/repository/interfaces"
	"stocks/pkg/postgresql"

	tmsql "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
)

type stockRepo struct {
	db     postgresql.Client
	getter *tmsql.CtxGetter
}

func NewRepository(db postgresql.Client, getter *tmsql.CtxGetter) interfaces.StockRepository {
	return &stockRepo{
		db:     db,
		getter: getter,
	}
}

func (r *stockRepo) AddItem(ctx context.Context, item models.StockItem) (string, error) {
	var (
		xmax   uint32
		result = "sku_created"
	)

	txOrDb := r.getter.DefaultTrOrDB(ctx, r.db.(*postgresql.PgClient))

	query := `
		INSERT INTO items (
			user_id, sku, count, price, location
		) VALUES (
			@user_id, @sku, @count, @price, @location
		) 
		ON CONFLICT (sku) DO UPDATE SET
			count = items.count + EXCLUDED.count,
			price = EXCLUDED.price,
			user_id = EXCLUDED.user_id,
			location = EXCLUDED.location,
			updated_at = CURRENT_TIMESTAMP
		RETURNING xmax
	`
	args := pgx.NamedArgs{
		"user_id":  item.UserID,
		"sku":      item.SKU,
		"count":    item.Count,
		"price":    item.Price,
		"location": item.Location,
	}

	err := txOrDb.QueryRow(ctx, query, args).Scan(&xmax)
	if err != nil {
		return result, err
	}

	if xmax != 0 {
		result = "sku_changed"
	}

	return result, nil
}

func (r *stockRepo) DeleteItem(ctx context.Context, sku uint32) error {
	query := "DELETE FROM items WHERE sku = @sku"

	args := pgx.NamedArgs{
		"sku": sku,
	}

	cmdTag, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return constants.ErrNotRowAffected
	}

	return nil
}

func (r *stockRepo) GetItemsByLocation(ctx context.Context, location string, userID, limit, offset int64) ([]models.StockItem, error) {
	var result []models.StockItem

	txOrDb := r.getter.DefaultTrOrDB(ctx, r.db.(*postgresql.PgClient))

	query := `
		SELECT 
			i.sku, i.count, s.name, 
			s.type, i.price, i.location 
		FROM items i
		LEFT JOIN sku s
			ON i.sku = s.sku_id
		WHERE i.location = @location AND i.user_id = @user_id
		LIMIT @limit OFFSET @offset
	`
	args := pgx.NamedArgs{
		"location": location,
		"user_id":  userID,
		"limit":    limit,
		"offset":   offset,
	}

	rows, err := txOrDb.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item DbStockItem

		err = rows.Scan(
			&item.SKU, &item.Count, &item.Name,
			&item.Type, &item.Price, &item.Location,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, item.ToDomain())
	}

	return result, nil
}

func (r *stockRepo) CountItemsByLocation(ctx context.Context, location string, userID int64) (int64, error) {
	var count int64

	txOrDb := r.getter.DefaultTrOrDB(ctx, r.db.(*postgresql.PgClient))

	query := `SELECT COUNT(*) FROM items WHERE location = @location AND user_id = @user_id`

	args := pgx.NamedArgs{
		"location": location,
		"user_id":  userID,
	}

	err := txOrDb.QueryRow(ctx, query, args).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *stockRepo) GetItemBySKU(ctx context.Context, sku uint32) (models.StockItem, error) {
	var item DbStockItem

	query := `
		SELECT 
			i.id, i.user_id, i.sku, i.count, s.name, 
			s.type, i.price, i.location
		FROM items i
		LEFT JOIN sku s
			ON i.sku = s.sku_id
		WHERE i.sku = @sku
	`
	args := pgx.NamedArgs{
		"sku": sku,
	}

	err := r.db.QueryRow(ctx, query, args).Scan(
		&item.ID, &item.UserID, &item.SKU, &item.Count,
		&item.Name, &item.Type, &item.Price, &item.Location,
	)

	if err != nil {
		return models.StockItem{}, err
	}

	return item.ToDomain(), nil
}

func (r *stockRepo) GetSKUByID(ctx context.Context, skuID uint32) (models.SKU, error) {
	txOrDb := r.getter.DefaultTrOrDB(ctx, r.db.(*postgresql.PgClient))
	var sku DbSKU

	query := `
		SELECT 
			s.sku_id, s.name, s.type, i.user_id
		FROM sku s
		LEFT JOIN items i
			ON s.sku_id = i.sku
		WHERE s.sku_id = @sku_id
	`
	args := pgx.NamedArgs{
		"sku_id": skuID,
	}

	err := txOrDb.QueryRow(ctx, query, args).Scan(&sku.SKUID, &sku.Name, &sku.Type, &sku.UserID)
	if err != nil {
		return models.SKU{}, err
	}

	return sku.ToDomain(), nil
}
