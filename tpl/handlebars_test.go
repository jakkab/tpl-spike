package tpl_test

import (
	"github.com/jakkab/tpl-spike/tpl"
	"strconv"
	"testing"
)

const testHandlebarsTemplate = `
<td>
	{{foo}}
</td>
{{#each items}}
<tr class="item">
	<td>
		{{description}}
	</td>
	
	<td>
		{{price}}
	</td>
</tr>
{{/each}}`

var (
	hc = &tpl.Handlebars{}
	tht = []byte(testHandlebarsTemplate)
)

func BenchmarkHandlebars_Compile_Serial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := hc.Compile(tht, td)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHandlebars_Compile_Parallel(b *testing.B) {
	for i := 1; i <= 8; i *= 2 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.SetParallelism(i)
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := hc.Compile(tht, td)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

