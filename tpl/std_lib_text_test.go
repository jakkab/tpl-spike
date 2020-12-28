package tpl_test

import (
	"github.com/jakkab/tpl-spike/tpl"
	"strconv"
	"testing"
)

var gtb = &tpl.GoBasicText{}

func BenchmarkGoBasicText_Compile_Serial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := gtb.Compile(tgt, td)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoBasicText_Compile_Parallel(b *testing.B) {
	for i := 1; i <= 8; i *= 2 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.SetParallelism(i)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := gtb.Compile(tgt, td)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}