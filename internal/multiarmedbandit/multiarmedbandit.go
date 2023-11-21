package multiarmedbandit

import (
	"math"
)

type Banner interface {
	GetID() int
	GetImpressions() float64
	GetClicks() float64
}

func PickBanner(banners []Banner) int {
	var (
		totalImpressions float64
		maximumRating    float64 = -1
		selectedBannerID         = 0
	)

	// Находим сумму всех impressions для последующего расчета
	for _, b := range banners {
		imp := b.GetImpressions()
		if imp == 0 {
			imp = 1
		}
		totalImpressions += imp
	}

	// Выбираем баннер с максимальным рейтингом
	for _, b := range banners {
		rating := calculateRating(b.GetClicks(), b.GetImpressions(), totalImpressions)
		if rating > maximumRating {
			maximumRating = rating
			selectedBannerID = b.GetID()
		}
	}

	return selectedBannerID
}

// calculateRating вычисляет рейтинг баннера
func calculateRating(clicks, impressions, totalImpressions float64) float64 {
	if impressions == 0 {
		impressions = 1
	}
	return clicks/impressions + math.Sqrt(2*math.Log(totalImpressions)/impressions)
}
