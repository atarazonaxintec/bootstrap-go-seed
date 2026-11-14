package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"microservice/internal/orders/domain/entity"
	domain "microservice/internal/orders/domain/repository"
	"microservice/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.OrdersRepository = (*PostgresOrdersRepository)(nil)

type PostgresOrdersRepository struct {
	db database.Database
}

func NewPostgresOrdersRepository(db database.Database) *PostgresOrdersRepository {
	return &PostgresOrdersRepository{db: db}
}

// Implementamos la interfaz del repository
func (r *PostgresOrdersRepository) Ping(ctx context.Context) error {
	pool := r.db.Native().(*pgxpool.Pool)
	return pool.Ping(ctx)
}

func (r *PostgresOrdersRepository) CreateOrder(ctx context.Context, rut string) (*entity.Order, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create order transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	order, err := scanOrder(tx.QueryRow(ctx, `
		INSERT INTO orders (rut, sub_total, status)
		VALUES ($1, 0, 'en_preparacion')
		RETURNING id::text, rut, sub_total, status, created_at::text, updated_at::text
	`, rut))
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO order_statuses (order_id, status)
		VALUES ($1, $2)
	`, order.Id, order.Status); err != nil {
		return nil, fmt.Errorf("record initial order status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create order transaction: %w", err)
	}
	return order, nil
}

func (r *PostgresOrdersRepository) UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	if order == nil {
		return nil, errors.New("update order: order cannot be nil")
	}

	pool := r.db.Native().(*pgxpool.Pool)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin update order transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var previousStatus string
	if err := tx.QueryRow(ctx, `
		SELECT status
		FROM orders
		WHERE id = $1
		FOR UPDATE
	`, order.Id).Scan(&previousStatus); err != nil {
		return nil, fmt.Errorf("find order to update: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM order_items WHERE order_id = $1`, order.Id); err != nil {
		return nil, fmt.Errorf("replace order items: %w", err)
	}

	for _, item := range order.OrderItems {
		if _, err := tx.Exec(ctx, `
			INSERT INTO order_items (order_id, sku, quantity)
			VALUES ($1, $2, $3)
		`, order.Id, item.Sku, item.Quantity); err != nil {
			return nil, fmt.Errorf("insert order item %q: %w", item.Sku, err)
		}
	}

	if _, err := tx.Exec(ctx, `
		UPDATE orders
		SET rut = $2,
		status = $3,
		sub_total = COALESCE((
			SELECT SUM(products.price * order_items.quantity)
			FROM order_items
			JOIN products ON products.sku = order_items.sku
			WHERE order_items.order_id = $1
		), 0),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, order.Id, order.Rut, order.Status); err != nil {
		return nil, fmt.Errorf("update order and recalculate subtotal: %w", err)
	}

	if previousStatus != string(order.Status) {
		if _, err := tx.Exec(ctx, `
			INSERT INTO order_statuses (order_id, status)
			VALUES ($1, $2)
		`, order.Id, order.Status); err != nil {
			return nil, fmt.Errorf("record order status change: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit update order transaction: %w", err)
	}
	return r.GetOrderById(ctx, order.Id)
}

func (r *PostgresOrdersRepository) GetOrderById(ctx context.Context, id string) (*entity.Order, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	order, err := scanOrder(pool.QueryRow(ctx, `
		SELECT id::text, rut, sub_total, status, created_at::text, updated_at::text
		FROM orders
		WHERE id = $1
	`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order %q not found: %w", id, err)
		}
		return nil, fmt.Errorf("get order %q: %w", id, err)
	}

	if err := loadOrderDetails(ctx, pool, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (r *PostgresOrdersRepository) GetOrdersByRut(ctx context.Context, rut string) ([]*entity.Order, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	rows, err := pool.Query(ctx, `
		SELECT id::text, rut, sub_total, status, created_at::text, updated_at::text
		FROM orders
		WHERE rut = $1
		ORDER BY created_at DESC
	`, rut)
	if err != nil {
		return nil, fmt.Errorf("get orders for rut %q: %w", rut, err)
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, fmt.Errorf("scan order for rut %q: %w", rut, err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders for rut %q: %w", rut, err)
	}

	for _, order := range orders {
		if err := loadOrderDetails(ctx, pool, order); err != nil {
			return nil, err
		}
	}
	return orders, nil
}

type orderScanner interface {
	Scan(dest ...any) error
}

type orderDetailQuerier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func scanOrder(row orderScanner) (*entity.Order, error) {
	order := &entity.Order{}
	var status string
	if err := row.Scan(
		&order.Id,
		&order.Rut,
		&order.SubTotal,
		&status,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return nil, err
	}
	order.Status = entity.OrderStatus(status)
	order.OrderItems = make([]entity.OrderItem, 0)
	order.Payments = make([]entity.OrderPayment, 0)
	return order, nil
}

func loadOrderDetails(ctx context.Context, querier orderDetailQuerier, order *entity.Order) error {
	itemRows, err := querier.Query(ctx, `
		SELECT order_id::text, sku, quantity, created_at::text, updated_at::text
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at
	`, order.Id)
	if err != nil {
		return fmt.Errorf("get items for order %q: %w", order.Id, err)
	}
	for itemRows.Next() {
		var item entity.OrderItem
		if err := itemRows.Scan(&item.OrderId, &item.Sku, &item.Quantity, &item.CreatedAt, &item.UpdatedAt); err != nil {
			itemRows.Close()
			return fmt.Errorf("scan item for order %q: %w", order.Id, err)
		}
		order.OrderItems = append(order.OrderItems, item)
	}
	if err := itemRows.Err(); err != nil {
		itemRows.Close()
		return fmt.Errorf("iterate items for order %q: %w", order.Id, err)
	}
	itemRows.Close()

	paymentRows, err := querier.Query(ctx, `
		SELECT id::text, order_id::text, payment_method, amount, paid_at::text, status
		FROM order_payments
		WHERE order_id = $1
		ORDER BY paid_at
	`, order.Id)
	if err != nil {
		return fmt.Errorf("get payments for order %q: %w", order.Id, err)
	}
	defer paymentRows.Close()

	for paymentRows.Next() {
		var payment entity.OrderPayment
		var status string
		if err := paymentRows.Scan(
			&payment.Id,
			&payment.OrderId,
			&payment.PaymentMethod,
			&payment.Amount,
			&payment.PaidAt,
			&status,
		); err != nil {
			return fmt.Errorf("scan payment for order %q: %w", order.Id, err)
		}
		payment.Status = entity.OrderPaymentStatus(status)
		order.Payments = append(order.Payments, payment)
	}
	if err := paymentRows.Err(); err != nil {
		return fmt.Errorf("iterate payments for order %q: %w", order.Id, err)
	}
	return nil
}
