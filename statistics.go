package main

// https://support.sunburst.com/hc/en-us/articles/229335208-Type-to-Learn-How-are-Words-Per-Minute-and-Accuracy-Calculated-
func (m *model) grossWPM() float64 {
	return (float64(m.totalKeysPressed) / 5) / m.elapsedTime.Minutes()
}

// range 0 to 1
func (m *model) accuracy() float64 {
	return (float64(m.correctKeysPressed) / float64(m.totalKeysPressed))
}

func (m *model) adjustedWPM() float64 {
	return m.grossWPM() * m.accuracy()
}
