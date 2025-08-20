-- Создание таблицы заказов
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) UNIQUE NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    locale VARCHAR(10),
    internal_signature TEXT,
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255),
    shard_key VARCHAR(10),
    sm_id INTEGER,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(10)
);

-- Создание таблицы доставки
CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(255),
    address TEXT,
    region VARCHAR(255),
    email VARCHAR(255),
    UNIQUE(order_uid)
);

-- Создание таблицы платежей
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(10),
    provider VARCHAR(255),
    amount INTEGER,
    payment_dt BIGINT,
    bank VARCHAR(255),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER,
    UNIQUE(order_uid)
);

-- Создание таблицы товаров
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INTEGER,
    track_number VARCHAR(255),
    price INTEGER,
    rid VARCHAR(255),
    name VARCHAR(500),
    sale INTEGER,
    size VARCHAR(50),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(255),
    status INTEGER
);

-- Создание индексов для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_orders_order_uid ON orders(order_uid);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_order_uid ON deliveries(order_uid);
CREATE INDEX IF NOT EXISTS idx_payments_order_uid ON payments(order_uid);
CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
