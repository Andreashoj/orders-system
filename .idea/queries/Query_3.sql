SELECT carts.id, cart_products.quantity, carts.last_updated, products.id, products.name, products.price FROM carts
    LEFT JOIN cart_products
        ON carts.id = cart_products.cart_id
                LEFT JOIN products
                    ON products.id = cart_products.product_id
                        WHERE user_id = 'a9c05e92-5c1c-4d86-b87b-cf5b43e5f6f7';