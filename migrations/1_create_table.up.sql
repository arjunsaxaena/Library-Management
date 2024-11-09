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
    name TEXT NOT NULL
);

CREATE TABLE books (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE SET NULL,
    is_checked_out BOOLEAN DEFAULT FALSE
);

CREATE TABLE issued_books (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    issue_date TIMESTAMP NOT NULL DEFAULT NOW(),
    return_date TIMESTAMP
    CHECK (issue_date <= COALESCE(return_date, NOW()))  -- Ensures return_date is not before issue_date
);

-- Unique index on (book_id) where return_date is NULL to enforce one active issue per book at a time
CREATE UNIQUE INDEX idx_issued_books_book_id ON issued_books (book_id) WHERE return_date IS NULL;