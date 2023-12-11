-- +goose Up
-- +goose StatementBegin

BEGIN;
CREATE TABLE metric_history
(
    id          bigserial PRIMARY KEY,
    name        text NOT NULL,
    type_id     bigint NOT NULL,
    value       double precision NOT NULL,
    created_at  timestamp NOT NULL DEFAULT NOW(),
    CONSTRAINT FK_metric_history_type FOREIGN KEY(type_id) REFERENCES metric_type(id)
                                                    ON DELETE RESTRICT
                                                    ON UPDATE CASCADE
);
CREATE INDEX idx_created_at ON metric_history(created_at);
CREATE INDEX idx_name ON metric_history(name);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_created_at;
DROP INDEX IF EXISTS idx_name;
DROP TABLE IF EXISTS metric_history;

-- +goose StatementEnd
