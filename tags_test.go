package vslparser_test

import (
	"strconv"
	"testing"

	"github.com/Showmax/vslparser"
)

func testTagArray() []vslparser.Tag {
	const (
		uniqueCnt    = 500
		nonUniqueCnt = 500
		mod          = 5
	)
	const totalCnt = uniqueCnt + nonUniqueCnt

	arr := make([]vslparser.Tag, totalCnt)
	for i := 0; i < totalCnt; i++ {
		val := strconv.Itoa(i)

		if i < uniqueCnt {
			arr[i] = vslparser.Tag{Key: val, Value: val}
		} else {
			key := "non-unique-" + strconv.Itoa(i%mod)
			arr[i] = vslparser.Tag{Key: key, Value: val}
		}
	}

	return arr
}

func BenchmarkTags_AllWithKey(b *testing.B) {
	tags := vslparser.Tags(testTagArray())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tags.AllWithKey("non-unique-1")
	}
}
