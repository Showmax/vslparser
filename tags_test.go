package vslparser_test

import (
	"strconv"
	"testing"

	"github.com/Showmax/vslparser"
)

func BenchmarkTags_New(b *testing.B) {
	const (
		uniqueCnt    = 50
		nonUniqueCnt = 50
		mod          = 5
	)

	const totalCnt = uniqueCnt + nonUniqueCnt

	tagsArr := make([]vslparser.Tag, totalCnt)
	for i := 0; i < totalCnt; i++ {
		val := strconv.Itoa(i)

		if i < uniqueCnt {
			tagsArr[i] = vslparser.Tag{Key: val, Value: val}
		} else {
			key := "non-unique-" + strconv.Itoa(i%mod)
			tagsArr[i] = vslparser.Tag{Key: key, Value: val}
		}
	}

	arrs := make([][]vslparser.Tag, b.N)
	for i := range arrs {
		arrs[i] = make([]vslparser.Tag, len(tagsArr))
		copy(arrs[i], tagsArr)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vslparser.NewTags(arrs[i])
	}
}
