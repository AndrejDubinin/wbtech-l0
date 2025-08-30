package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) InTx(ctx context.Context, f func(tx pgx.Tx) error) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("tx.Rollback: %v", err)
		}
	}(tx, ctx)

	err = f(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) AddOrder(ctx context.Context, order domain.Order) error {
	err := r.InTx(ctx, func(tx pgx.Tx) error {
		var err error

		err = r.addOrder(ctx, order)
		if err != nil {
			return err
		}

		err = r.addDelivery(ctx, order.OrderUID, order.Delivery)
		if err != nil {
			return err
		}

		err = r.addPayment(ctx, order.OrderUID, order.Payment)
		if err != nil {
			return err
		}

		err = r.addItems(ctx, order.OrderUID, order.Items)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *Repository) addOrder(ctx context.Context, order domain.Order) error {
	const query = `
	INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
	delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.conn.Exec(ctx, query, order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID,
		order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) addDelivery(ctx context.Context, orderUUID string, delivery domain.Delivery) error {
	const query = `
	INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.conn.Exec(ctx, query, orderUUID, delivery.Name, delivery.Phone, delivery.Zip,
		delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) addPayment(ctx context.Context, orderUUID string, payment domain.Payment) error {
	const query = `
	INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt,
	bank, delivery_cost, goods_total, custom_fee)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.conn.Exec(ctx, query, orderUUID, payment.Transaction, payment.RequestID,
		payment.Currency, payment.Provider, payment.Amount, payment.PaymentDT, payment.Bank,
		payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) addItems(ctx context.Context, orderUUID string, items []domain.Item) error {
	vals := []any{}
	placeholders := []string{}
	for i, it := range items {
		start := i*12 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			start, start+1, start+2, start+3, start+4, start+5, start+6, start+7, start+8, start+9,
			start+10, start+11))
		vals = append(vals,
			orderUUID, it.ChrtID, it.TrackNumber, it.Price, it.RID, it.Name, it.Sale, it.Size,
			it.TotalPrice, it.NmID, it.Brand, it.Status)
	}

	query := `
	INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale,
	size, total_price, nm_id, brand, status) VALUES ` + strings.Join(placeholders, ",")

	_, err := r.conn.Exec(ctx, query, vals...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetOrders(ctx context.Context, amount int64) ([]*domain.Order, error) {
	const query = `
	SELECT
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
		o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,

		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,

		p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
		p.bank, p.delivery_cost, p.goods_total, p.custom_fee,

		i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale,
		i.size, i.total_price, i.nm_id, i.brand, i.status
	FROM orders o
	INNER JOIN delivery d ON o.order_uid = d.order_uid
	INNER JOIN payment p ON o.order_uid = p.order_uid
	LEFT JOIN items i ON o.order_uid = i.order_uid
	ORDER BY o.date_created DESC
	LIMIT $1
	`
	rows, err := r.conn.Query(ctx, query, amount)
	if err != nil {
		return []*domain.Order{}, err
	}
	defer rows.Close()

	ordersMap := make(map[string]*domain.Order)

	for rows.Next() {
		var order domain.Order
		var delivery domain.Delivery
		var payment domain.Payment
		var item domain.Item

		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID,
			&order.DateCreated, &order.OofShard,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
			&delivery.Region, &delivery.Email,
			&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
			&payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost,
			&payment.GoodsTotal, &payment.CustomFee,
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}

		o, inMap := ordersMap[order.OrderUID]
		if !inMap {
			order.Delivery = delivery
			order.Payment = payment
			order.Items = []domain.Item{}
			ordersMap[order.OrderUID] = &order
			o = &order
		}

		o.Items = append(o.Items, item)
	}

	orders := make([]*domain.Order, 0, len(ordersMap))
	for _, o := range ordersMap {
		orders = append(orders, o)
	}

	return orders, nil
}

func (r *Repository) GetOrder(ctx context.Context, orderUID string) (*domain.Order, error) {
	const query = `
	SELECT
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
		o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,

		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,

		p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
		p.bank, p.delivery_cost, p.goods_total, p.custom_fee,

		i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale,
		i.size, i.total_price, i.nm_id, i.brand, i.status
	FROM orders o
	INNER JOIN delivery d ON o.order_uid = d.order_uid
	INNER JOIN payment p ON o.order_uid = p.order_uid
	LEFT JOIN items i ON o.order_uid = i.order_uid
	WHERE o.order_uid = $1
	`
	rows, err := r.conn.Query(ctx, query, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderResult *domain.Order

	for rows.Next() {
		var order domain.Order
		var delivery domain.Delivery
		var payment domain.Payment
		var item domain.Item

		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID,
			&order.DateCreated, &order.OofShard,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
			&delivery.Region, &delivery.Email,
			&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
			&payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost,
			&payment.GoodsTotal, &payment.CustomFee,
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}

		if orderResult == nil {
			order.Delivery = delivery
			order.Payment = payment
			order.Items = []domain.Item{}
			orderResult = &order
		}

		orderResult.Items = append(orderResult.Items, item)
	}

	return orderResult, nil
}
