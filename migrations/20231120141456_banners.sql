-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banners
(
    id         serial constraint banners_pk primary key,
    name       varchar   not null,
    created_at timestamp not null
);

DO $$
BEGIN
    IF (SELECT COUNT(*) FROM banners) = 0 THEN
        -- Вставляем данные только если таблица пуста
        INSERT INTO banners (name, created_at) VALUES
            ('Banner 1', NOW()),
            ('Banner 2', NOW()),
            ('Banner 3', NOW()),
            ('Banner 4', NOW()),
            ('Banner 5', NOW()),
            ('Banner 6', NOW()),
            ('Banner 7', NOW()),
            ('Banner 8', NOW()),
            ('Banner 9', NOW()),
            ('Banner 10', NOW());
END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banners;
-- +goose StatementEnd
