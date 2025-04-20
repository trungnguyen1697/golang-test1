-- internal/database/migrations/01_init_schema.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    full_name VARCHAR(100),
    role VARCHAR(20) DEFAULT 'user',
    preferences JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
    );

-- Create indexes for user lookup
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Products table
CREATE TABLE IF NOT EXISTS products (
                                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    sale_price DECIMAL(10, 2),
    cost_price DECIMAL(10, 2),
    stock_quantity INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    attributes JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
    );

-- Create indexes for product search and filtering
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_price ON products(price);
CREATE INDEX IF NOT EXISTS idx_products_stock ON products(stock_quantity);
CREATE INDEX IF NOT EXISTS idx_products_attributes ON products USING GIN (attributes);

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
                                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES categories(id),
    is_active BOOLEAN DEFAULT TRUE,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
    );

CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);

-- Product categories junction table
CREATE TABLE IF NOT EXISTS product_categories (
                                                  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
    );

CREATE INDEX IF NOT EXISTS idx_product_categories_product ON product_categories(product_id);
CREATE INDEX IF NOT EXISTS idx_product_categories_category ON product_categories(category_id);

-- Reviews table
CREATE TABLE IF NOT EXISTS reviews (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(100),
    comment TEXT,
    is_verified_purchase BOOLEAN DEFAULT FALSE,
    helpful_votes INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
    );

CREATE INDEX IF NOT EXISTS idx_reviews_product ON reviews(product_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at);

-- Wishlist table
CREATE TABLE IF NOT EXISTS wishlist (
                                        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, product_id)
    );

CREATE INDEX IF NOT EXISTS idx_wishlist_user ON wishlist(user_id);
CREATE INDEX IF NOT EXISTS idx_wishlist_product ON wishlist(product_id);

-- inventory_movements table for tracking stock changes
CREATE TABLE IF NOT EXISTS inventory_movements (
                                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    movement_type VARCHAR(20) NOT NULL, -- purchase, sale, adjustment, return
    reference_id UUID, -- order_id, purchase_id, etc.
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX IF NOT EXISTS idx_inventory_product ON inventory_movements(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movement_type ON inventory_movements(movement_type);
CREATE INDEX IF NOT EXISTS idx_inventory_created_at ON inventory_movements(created_at);

-- Create update_modified_column function for maintaining updated_at timestamps
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to update updated_at timestamps automatically
CREATE TRIGGER update_users_modtime
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_products_modtime
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_categories_modtime
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_reviews_modtime
    BEFORE UPDATE ON reviews
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

-- internal/database/migrations/02_seed_data.sql

-- Insert admin user (password: demo#123)
INSERT INTO users (id, username, email, password_hash, full_name, role, is_active, created_at, updated_at, is_deleted)
VALUES (
           gen_random_uuid(),
           'admin',
           'admin@example.com',
           '$2a$10$AWuR/dHlOdwYY0Vwtez28.thr67ir8LoB964QQr8QS2tX/eYKh8yS', -- bcrypt hash of 'demo#123'
           'System Administrator',
           'admin',
           TRUE,
           CURRENT_TIMESTAMP,
           CURRENT_TIMESTAMP,
              FALSE
       ) ON CONFLICT (username) DO NOTHING;

-- Insert default categories
INSERT INTO categories (id, name, slug, description, is_active, display_order, created_at, updated_at, is_deleted)
VALUES
    (gen_random_uuid(), 'Electronics', 'electronics', 'Electronic devices and gadgets', TRUE, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, FALSE),
    (gen_random_uuid(), 'Computers', 'computers', 'Laptops, desktops, and accessories', TRUE, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, FALSE),
    (gen_random_uuid(), 'Audio', 'audio', 'Headphones, speakers, and audio equipment', TRUE, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, FALSE),
    (gen_random_uuid(), 'Gaming', 'gaming', 'Gaming consoles and accessories', TRUE, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, FALSE),
    (gen_random_uuid(), 'Networking', 'networking', 'Routers, switches, and networking gear', TRUE, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, FALSE)
ON CONFLICT (slug) DO NOTHING;

-- Insert sample products
DO $$
DECLARE
electronics_id UUID;
    computers_id UUID;
    audio_id UUID;
    gaming_id UUID;
    product_id UUID;
BEGIN
    -- Get category IDs
SELECT id INTO electronics_id FROM categories WHERE slug = 'electronics' LIMIT 1;
SELECT id INTO computers_id FROM categories WHERE slug = 'computers' LIMIT 1;
SELECT id INTO audio_id FROM categories WHERE slug = 'audio' LIMIT 1;
SELECT id INTO gaming_id FROM categories WHERE slug = 'gaming' LIMIT 1;

-- Sample product 1
INSERT INTO products (id, sku, name, description, price, sale_price, stock_quantity, status, created_at, updated_at, is_deleted)
VALUES (
           gen_random_uuid(),
           'SW-PRO-001',
           'SmartWatch Pro',
           'Advanced smartwatch with fitness tracking.',
           299.99,
           NULL,
           60,
           'active',
           CURRENT_TIMESTAMP,
           CURRENT_TIMESTAMP,
              FALSE
       ) ON CONFLICT (sku) DO NOTHING
    RETURNING id INTO product_id;

IF product_id IS NOT NULL THEN
        INSERT INTO product_categories (product_id, category_id) VALUES (product_id, electronics_id);
END IF;

    -- Sample product 2
INSERT INTO products (id, sku, name, description, price, sale_price, stock_quantity, status, created_at, updated_at, is_deleted)
VALUES (
           gen_random_uuid(),
           'WM-X-002',
           'Wireless Mouse X',
           'Ergonomic wireless mouse with silent clicks.',
           25.50,
           19.99,
           0,
           'out_of_stock',
           CURRENT_TIMESTAMP,
           CURRENT_TIMESTAMP,
                FALSE
       ) ON CONFLICT (sku) DO NOTHING
    RETURNING id INTO product_id;

IF product_id IS NOT NULL THEN
        INSERT INTO product_categories (product_id, category_id) VALUES (product_id, computers_id);
END IF;

    -- Sample product 3
INSERT INTO products (id, sku, name, description, price, sale_price, stock_quantity, status, created_at, updated_at, is_deleted)
VALUES (
           gen_random_uuid(),
           'UB-AIR-003',
           'UltraBook Air',
           'Lightweight laptop with long battery life.',
           1199.00,
           NULL,
           42,
           'active',
           CURRENT_TIMESTAMP,
           CURRENT_TIMESTAMP,
                FALSE
       ) ON CONFLICT (sku) DO NOTHING
    RETURNING id INTO product_id;

IF product_id IS NOT NULL THEN
        INSERT INTO product_categories (product_id, category_id) VALUES (product_id, computers_id);
END IF;

    -- Sample product 4
INSERT INTO products (id, sku, name, description, price, sale_price, stock_quantity, status, created_at, updated_at, is_deleted)
VALUES (
           gen_random_uuid(),
           'V24-MON-004',
           'Vision 24 Monitor',
           '24-inch Full HD monitor with slim bezel.',
           179.99,
           159.99,
           5,
           'active',
           CURRENT_TIMESTAMP,
           CURRENT_TIMESTAMP,
                FALSE
       ) ON CONFLICT (sku) DO NOTHING
    RETURNING id INTO product_id;

IF product_id IS NOT NULL THEN
        INSERT INTO product_categories (product_id, category_id) VALUES (product_id, computers_id);
END IF;

    -- Sample product 5
INSERT INTO products (id, sku, name, description, price, sale_price, stock_quantity, status, created_at, updated_at, is_deleted)
VALUES (
           gen_random_uuid(),
           'NA-EAR-005',
           'NoiseAway Earbuds',
           'Wireless earbuds with active noise cancellation.',
           79.95,
           NULL,
           89,
           'active',
           CURRENT_TIMESTAMP,
           CURRENT_TIMESTAMP,
                FALSE
       ) ON CONFLICT (sku) DO NOTHING
    RETURNING id INTO product_id;

IF product_id IS NOT NULL THEN
        INSERT INTO product_categories (product_id, category_id) VALUES (product_id, audio_id);
INSERT INTO product_categories (product_id, category_id) VALUES (product_id, electronics_id);
END IF;
END $$;