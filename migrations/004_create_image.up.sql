CREATE TABLE IF NOT EXISTS
    images (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
        url TEXT NOT NULL
    );  