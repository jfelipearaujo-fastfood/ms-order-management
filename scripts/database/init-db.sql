CREATE TABLE IF NOT EXISTS orders (
    id varchar(255),
    customer_id varchar(255),
    track_id varchar(255),
    state int,
    state_updated_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS order_items (
    order_id varchar(255),
    product_id varchar(255),
    quantity int,
    price DECIMAL(10, 2),
    PRIMARY KEY (order_id, product_id)
);

CREATE TABLE IF NOT EXISTS order_payments (
    order_id varchar(255),
    payment_id varchar(255),
    total_items int,
    amount DECIMAL(10, 2),
    state int,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (order_id, payment_id)
);