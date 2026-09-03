CREATE extension IF NOT EXISTS "pgcrypto";

CREATE TYPE order_status AS ENUM (
    'pending', 'paid', 'processing','shipped','delivered','cancelled', 'refunded'
);

CREATE TYPE payment_status AS ENUM ('pending', 'succeeded', 'failed', 'refunded');

CREATE EXTENSION citext;
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name          TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ 
)

CREATE TABLE addresses (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  line1       TEXT NOT NULL,
  line2       TEXT,
  city        TEXT NOT NULL,
  state       TEXT,
  postal_code TEXT NOT NULL,
  country     TEXT NOT NULL,
  is_default  BOOLEAN NOT NULL DEFAULT false,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE categories (
  id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name      TEXT NOT NULL,
  parent_id UUID REFERENCES categories(id)  -- self-referencing, for nested categories
);

CREATE TABLE products (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  category_id UUID REFERENCES categories(id),
  name        TEXT NOT NULL,
  description TEXT,
  base_price  NUMERIC(10,2) NOT NULL CHECK (base_price >= 0),
  is_active   BOOLEAN NOT NULL DEFAULT true,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE product_variants (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  sku        TEXT UNIQUE NOT NULL,
  attributes JSONB NOT NULL DEFAULT '{}',  -- {"size": "M", "color": "Red"}
  price      NUMERIC(10,2) NOT NULL CHECK (price >= 0),
  stock      INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_variants_product ON product_variants(product_id);
CREATE INDEX idx_variants_attributes ON product_variants USING GIN (attributes);

CREATE TABLE carts (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID REFERENCES users(id) ON DELETE CASCADE,  -- nullable: guest carts
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE cart_items (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cart_id    UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
  variant_id UUID NOT NULL REFERENCES product_variants(id),
  quantity   INTEGER NOT NULL CHECK (quantity > 0),
  UNIQUE (cart_id, variant_id)  -- same variant can't appear twice as separate rows
);

CREATE TABLE orders (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           UUID NOT NULL REFERENCES users(id),
  status            order_status NOT NULL DEFAULT 'pending',
  shipping_address  JSONB NOT NULL,   -- snapshot, not a FK — see below
  subtotal          NUMERIC(10,2) NOT NULL CHECK (subtotal >= 0),
  discount_total    NUMERIC(10,2) NOT NULL DEFAULT 0,
  tax_total         NUMERIC(10,2) NOT NULL DEFAULT 0,
  shipping_total    NUMERIC(10,2) NOT NULL DEFAULT 0,
  grand_total       NUMERIC(10,2) NOT NULL CHECK (grand_total >= 0),
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id          UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  variant_id        UUID NOT NULL REFERENCES product_variants(id),
  product_name      TEXT NOT NULL,        -- snapshot
  variant_attributes JSONB NOT NULL,      -- snapshot
  unit_price        NUMERIC(10,2) NOT NULL CHECK (unit_price >= 0),
  quantity          INTEGER NOT NULL CHECK (quantity > 0),
  line_total        NUMERIC(10,2) NOT NULL CHECK (line_total >= 0)
);

CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_order_items_order ON order_items(order_id);

CREATE TABLE payments (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id         UUID NOT NULL REFERENCES orders(id),
  provider         TEXT NOT NULL,           -- 'stripe', 'razorpay', etc.
  provider_ref     TEXT NOT NULL,           -- their transaction/charge id
  status           payment_status NOT NULL DEFAULT 'pending',
  amount           NUMERIC(10,2) NOT NULL CHECK (amount >= 0),
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (provider, provider_ref)
);

CREATE TABLE coupons (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code            TEXT UNIQUE NOT NULL,
  discount_type   TEXT NOT NULL CHECK (discount_type IN ('percent', 'fixed')),
  discount_value  NUMERIC(10,2) NOT NULL,
  max_uses        INTEGER,             -- NULL = unlimited
  times_used      INTEGER NOT NULL DEFAULT 0,
  expires_at      TIMESTAMPTZ,
  is_active       BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE order_coupons (
  order_id  UUID NOT NULL REFERENCES orders(id),
  coupon_id UUID NOT NULL REFERENCES coupons(id),
  discount_amount NUMERIC(10,2) NOT NULL,
  PRIMARY KEY (order_id, coupon_id)
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_products_updated_at
  BEFORE UPDATE ON products
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_orders_updated_at
  BEFORE UPDATE ON orders
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();