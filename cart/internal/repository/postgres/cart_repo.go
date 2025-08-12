package postgres

import (
	"cart/internal/constants"
	"cart/internal/models"
	"cart/internal/repository/interfaces"
	"cart/pkg/postgresql"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type cartRepo struct {
	db postgresql.Client
}

func NewRepository(db postgresql.Client) interfaces.CartRepository {
	return &cartRepo{db: db}
}

func (r *cartRepo) AddItem(ctx context.Context, item models.CartItem) (int64, error) {
	var cartId int64
	query := `
		INSERT INTO cart (user_id, sku, count)
		VALUES (@userID, @sku, @count)
		ON CONFLICT (user_id, sku)
		DO UPDATE SET count = cart.count + EXCLUDED.count
		RETURNING cart.id
	`
	args := pgx.NamedArgs{
		"userID": item.UserID,
		"sku":    item.SKU,
		"count":  item.Count,
	}

	err := r.db.QueryRow(ctx, query, args).Scan(&cartId)
	if err != nil {
		return 0, err
	}

	return cartId, nil
}

func (r *cartRepo) CartItemCount(ctx context.Context, userID int64, sku uint32) (uint32, error) {
	var itemCount uint32

	query := `
		SELECT
			count
		FROM cart
		WHERE user_id = @user_id AND sku = @sku
	`
	args := pgx.NamedArgs{
		"user_id": userID,
		"sku":     sku,
	}

	err := r.db.QueryRow(ctx, query, args).Scan(&itemCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return itemCount, nil
}

func (r *cartRepo) ListItems(ctx context.Context, userID int64) ([]models.CartItem, error) {
	var items []models.CartItem

	query := `SELECT sku, count FROM cart WHERE user_id = @userID`

	args := pgx.NamedArgs{
		"userID": userID,
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item DbCartItem
		item.UserID = userID

		if err := rows.Scan(&item.SKU, &item.Count); err != nil {
			return nil, err
		}

		items = append(items, item.ToDomain())
	}

	return items, nil
}

func (r *cartRepo) DeleteCartItem(ctx context.Context, userID int64, sku uint32) error {
	query := `DELETE FROM cart WHERE user_id = @userID AND sku = @sku`

	args := pgx.NamedArgs{
		"userID": userID,
		"sku":    sku,
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

func (r *cartRepo) ClearCart(ctx context.Context, userID int64) error {
	query := `DELETE FROM cart WHERE user_id = @userID`

	args := pgx.NamedArgs{
		"userID": userID,
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
