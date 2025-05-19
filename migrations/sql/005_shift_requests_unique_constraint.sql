-- +goose Up
ALTER TABLE shift_requests ADD CONSTRAINT unique_shift_worker UNIQUE (shift_id, worker_id);

-- +goose Down
ALTER TABLE shift_requests DROP CONSTRAINT unique_shift_worker;
