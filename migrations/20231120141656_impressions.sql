-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS impressions
(
    id           serial constraint impressions_pk primary key,
    slot_id      int       not null constraint impressions_slots_id_fk references slots on update cascade on delete cascade,
    banner_id    int       not null constraint impressions_banners_id_fk references banners on update cascade on delete cascade,
    usergroup_id int       not null constraint impressions_usergroups_id_fk references usergroups on update cascade on delete cascade,
    created_at   timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS impressions;
-- +goose StatementEnd
