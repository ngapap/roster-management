-- +goose Up
CREATE TABLE shifts (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    role VARCHAR(100) NOT NULL,
    assigned_to UUID REFERENCES users(id),
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_shift_times CHECK (end_time > start_time)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER shift_updated_at
    BEFORE UPDATE
    ON "shifts"
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose Down
DROP TABLE shifts;