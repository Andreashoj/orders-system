CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(255) NOT NULL,
                                        price INT NOT NULL
);

CREATE TABLE IF NOT EXISTS product_inventory (
                                                 id SERIAL PRIMARY KEY,
                                                 product_id INT NOT NULL REFERENCES products(id),
                                                 quantity INT NOT NULL,
                                                 next_shipment_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS carts (
                                     id SERIAL PRIMARY KEY,
                                     user_id INT NOT NULL REFERENCES users(id),
                                     quantity INT NOT NULL,
                                     last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cart_products (
                                             cart_id INT NOT NULL REFERENCES carts(id),
                                             product_id INT NOT NULL REFERENCES products(id),
                                             quantity INT NOT NULL,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             PRIMARY KEY (cart_id, product_id)
);

CREATE TABLE IF NOT EXISTS orders (
                                      id SERIAL PRIMARY KEY,
                                      user_id INT NOT NULL REFERENCES users(id),
                                      complete BOOLEAN DEFAULT FALSE,
                                      completed_at TIMESTAMP,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_products (
                                              order_id INT NOT NULL REFERENCES orders(id),
                                              product_id INT NOT NULL REFERENCES products(id),
                                              quantity INT NOT NULL,
                                              PRIMARY KEY (order_id, product_id)
);

CREATE TABLE IF NOT EXISTS transactions (
                                            id SERIAL PRIMARY KEY,
                                            order_id INT NOT NULL REFERENCES orders(id),
                                            payment_type VARCHAR(50) NOT NULL,
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS shipments (
                                         id SERIAL PRIMARY KEY,
                                         order_id INT NOT NULL REFERENCES orders(id),
                                         status VARCHAR(50) NOT NULL,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- INSERT INTO products (name, price) VALUES
--                                        ('Wireless Headphones', 7999),
--                                        ('USB-C Cable', 1299),
--                                        ('Phone Case', 1999),
--                                        ('Screen Protector', 899),
--                                        ('Portable Charger', 3499),
--                                        ('Keyboard', 5999),
--                                        ('Mouse Pad', 499),
--                                        ('USB Hub', 2499),
--                                        ('Phone Stand', 1599),
--                                        ('Laptop Sleeve', 2999);