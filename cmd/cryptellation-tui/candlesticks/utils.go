package candlesticks

func normalizePoint(point, min, max float64, height int) float64 {
	point = (point - min) / (max - min)
	point = point * float64(height)
	return point
}
