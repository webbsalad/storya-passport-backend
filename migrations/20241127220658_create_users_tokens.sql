-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_tokens
(
    id UUID DEFAULT uuid_generate_v4() NOT NULL CONSTRAINT user_tokens_pkey PRIMARY KEY,
    user_id UUID NOT NULL,
    device_id UUID NOT NULL,
    version INTEGER NOT NULL,
    CONSTRAINT unique_user_device UNIQUE (user_id, device_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_tokens;
-- +goose StatementEnd
