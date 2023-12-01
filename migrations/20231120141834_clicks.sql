-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS clicks
(
    id           serial constraint clicks_pk primary key,
    slot_id      int       not null constraint clicks_slots_id_fk references slots on update cascade on delete cascade,
    banner_id    int       not null constraint clicks_banners_id_fk references banners on update cascade on delete cascade,
    usergroup_id int       not null constraint clicks_usergroups_id_fk references usergroups on update cascade on delete cascade,
    created_at   timestamp not null
);

DO $$
BEGIN
    IF (SELECT COUNT(*) FROM clicks) = 0 THEN
        -- Вставляем данные только если таблица пуста
        INSERT INTO clicks (slot_id, banner_id, usergroup_id, created_at) VALUES
            (1, 1, 1, NOW()),
            (2, 2, 2, NOW()),
            (3, 3, 3, NOW());
END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clicks;
-- +goose StatementEnd
