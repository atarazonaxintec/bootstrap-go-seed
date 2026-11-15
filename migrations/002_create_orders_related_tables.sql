CREATE SCHEMA IF NOT EXISTS go_seed;

CREATE TABLE IF NOT EXISTS go_seed.products (
    sku VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock_quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO go_seed.products (sku, name, description, price, stock_quantity)
VALUES
    ('SKU-001', 'Teclado mecánico', 'Teclado mecánico USB', 49990.00, 25),
    ('SKU-002', 'Mouse inalámbrico', 'Mouse inalámbrico óptico', 14990.00, 40),
    ('SKU-003', 'Monitor 24 pulgadas', 'Monitor LED Full HD', 129990.00, 10),
    ('SKU-004', 'Audífonos Bluetooth', 'Audífonos inalámbricos con micrófono', 34990.00, 18),
    ('SKU-005', 'Webcam HD', 'Cámara web USB con resolución 1080p', 27990.00, 12),
    ('SKU-006', 'Hub USB-C', 'Adaptador multipuerto USB-C de 5 puertos', 19990.00, 30)
ON CONFLICT (sku) DO NOTHING;

CREATE TABLE IF NOT EXISTS go_seed.orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rut VARCHAR(10) NOT NULL,
    sub_total DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS go_seed.order_items (
    order_id UUID NOT NULL REFERENCES go_seed.orders(id),
    sku VARCHAR(100) NOT NULL REFERENCES go_seed.products(sku),
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS go_seed.order_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES go_seed.orders(id),
    status VARCHAR(50) NOT NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS go_seed.order_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES go_seed.orders(id),
    payment_method VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    paid_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL
);
