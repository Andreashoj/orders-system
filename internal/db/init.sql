CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
                                        id UUID PRIMARY KEY,
                                        name VARCHAR(255) NOT NULL,
                                        price INT NOT NULL
);

CREATE TABLE IF NOT EXISTS product_inventory (
                                                 id UUID PRIMARY KEY,
                                                 product_id UUID NOT NULL REFERENCES products(id),
                                                 quantity INT NOT NULL,
                                                 next_shipment_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS carts (
                                     id UUID PRIMARY KEY,
                                     user_id UUID NOT NULL REFERENCES users(id),
                                     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cart_products (
                                             cart_id UUID NOT NULL REFERENCES carts(id),
                                             product_id UUID NOT NULL REFERENCES products(id),
                                             quantity INT NOT NULL,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             PRIMARY KEY (cart_id, product_id)
);

CREATE TABLE IF NOT EXISTS orders (
                                      id UUID PRIMARY KEY,
                                      user_id UUID NOT NULL REFERENCES users(id),
                                      complete BOOLEAN DEFAULT FALSE,
                                      completed_at TIMESTAMP,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_products (
                                              order_id UUID NOT NULL REFERENCES orders(id),
                                              product_id UUID NOT NULL REFERENCES products(id),
                                              quantity INT NOT NULL,
                                              PRIMARY KEY (order_id, product_id)
);

CREATE TABLE IF NOT EXISTS transactions (
                                            id UUID PRIMARY KEY,
                                            order_id UUID NOT NULL REFERENCES orders(id),
                                            payment_type VARCHAR(50) NOT NULL,
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS shipments (
                                         id UUID PRIMARY KEY,
                                         order_id UUID NOT NULL REFERENCES orders(id),
                                         status VARCHAR(50) NOT NULL,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- INSERT INTO products (id, name, price) VALUES
--                                            (gen_random_uuid(), 'Wireless Headphones', 7999),
--                                            (gen_random_uuid(), 'USB-C Cable', 1299),
--                                            (gen_random_uuid(), 'Phone Case', 1999),
--                                            (gen_random_uuid(), 'Screen Protector', 899),
--                                            (gen_random_uuid(), 'Portable Charger', 3499),
--                                            (gen_random_uuid(), 'Keyboard', 5999),
--                                            (gen_random_uuid(), 'Mouse Pad', 499),
--                                            (gen_random_uuid(), 'USB Hub', 2499),
--                                            (gen_random_uuid(), 'Phone Stand', 1599),
--                                            (gen_random_uuid(), 'Laptop Sleeve', 2999);