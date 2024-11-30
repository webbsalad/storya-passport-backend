-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_tokens
(
    device_id UUID DEFAULT uuid_generate_v4() NOT NULL CONSTRAINT user_tokens_pkey PRIMARY KEY,
    user_id UUID NOT NULL,
    version INTEGER DEFAULT 1 NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_tokens;
-- +goose StatementEnd
