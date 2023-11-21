package multiarmedbandit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type bnr struct {
	ID          int
	impressions int
	clicks      int
}

func (b *bnr) GetID() int {
	return b.ID
}

func (b *bnr) GetImpressions() float64 {
	return float64(b.impressions)
}

func (b *bnr) GetClicks() float64 {
	return float64(b.clicks)
}

func TestPickBanner(t *testing.T) {
	tests := []struct {
		name    string
		banners []Banner
		want    int
	}{
		{
			name: "all banners have no impressions, pick first banner",
			banners: []Banner{
				&bnr{ID: 1, impressions: 0, clicks: 0},
				&bnr{ID: 2, impressions: 0, clicks: 0},
				&bnr{ID: 3, impressions: 0, clicks: 0},
			},
			want: 1,
		},
		{
			name: "one banner has impressions, but not clicks, pick second banner",
			banners: []Banner{
				&bnr{ID: 1, impressions: 2, clicks: 0},
				&bnr{ID: 2, impressions: 0, clicks: 0},
				&bnr{ID: 3, impressions: 0, clicks: 0},
			},
			want: 2,
		},
		{
			name: "two banners have impressions, but not clicks, pick third banner",
			banners: []Banner{
				&bnr{ID: 1, impressions: 2, clicks: 0},
				&bnr{ID: 2, impressions: 2, clicks: 0},
				&bnr{ID: 3, impressions: 0, clicks: 0},
			},
			want: 3,
		},
		{
			name: "all banners have different amount of impressions, but not clicks, pick third banner",
			banners: []Banner{
				&bnr{ID: 1, impressions: 3, clicks: 0},
				&bnr{ID: 2, impressions: 4, clicks: 0},
				&bnr{ID: 3, impressions: 2, clicks: 0},
			},
			want: 3,
		},
		{
			name: "one banner has clicks, pick first banner",
			banners: []Banner{
				&bnr{ID: 1, impressions: 5, clicks: 1},
				&bnr{ID: 2, impressions: 4, clicks: 0},
				&bnr{ID: 3, impressions: 4, clicks: 0},
			},
			want: 1,
		},
		{
			name: "one banner has clicks, but too many impressions, pick second banner",
			banners: []Banner{
				&bnr{ID: 1, impressions: 6, clicks: 1},
				&bnr{ID: 2, impressions: 4, clicks: 0},
				&bnr{ID: 3, impressions: 4, clicks: 0},
			},
			want: 2,
		},
		{
			name: "all banners have clicks, pick second banner",
			banners: []Banner{
				&bnr{ID: 1, impressions: 7, clicks: 2},
				&bnr{ID: 2, impressions: 4, clicks: 1},
				&bnr{ID: 3, impressions: 4, clicks: 0},
			},
			want: 2,
		},
		{
			name: "all banners have clicks, pick third banner",
			banners: []Banner{
				&bnr{ID: 1, impressions: 7, clicks: 2},
				&bnr{ID: 2, impressions: 5, clicks: 1},
				&bnr{ID: 3, impressions: 4, clicks: 1},
			},
			want: 3,
		},
		{
			name: "one banner has many impressions and clicks",
			banners: []Banner{
				&bnr{ID: 1, impressions: 16000, clicks: 799},
				&bnr{ID: 2, impressions: 9000, clicks: 59},
				&bnr{ID: 3, impressions: 3000, clicks: 9},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bannerID := PickBanner(tt.banners)
			require.Equal(t, tt.want, bannerID)
		})
	}
}
