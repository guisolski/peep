package model

import "testing"

// BenchmarkFilterModel_TypeExpression simulates typing a jq expression one
// character at a time into an already-open filter over a large document —
// the hot path that used to re-parse the whole document on every keystroke.
// The FilterModel is constructed once outside the timer so only the
// per-keystroke cost of apply() is measured.
func BenchmarkFilterModel_TypeExpression(b *testing.B) {
	data := genLargeArrayJSON(100_000)
	fm := NewFilterModel(data, 80, 20)
	const expr = ".[50000]"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 1; j <= len(expr); j++ {
			fm.input.SetValue(expr[:j])
			fm.apply()
		}
	}
}
