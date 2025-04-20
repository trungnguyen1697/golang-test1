-- Drop triggers first
DROP TRIGGER IF EXISTS update_users_modtime ON users;
DROP TRIGGER IF EXISTS update_products_modtime ON products;
DROP TRIGGER IF EXISTS update_categories_modtime ON categories;
DROP TRIGGER IF EXISTS update_reviews_modtime ON reviews;

-- Drop function
DROP FUNCTION IF EXISTS update_modified_column();

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS inventory_movements;
DROP TABLE IF EXISTS wishlist;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

-- Drop extension (optional, comment out if you want to keep it)
-- DROP EXTENSION IF EXISTS "uuid-ossp";