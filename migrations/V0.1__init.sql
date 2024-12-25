-- Users table
CREATE TABLE
    users (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        email VARCHAR(255) NOT NULL,
        is_email_verified BOOLEAN NOT NULL,
        password VARCHAR(255) NOT NULL,
        role VARCHAR(50) NOT NULL,
        first_name VARCHAR(255),
        last_name VARCHAR(255)
    );

-- Verification Codes table
CREATE TABLE
    verification_codes (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        user_id UUID REFERENCES users (id),
        code VARCHAR(255) NOT NULL,
        purpose VARCHAR(50) NOT NULL
    );

-- Categories table
CREATE TABLE
    categories (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        slug VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        parent_id UUID REFERENCES categories (id)
    );

-- Products table
CREATE TABLE
    products (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        category_id UUID REFERENCES categories (id),
        slug VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        specification TEXT NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        stock INT NOT NULL,
        brand VARCHAR(255) NOT NULL,
        slash_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00
    );

-- Images table
CREATE TABLE
    images (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        product_id UUID REFERENCES products (id),
        key VARCHAR(255) NOT NULL
    );

-- Transactions table
CREATE TABLE
    transactions (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        user_id UUID REFERENCES users (id),
        amount DECIMAL(10, 2) NOT NULL,
        type VARCHAR(50) NOT NULL,
        reference VARCHAR(255) NOT NULL,
        status VARCHAR(50) NOT NULL,
        method VARCHAR(50) NOT NULL,
        description TEXT,
        vendor VARCHAR(255),
        purpose VARCHAR(255)
    );

-- Order Statuses table
CREATE TABLE
    order_statuses (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        name VARCHAR(50) NOT NULL,
        short_name VARCHAR(50) NOT NULL
    );

-- Orders table
CREATE TABLE
    orders (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        user_id UUID REFERENCES users (id),
        transaction_id UUID REFERENCES transactions (id),
        payment_method VARCHAR(50) NOT NULL,
        reference VARCHAR(255) NOT NULL,
        total_price DECIMAL(10, 2) NOT NULL,
        status_id UUID REFERENCES order_statuses (id)
    );

-- Order Items table
CREATE TABLE
    order_items (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        order_id UUID REFERENCES orders (id),
        product_id UUID REFERENCES products (id),
        quantity INT NOT NULL,
        price DECIMAL(10, 2) NOT NULL
    );

-- Order Status Histories table
CREATE TABLE
    order_status_histories (
        id UUID PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        deleted_at TIMESTAMPTZ,
        order_id UUID REFERENCES orders (id),
        status_id UUID REFERENCES order_statuses (id)
    );