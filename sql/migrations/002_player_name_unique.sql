-- +goose Up
-- +goose StatementBegin
ALTER TABLE player ADD CONSTRAINT UC_player UNIQUE (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE player DROP CONSTRAINT UC_player;
-- +goose StatementEnd
