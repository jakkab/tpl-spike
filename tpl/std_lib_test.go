package tpl_test

import (
	"github.com/jakkab/tpl-spike/tpl"
	"strconv"
	"testing"
)

const testGoTemplate = `
<td>
	{{.foo}}
</td>
{{range .items}}
<tr class="item">
	<td>
		{{.description}}
	</td>
	
	<td>
		{{.price}}
	</td>
</tr>
{{end}}`

var (
	gc  = &tpl.GoBasic{}
	tgt = []byte(testGoTemplate)
)

func BenchmarkGoBasic_Compile_Serial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := gc.Compile(tgt, td)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoBasic_Compile_Parallel(b *testing.B) {
	for i := 1; i <= 8; i *= 2 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.SetParallelism(i)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := gc.Compile(tgt, td)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}