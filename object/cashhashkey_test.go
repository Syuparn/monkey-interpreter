package object

import (
	"fmt"
	"hash/fnv"
	"testing"
)

const N_STRING = 10000

//      [test results]
// TL;DR Hash is faster in every condition!
// goos: windows
// goarch: amd64
// number of strings: 10000
//         |   kind of strings |
//  [us]   |  1 | 10 | 100|1000|
// -----------------------------
//    hash | 112| 206| 233| 575|
// no hash | 891| 802| 816| 732|

func BenchmarkHashKeyWithCash1kind(b *testing.B) {
	strings := sampleStrings(N_STRING, 1)
	b.ResetTimer()
	benchmarkHashKeyWithCash(b, strings)
}

func BenchmarkHashKeyWithCash10kinds(b *testing.B) {
	strings := sampleStrings(N_STRING, 10)
	b.ResetTimer()
	benchmarkHashKeyWithCash(b, strings)
}

func BenchmarkHashKeyWithCash100kinds(b *testing.B) {
	strings := sampleStrings(N_STRING, 100)
	b.ResetTimer()
	benchmarkHashKeyWithCash(b, strings)
}

func BenchmarkHashKeyWithCash1000kinds(b *testing.B) {
	strings := sampleStrings(N_STRING, 1000)
	b.ResetTimer()
	benchmarkHashKeyWithCash(b, strings)
}

func BenchmarkHashKeyWithoutCash1kind(b *testing.B) {
	strings := sampleStrings(N_STRING, 1)
	b.ResetTimer()
	benchmarkHashKeyWithoutCash(b, strings)
}

func BenchmarkHashKeyWithoutCash10kinds(b *testing.B) {
	strings := sampleStrings(N_STRING, 10)
	b.ResetTimer()
	benchmarkHashKeyWithoutCash(b, strings)
}

func BenchmarkHashKeyWithoutCash100kinds(b *testing.B) {
	strings := sampleStrings(N_STRING, 100)
	b.ResetTimer()
	benchmarkHashKeyWithoutCash(b, strings)
}

func BenchmarkHashKeyWithoutCash1000kinds(b *testing.B) {
	strings := sampleStrings(N_STRING, 1000)
	b.ResetTimer()
	benchmarkHashKeyWithoutCash(b, strings)
}

func benchmarkHashKeyWithCash(b *testing.B, strs []string) {
	for i := 0; i < b.N; i++ {
		hashes := make(map[string]HashKey)

		for _, str := range strs {
			if hash, ok := hashes[str]; ok {
				_ = HashKey{Type: STRING_OBJ, Value: hash.Value}
			} else {
				h := fnv.New64a()
				h.Write([]byte(str))
				hashes[str] = HashKey{Type: STRING_OBJ, Value: h.Sum64()}
			}
		}
	}
}

func benchmarkHashKeyWithoutCash(b *testing.B, strs []string) {
	for i := 0; i < b.N; i++ {
		for _, str := range strs {
			h := fnv.New64a()
			h.Write([]byte(str))
			_ = HashKey{Type: STRING_OBJ, Value: h.Sum64()}
		}
	}
}

func sampleStrings(numStrings int, numStringKinds int) []string {
	strings := make([]string, numStrings)

	stringIdx := 0
	for {
		for i := 0; i < numStringKinds; i++ {
			// NOTE: 文字列を冗長にする
			// 短すぎるとハッシュキー生成時間が上手く測定できない可能性があるため
			strings[stringIdx] = fmt.Sprintf("hello, I am %06d", i)
			stringIdx++

			if stringIdx >= numStrings {
				return strings
			}
		}
	}
}
