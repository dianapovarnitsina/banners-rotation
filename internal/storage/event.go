package storage

import "time"

type Click struct {
	ID          int       `json:"id"`
	SlotID      int       `json:"slot_id"`
	BannerID    int       `json:"banner_id"`
	UserGroupID int       `json:"usergroup_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Impress struct {
	ID          int       `json:"id"`
	SlotID      int       `json:"slot_id"`
	BannerID    int       `json:"banner_id"`
	UserGroupID int       `json:"usergroup_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Notification struct {
	TypeEvent   string    `json:"type_event"`   //nolint:tagliatelle
	SlotId      int       `json:"slot_id"`      //nolint:tagliatelle
	BannerId    int       `json:"banner_id"`    //nolint:tagliatelle
	UsergroupId int       `json:"usergroup_id"` //nolint:tagliatelle
	DateTime    time.Time `json:"date_time"`    //nolint:tagliatelle
}
