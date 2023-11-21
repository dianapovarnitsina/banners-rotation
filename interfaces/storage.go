package interfaces

import "context"

type Storage interface {
	Connect(ctx context.Context, dbPort int, dbHost, dbUser, dbPassword, dbName string) error
	Close(ctx context.Context) error
	Migrate(ctx context.Context, migrate string) error
	AddBanner(ctx context.Context, bannerID, slotID int) error
	RemoveBanner(ctx context.Context, bannerID, slotID int) error
	ClickBanner(ctx context.Context, bannerID, slotID, usergroupID int) error
	PickBanner(ctx context.Context, slotID, usergroupID int) (int, error)
	IsBannerAssignedToSlot(ctx context.Context, bannerID, slotID int) (bool, error)
	BannerExists(ctx context.Context, bannerID int) bool
	SlotExists(ctx context.Context, slotID int) bool
	UserGroupExists(ctx context.Context, userGroupId int) bool
}
