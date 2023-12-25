-- +goose Up
-- +goose StatementBegin

BEGIN;
CREATE TABLE metric_actual
(
    id          bigserial PRIMARY KEY,
    name        text UNIQUE NOT NULL,
    type_id     bigint NOT NULL,
    value       double precision NOT NULL,
    created_at  timestamp NOT NULL DEFAULT NOW(),
    updated_at  timestamp NOT NULL DEFAULT NOW(),
    CONSTRAINT FK_metric_actual_type FOREIGN KEY(type_id) REFERENCES metric_type(id)
                                                    ON DELETE RESTRICT
                                                    ON UPDATE CASCADE
);
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS metric_actual;

-- +goose StatementEnd
