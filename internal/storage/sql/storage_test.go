package sql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

	mock.ExpectExec("INSERT INTO clicks").
		WithArgs(1, 2, 3).
		WillReturnResult(sqlmock.NewResult(0, 1)).
		WillReturnError(nil)

	ctx := context.Background()

	// Тестируем ClickBanner
	err = storage.ClickBanner(ctx, 2, 1, 3)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
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

	mock.ExpectExec("INSERT INTO impressions").
		WithArgs(1, 2, 3).
		WillReturnResult(sqlmock.NewResult(0, 1)).
		WillReturnError(nil)
	ctx := context.Background()

	// Тестируем ImpressBanner
	err = storage.ImpressBanner(ctx, 2, 1, 3)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
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
	expectedSlotID := 1
	expectedUserGroupID := 1

	rows := sqlmock.NewRows([]string{"banner_id", "impressions", "clicks"}).
		AddRow(expectedBannerID, 1, 1)

	mock.ExpectQuery("SELECT").
		WithArgs(expectedUserGroupID, expectedUserGroupID, expectedSlotID).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO impressions").
		WithArgs(expectedSlotID, expectedBannerID, expectedUserGroupID).
		WillReturnResult(sqlmock.NewResult(0, 1)).
		WillReturnError(nil)

	ctx := context.Background()

	// Тестируем PickBanner
	bannerID, err := storage.PickBanner(ctx, expectedSlotID, expectedUserGroupID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if bannerID != expectedBannerID {
		t.Errorf("expected bannerID %d, got %d", expectedBannerID, bannerID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
