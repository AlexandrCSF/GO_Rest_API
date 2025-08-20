package sqlstore

import (
	"database/sql"
	"wb_cource/internal/app/model"
)

type OrderRepository struct {
	store *Store
}

func (r *OrderRepository) Create(o *model.Order) error {
	if err := o.Validate(); err != nil {
		return err
	}

	tx, err := r.store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Вставляем основной заказ
	err = tx.QueryRow(
		`INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard,
	).Scan(&o.OrderUID)
	if err != nil {
		return err
	}

	// deliveryJSON, _ := json.Marshal(o.Delivery)
	_, err = tx.Exec(
		`INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email,
	)
	if err != nil {
		return err
	}

	// Вставляем payment
	_, err = tx.Exec(
		`INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider, amount, 
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		o.OrderUID, o.Payment.Transaction, o.Payment.RequestID, o.Payment.Currency,
		o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank,
		o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	// Вставляем items
	for _, item := range o.Items {
		_, err = tx.Exec(
			`INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name, sale, 
				size, total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			o.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) FindByID(orderUID string) (*model.Order, error) {
	o := &model.Order{}

	// Получаем основной заказ
	err := r.store.db.QueryRow(
		`SELECT order_uid, track_number, entry, locale, internal_signature, 
		 customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard 
		 FROM orders WHERE order_uid = $1`,
		orderUID,
	).Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
		&o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	// Получаем delivery
	err = r.store.db.QueryRow(
		`SELECT name, phone, zip, city, address, region, email 
		 FROM deliveries WHERE order_uid = $1`,
		orderUID,
	).Scan(
		&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City,
		&o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Получаем payment
	err = r.store.db.QueryRow(
		`SELECT transaction, request_id, currency, provider, amount, 
		 payment_dt, bank, delivery_cost, goods_total, custom_fee 
		 FROM payments WHERE order_uid = $1`,
		orderUID,
	).Scan(
		&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider,
		&o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank, &o.Payment.DeliveryCost,
		&o.Payment.GoodsTotal, &o.Payment.CustomFee,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Получаем items
	rows, err := r.store.db.Query(
		`SELECT chrt_id, track_number, price, rid, name, sale, 
		 size, total_price, nm_id, brand, status 
		 FROM items WHERE order_uid = $1`,
		orderUID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := model.Item{}
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return o, nil
}

func (r *OrderRepository) GetAll() ([]*model.Order, error) {
	rows, err := r.store.db.Query(`SELECT order_uid FROM orders ORDER BY date_created DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return nil, err
		}

		order, err := r.FindByID(orderUID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
