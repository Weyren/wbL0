package postgres

import (
	"WBL0/internal/common"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderPostgres struct {
	db *pgxpool.Pool
}

func NewOrderPostgres(db *pgxpool.Pool) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func ConnectDB(user, password, host, port, name string) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, name)
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (o *OrderPostgres) CreateOrder(ctx context.Context, order *common.Order) error {
	_, err := o.db.Exec(ctx,
		`INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	_, err = o.db.Exec(
		ctx,
		`INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return err
	}

	_, err = o.db.Exec(
		ctx,
		`INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = o.db.Exec(
			ctx,
			`INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			item.ChrtID, order.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OrderPostgres) GetOrder(ctx context.Context, orderUID string) (common.Order, error) {
	order := common.Order{}

	// Получаем данные заказа
	if err := o.db.QueryRow(ctx, `
		SELECT *
		FROM orders
		WHERE order_uid = $1
	`, orderUID).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated,
		&order.OofShard); err != nil {
		return common.Order{}, fmt.Errorf("error getting order details: %w", err)
	}

	//// Получаем данные доставки
	if err := o.db.QueryRow(ctx, `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1
	`, orderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
	); err != nil {
		return common.Order{}, fmt.Errorf("error getting delivery details: %w", err)
	}

	//// Получаем данные оплаты
	if err := o.db.QueryRow(ctx, `
		SELECT transaction, request_id, currency, provider, amount,
		       payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE transaction = $1
	`, orderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDT,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
	); err != nil {
		return common.Order{}, fmt.Errorf("error getting payment details: %w", err)
	}
	//
	//// Получаем данные товаров
	rows, err := o.db.Query(ctx, `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE track_number = $1
	`, order.TrackNumber)
	if err != nil {
		return common.Order{}, fmt.Errorf("error getting item details: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := common.Item{}
		if err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		); err != nil {
			return common.Order{}, err
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return common.Order{}, err
	}
	return order, nil
}

func (o *OrderPostgres) GetAllOrders(ctx context.Context) ([]common.Order, error) {
	var orders []common.Order
	rows, err := o.db.Query(ctx, `SELECT order_uid FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uid string
		err := rows.Scan(&uid)
		if err != nil {
			return nil, err
		}
		order, err := o.GetOrder(context.TODO(), uid)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
