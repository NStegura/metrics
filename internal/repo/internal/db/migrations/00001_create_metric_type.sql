-- +goose Up
-- +goose StatementBegin

BEGIN;
CREATE TYPE metric_types AS ENUM ('gauge', 'counter');
CREATE TABLE metric_type
(
    id          bigserial PRIMARY KEY,
    name        metric_types NOT NULL,
    description text NOT NULL,
    created_at  timestamp NOT NULL DEFAULT NOW()
);
INSERT INTO metric_type
VALUES
    (1, 'gauge', 'A gauge is a metric that represents a single numerical value that can arbitrarily go up and down.'),
    (2, 'counter', 'A counter is a cumulative metric that represents a single monotonically increasing counter.');
COMMIT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS metric_type;
DROP TYPE IF EXISTS metric_types;

-- +goose StatementEnd
