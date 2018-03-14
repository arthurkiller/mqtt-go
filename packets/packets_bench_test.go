package packets

import (
	"math/rand"
	"testing"
)

func BenchmarkNewControlPacket(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		n := byte(rand.Int()%14 + 1)
		b.StartTimer()
		NewControlPacket(n).Close()
	}
}

func BenchmarkParallelNewControlPacket(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		b.StopTimer()
		n := byte(rand.Int()%14 + 1)
		b.StartTimer()
		for pb.Next() {
			NewControlPacket(n).Close()
		}
	})
}
