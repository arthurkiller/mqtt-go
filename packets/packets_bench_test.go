package packets

import (
	"math/rand"
	"testing"
)

func BenchmarkNewControlPacket(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewControlPacket(byte(rand.Int()%14 + 1)).Close()
	}
}

func BenchmarkParallelNewControlPacket(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NewControlPacket(byte(rand.Int()%14 + 1)).Close()
		}
	})
}
