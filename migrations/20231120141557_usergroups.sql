-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS usergroups
(
    id         serial constraint usergroups_pk primary key,
    name       varchar   not null,
    created_at timestamp not null
);

DO $$
BEGIN
    IF (SELECT COUNT(*) FROM usergroups) = 0 THEN
        -- Вставляем данные только если таблица пуста
        INSERT INTO usergroups (name, created_at) VALUES
            ('Groups 1', NOW()),
            ('Groups 2', NOW()),
            ('Groups 3', NOW()),
            ('Groups 4', NOW()),
            ('Groups 5', NOW());
END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS usergroups;
-- +goose StatementEnd
