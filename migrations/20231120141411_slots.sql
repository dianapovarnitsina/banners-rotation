-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS slots
(
    id         serial constraint slots_pk primary key,
    name       varchar   not null,
    created_at timestamp not null
);

DO $$
BEGIN
    IF (SELECT COUNT(*) FROM slots) = 0 THEN
        -- Вставляем данные только если таблица пуста
        INSERT INTO slots (name, created_at) VALUES
            ('Slot 1', NOW()),
            ('Slot 2', NOW()),
            ('Slot 3', NOW());
END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS slots;
-- +goose StatementEnd
