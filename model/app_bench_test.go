package model

import "testing"

// BenchmarkNewApp measures the cost of opening a large document — the
// point at which RawModel/SearchModel/FilterModel used to be built eagerly
// (each doing a full-document pass: Indent, FlatAll, Indent again) even
// though most sessions only ever use the tree/graph views.
func BenchmarkNewApp(b *testing.B) {
	data := genLargeArrayJSON(100_000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := NewApp(data, "bench.json", 120, 40); err != nil {
			b.Fatalf("NewApp error: %v", err)
		}
	}
}
