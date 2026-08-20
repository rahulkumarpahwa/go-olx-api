CREATE TYPE listing_status AS ENUM ('active', 'inactive');

CREATE TABLE IF NOT EXISTS
    listings (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        title TEXT NOT NULL,
        description TEXT,
        price BIGINT NOT NULL, -- we are using the BIGINT instead of the NUMREIC as we are storing the values in Paise
        city TEXT NOT NULL,
        status listing_status NOT NULL DEFAULT 'inactive',
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        category_id UUID NOT NULL REFERENCES category (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ
    );