-- +goose Up
CREATE TABLE shift_requests (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    shift_id UUID NOT NULL REFERENCES shifts(id),
    worker_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER shift_requests_updated_at
    BEFORE UPDATE
    ON "shift_requests"
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();


-- +goose Down
DROP TABLE shift_requests;