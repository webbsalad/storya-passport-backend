-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id UUID DEFAULT uuid_generate_v4() NOT NULL CONSTRAINT users_pkey PRIMARY KEY,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    email_id UUID UNIQUE,
    
    CONSTRAINT fk_email FOREIGN KEY (email_id) REFERENCES confirmed_emails(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
