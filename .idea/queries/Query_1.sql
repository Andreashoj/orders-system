SELECT products.name, order_products.order_id FROM order_products
    LEFT JOIN products
        ON order_products.product_id = products.id
                    WHERE order_products =