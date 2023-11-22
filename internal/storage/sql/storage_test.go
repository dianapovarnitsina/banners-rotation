package sql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	stor "github.com/dianapovarnitsina/banners-rotation/internal/storage"
)

func TestAddBanner(t *testing.T) {
	// Инициализация SQL Mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	storage := &Storage{db: db}

	// Ожидаемый запрос
	mock.ExpectExec("INSERT INTO rotations").
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	ctx := context.Background()

	// Тестируем AddBanner
	err = storage.AddBanner(ctx, 2, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRemoveBanner(t *testing.T) {
	// Инициализация SQL Mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	storage := &Storage{db: db}

	// Ожидаемый запрос
	mock.ExpectExec("DELETE FROM rotations").
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(0, 1)).
		WillReturnError(nil)
	ctx := context.Background()

	// Тестируем RemoveBanner
	err = storage.RemoveBanner(ctx, 2, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestClickBanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	storage := &Storage{db: db}

	expectedClick := &stor.Click{
		ID:          1,
		SlotID:      2,
		BannerID:    3,
		UserGroupID: 4,
		CreatedAt:   time.Now(),
	}

	mock.ExpectQuery("INSERT INTO clicks").
		WithArgs(2, 3, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "slot_id", "banner_id", "usergroup_id", "created_at"}).
			AddRow(expectedClick.ID, expectedClick.SlotID, expectedClick.BannerID, expectedClick.UserGroupID, expectedClick.CreatedAt))

	ctx := context.Background()

	click, err := storage.ClickBanner(ctx, 3, 2, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if click == nil {
		t.Errorf("expected a non-nil click")
		return
	}

	if *click != *expectedClick {
		t.Errorf("unexpected values in Сlick")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestImpressBanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	storage := &Storage{db: db}

	expectedImpress := &stor.Impress{
		ID:          1,
		SlotID:      2,
		BannerID:    3,
		UserGroupID: 4,
		CreatedAt:   time.Now(),
	}

	mock.ExpectQuery("INSERT INTO impressions").
		WithArgs(2, 3, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "slot_id", "banner_id", "usergroup_id", "created_at"}).
			AddRow(expectedImpress.ID, expectedImpress.SlotID, expectedImpress.BannerID, expectedImpress.UserGroupID, expectedImpress.CreatedAt))

	ctx := context.Background()

	impress, err := storage.ImpressBanner(ctx, 3, 2, 1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if impress == nil {
		t.Errorf("expected a non-nil Impress")
		return
	}

	if *impress != *expectedImpress {
		t.Errorf("unexpected values in Impress")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPickBanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	storage := &Storage{db: db}

	expectedBannerID := 1
	expectedSlotID := 2
	expectedUserGroupID := 3

	rows := sqlmock.NewRows([]string{"banner_id", "impressions", "clicks"}).
		AddRow(expectedBannerID, 10, 5) // Example values for simulating a banner

	mock.ExpectQuery("SELECT").
		WithArgs(expectedUserGroupID, expectedSlotID).
		WillReturnRows(rows)

	expectedImpress := &stor.Impress{
		ID:          1,
		SlotID:      expectedSlotID,
		BannerID:    expectedBannerID,
		UserGroupID: expectedUserGroupID,
		CreatedAt:   time.Now(),
	}

	mock.ExpectQuery("INSERT INTO impressions").
		WithArgs(expectedSlotID, expectedBannerID, expectedUserGroupID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "slot_id", "banner_id", "usergroup_id", "created_at"}).
			AddRow(expectedImpress.ID, expectedImpress.SlotID, expectedImpress.BannerID, expectedImpress.UserGroupID, expectedImpress.CreatedAt))

	ctx := context.Background()

	impress, bannerID, err := storage.PickBanner(ctx, expectedSlotID, expectedUserGroupID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if impress == nil {
		t.Errorf("expected a non-nil Impress object")
		return
	}

	if impress.ID != expectedImpress.ID || impress.BannerID != expectedImpress.BannerID {
		t.Errorf("unexpected values in Impress")
		return
	}

	if bannerID != expectedBannerID {
		t.Errorf("unexpected BannerID")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}
