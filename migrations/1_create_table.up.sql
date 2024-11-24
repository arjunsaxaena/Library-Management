CREATE TABLE authors (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE locations (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    standard TEXT NOT NULL
);

CREATE TABLE books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE SET NULL,
    is_checked_out BOOLEAN DEFAULT FALSE,
    book_type TEXT NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW() 
);

CREATE TABLE issued_books (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    issue_date TIMESTAMP NOT NULL DEFAULT NOW(),
    return_date TIMESTAMP,
    late_fees NUMERIC DEFAULT 0
    CHECK (return_date IS NULL OR issue_date <= return_date)
);

CREATE UNIQUE INDEX idx_unique_active_issue ON issued_books (book_id) WHERE return_date IS NULL;

CREATE OR REPLACE VIEW issued_books_with_fees AS
SELECT 
    id,
    book_id,
    user_id,
    issue_date,
    return_date,
    CASE 
        WHEN return_date IS NULL THEN 
            GREATEST(0, (EXTRACT(DAY FROM (CURRENT_DATE - issue_date)) - 15) * 2)
        WHEN EXTRACT(DAY FROM (return_date - issue_date)) > 15 THEN 
            (EXTRACT(DAY FROM (return_date - issue_date)) - 15) * 2
        ELSE 
            0
    END AS late_fees
FROM issued_books;

