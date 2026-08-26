package model

import "testing"

func BenchmarkParseJSON(b *testing.B) {
	cases := []struct {
		name string
		data []byte
	}{
		{"wide_object_1000_keys", genWideObjectJSON(1000)},
		{"deep_object_1000_levels", genDeepObjectJSON(1000)},
		{"large_array_100000_items", genLargeArrayJSON(100000)},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := ParseJSON(tc.data); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
