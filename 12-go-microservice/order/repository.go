package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) (err error) {
	// As Order logic has many dependencies, it would be safer to use transaction in this context
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert order
	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO orders(id, created_at, account_id, total_price) 
			VALUES ($1, $2, $3, $4)
		`,
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)

	if err != nil {
		return err
	}

	// Only insert products if there are any
	if len(o.Products) > 0 {
		stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, p := range o.Products {
			_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
			if err != nil {
				return err
			}
		}

		_, err = stmt.ExecContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			o.id,
			o.created_at,
			o.account_id,
			o.total_price::money::numeric::float8,
			op.product_id,
			op.quantity
		FROM orders o
		JOIN order_products op
			ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY 
			o.id
		`,
		accountID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}
	var currentOrder *Order
	var products []OrderedProduct

	for rows.Next() {
		order := &Order{}
		orderedProduct := &OrderedProduct{}

		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		// If we're processing a new order (different from the current one)
		if currentOrder == nil || currentOrder.ID != order.ID {
			// Save the previous order if it exists
			if currentOrder != nil {
				newOrder := Order{
					ID:         currentOrder.ID,
					AccountID:  currentOrder.AccountID,
					CreatedAt:  currentOrder.CreatedAt,
					TotalPrice: currentOrder.TotalPrice,
					Products:   make([]OrderedProduct, len(products)),
				}
				copy(newOrder.Products, products)
				orders = append(orders, newOrder)
			}

			// Start a new order
			currentOrder = order
			products = []OrderedProduct{}
		}

		// Add the current product to the products slice
		products = append(products, *orderedProduct)
	}

	// Add the last order if it exists
	if currentOrder != nil {
		finalOrder := Order{
			ID:         currentOrder.ID,
			AccountID:  currentOrder.AccountID,
			CreatedAt:  currentOrder.CreatedAt,
			TotalPrice: currentOrder.TotalPrice,
			Products:   make([]OrderedProduct, len(products)),
		}
		copy(finalOrder.Products, products)
		orders = append(orders, finalOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}
