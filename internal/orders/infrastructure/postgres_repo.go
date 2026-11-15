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
var _ domain.OrderProductsRepository = (*PostgresProductsRepository)(nil)
var _ domain.OrderPaymentRepository = (*PostgresPaymentsRepository)(nil)

type PostgresOrdersRepository struct {
	db database.Database
}

type PostgresProductsRepository struct {
	db database.Database
}

type PostgresPaymentsRepository struct {
	db database.Database
}

// Constructor functions for the repositories
func NewPostgresOrdersRepository(db database.Database) *PostgresOrdersRepository {
	return &PostgresOrdersRepository{db: db}
}

func NewPostgresProductsRepository(db database.Database) *PostgresProductsRepository {
	return &PostgresProductsRepository{db: db}
}

func NewPostgresPaymentsRepository(db database.Database) *PostgresPaymentsRepository {
	return &PostgresPaymentsRepository{db: db}
}

// Implementamos la interfaz del repository
func (r *PostgresOrdersRepository) Ping(ctx context.Context) error {
	pool := r.db.Native().(*pgxpool.Pool)
	return pool.Ping(ctx)
}

func (r *PostgresOrdersRepository) CreateOrder(ctx context.Context, rut string) (*entity.OrderEntity, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create order transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	order, err := scanOrder(tx.QueryRow(ctx, `
		INSERT INTO go_seed.orders (rut, sub_total, status)
		VALUES ($1, 0, 'inicializado')
		RETURNING id::text, rut, sub_total, status, created_at::text, updated_at::text
	`, rut))
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO go_seed.order_statuses (order_id, status)
		VALUES ($1, $2)
	`, order.Id, order.Status); err != nil {
		return nil, fmt.Errorf("record initial order status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create order transaction: %w", err)
	}
	return order, nil
}

func (r *PostgresOrdersRepository) UpdateOrder(ctx context.Context, order *entity.OrderEntity) (*entity.OrderEntity, error) {
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
		FROM go_seed.orders
		WHERE id = $1
		FOR UPDATE
	`, order.Id).Scan(&previousStatus); err != nil {
		return nil, fmt.Errorf("find order to update: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM go_seed.order_items WHERE order_id = $1`, order.Id); err != nil {
		return nil, fmt.Errorf("replace order items: %w", err)
	}

	for _, item := range order.OrderItems {
		if _, err := tx.Exec(ctx, `
			INSERT INTO go_seed.order_items (order_id, sku, quantity)
			VALUES ($1, $2, $3)
		`, order.Id, item.Sku, item.Quantity); err != nil {
			return nil, fmt.Errorf("insert order item %q: %w", item.Sku, err)
		}
	}

	if _, err := tx.Exec(ctx, `
		UPDATE go_seed.orders
		SET rut = $2,
		status = $3,
		sub_total = COALESCE((
			SELECT SUM(product.price * item.quantity)
			FROM go_seed.order_items AS item
			JOIN go_seed.products AS product ON product.sku = item.sku
			WHERE item.order_id = $1
		), 0),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, order.Id, order.Rut, order.Status); err != nil {
		return nil, fmt.Errorf("update order and recalculate subtotal: %w", err)
	}

	if previousStatus != string(order.Status) {
		if _, err := tx.Exec(ctx, `
			INSERT INTO go_seed.order_statuses (order_id, status)
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

func (r *PostgresOrdersRepository) GetOrderById(ctx context.Context, id string) (*entity.OrderEntity, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	order, err := scanOrder(pool.QueryRow(ctx, `
		SELECT id::text, rut, sub_total, status, created_at::text, updated_at::text
		FROM go_seed.orders
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

func (r *PostgresOrdersRepository) GetOrdersByRut(ctx context.Context, rut string) ([]*entity.OrderEntity, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	rows, err := pool.Query(ctx, `
		SELECT id::text, rut, sub_total, status, created_at::text, updated_at::text
		FROM go_seed.orders
		WHERE rut = $1
		ORDER BY created_at DESC
	`, rut)
	if err != nil {
		return nil, fmt.Errorf("get orders for rut %q: %w", rut, err)
	}
	defer rows.Close()

	orders := make([]*entity.OrderEntity, 0)
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

func scanOrder(row orderScanner) (*entity.OrderEntity, error) {
	order := &entity.OrderEntity{}
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
	order.OrderItems = make([]entity.OrderItemEntity, 0)
	order.Payments = make([]entity.OrderPaymentEntity, 0)
	return order, nil
}

func loadOrderDetails(ctx context.Context, querier orderDetailQuerier, order *entity.OrderEntity) error {
	itemRows, err := querier.Query(ctx, `
		SELECT order_id::text, sku, quantity, created_at::text, updated_at::text
		FROM go_seed.order_items
		WHERE order_id = $1
		ORDER BY created_at
	`, order.Id)
	if err != nil {
		return fmt.Errorf("get items for order %q: %w", order.Id, err)
	}
	for itemRows.Next() {
		var item entity.OrderItemEntity
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
		FROM go_seed.order_payments
		WHERE order_id = $1
		ORDER BY paid_at
	`, order.Id)
	if err != nil {
		return fmt.Errorf("get payments for order %q: %w", order.Id, err)
	}
	defer paymentRows.Close()

	for paymentRows.Next() {
		var payment entity.OrderPaymentEntity
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

func (r *PostgresOrdersRepository) GetOrderAndReserveForUpdate(ctx context.Context, id string) (*entity.OrderEntity, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction to get order for update: %w", err)
	}
	defer tx.Rollback(ctx)

	order, err := scanOrder(tx.QueryRow(ctx, `
		SELECT id::text, rut, sub_total, status, created_at::text, updated_at::text
		FROM go_seed.orders
		WHERE id = $1
		FOR UPDATE
	`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order %q not found: %w", id, err)
		}
		return nil, fmt.Errorf("get order %q for update: %w", id, err)
	}

	if err := loadOrderDetails(ctx, tx, order); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction to get order for update: %w", err)
	}
	return order, nil
}

// ProductsRepository interface implementation
func (r *PostgresProductsRepository) FindBySku(ctx context.Context, sku string) (*entity.ProductEntity, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	product := &entity.ProductEntity{}
	if err := pool.QueryRow(ctx, `
		SELECT sku, name, price, created_at::text, updated_at::text, description, stock_quantity
		FROM go_seed.products
		WHERE sku = $1
	`, sku).Scan(&product.Sku, &product.Name, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.Description, &product.StockQuantity); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product %q not found: %w", sku, err)
		}
		return nil, fmt.Errorf("get product %q: %w", sku, err)
	}
	return product, nil
}

func (r *PostgresProductsRepository) FindAll(ctx context.Context) ([]*entity.ProductEntity, error) {
	pool := r.db.Native().(*pgxpool.Pool)

	rows, err := pool.Query(ctx, `
		SELECT sku, name, price, created_at::text, updated_at::text, description, stock_quantity
		FROM go_seed.products
	`)
	if err != nil {
		return nil, fmt.Errorf("find all products: %w", err)
	}
	defer rows.Close()

	products := make([]*entity.ProductEntity, 0)
	for rows.Next() {
		product := &entity.ProductEntity{}
		if err := rows.Scan(&product.Sku, &product.Name, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.Description, &product.StockQuantity); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}
	return products, nil
}

// PaymentRepository interface implementation
func (r *PostgresPaymentsRepository) CreateEmpty(ctx context.Context, orderID string) (*entity.OrderPaymentEntity, error) {
	sql := `
		INSERT INTO go_seed.order_payments (order_id, status)
		VALUES ($1, $2)
		RETURNING id::text, order_id::text, payment_method, amount, paid_at::text, status
	`

	pool := r.db.Native().(*pgxpool.Pool)
	payment := &entity.OrderPaymentEntity{}
	var status string
	if err := pool.QueryRow(ctx, sql, orderID, entity.Empty).Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.PaymentMethod,
		&payment.Amount,
		&payment.PaidAt,
		&status,
	); err != nil {
		return nil, fmt.Errorf("create empty payment for order %q: %w", orderID, err)
	}
	return payment, nil
}

func (r *PostgresPaymentsRepository) FindByOrderID(ctx context.Context, id string) ([]*entity.OrderPaymentEntity, error) {
	pool := r.db.Native().(*pgxpool.Pool)
	payment := []*entity.OrderPaymentEntity{}
	var status string
	sql := `
		SELECT id::text, order_id::text, payment_method, amount, paid_at::text, status
		FROM go_seed.order_payments
		WHERE order_id = $1
	`
	rows, err := pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("find payments for order %q: %w", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		p := &entity.OrderPaymentEntity{}
		if err := rows.Scan(
			&p.Id,
			&p.OrderId,
			&p.PaymentMethod,
			&p.Amount,
			&p.PaidAt,
			&status,
		); err != nil {
			return nil, fmt.Errorf("scan payment for order %q: %w", id, err)
		}
		p.Status = entity.OrderPaymentStatus(status)
		payment = append(payment, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate payments for order %q: %w", id, err)
	}
	return payment, nil
}

func (r *PostgresPaymentsRepository) Pay(ctx context.Context, id string, paymentMethod string, amount float64) (*entity.OrderPaymentEntity, error) {

	sql := `
		UPDATE go_seed.order_payments
		SET payment_method = $2, amount = $3, paid_at = CURRENT_TIMESTAMP, status = $4
		WHERE order_id = $1
		RETURNING id::text, order_id::text, payment_method, amount, paid_at::text, status
	`

	pool := r.db.Native().(*pgxpool.Pool)
	payment := &entity.OrderPaymentEntity{}
	var status string
	if err := pool.QueryRow(ctx, sql, id, paymentMethod, amount, entity.Approved).Scan(
		&payment.Id,
		&payment.OrderId,
		&payment.PaymentMethod,
		&payment.Amount,
		&payment.PaidAt,
		&status,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("payment for order %q not found: %w", id, err)
		}
		return nil, fmt.Errorf("pay for order %q: %w", id, err)
	}
	payment.Status = entity.OrderPaymentStatus(status)

	return payment, nil
}

func (r *PostgresPaymentsRepository) UpdateStatus(ctx context.Context, id string, status entity.OrderPaymentStatus) error {
	pool := r.db.Native().(*pgxpool.Pool)
	if _, err := pool.Exec(ctx, `
		UPDATE go_seed.order_payments
		SET status = $2
		WHERE order_id = $1
	`, id, status); err != nil {
		return fmt.Errorf("update payment status for order %q: %w", id, err)
	}
	return nil
}
