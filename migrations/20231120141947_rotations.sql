-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rotations
(
    slot_id    int       not null constraint rotations_slots_id_fk references slots on update cascade on delete cascade,
    banner_id  int       not null constraint rotations_banners_id_fk references banners on update cascade on delete cascade,
    created_at timestamp not null,
    constraint rotations_pk
    primary key (slot_id, banner_id)
);

DO $$
BEGIN
    IF (SELECT COUNT(*) FROM rotations) = 0 THEN
        -- Вставляем данные только если таблица пуста
        INSERT INTO rotations (slot_id, banner_id, created_at) VALUES
            (1, 1, NOW()),
            (2, 2, NOW()),
            (3, 3, NOW()),
            (1, 4, NOW()),
            (2, 5, NOW()),
            (3, 6, NOW());
END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rotations;
-- +goose StatementEnd
