package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/aki017/gq"
)

func main() {
	f, _ := os.Create("out.pprof")
	jq := gq.NewJQ()
	program := os.Args[1]
	var text []byte
	if len(os.Args) == 3 {
		fp, _ := os.Open(os.Args[2])
		text, _ = ioutil.ReadAll(fp)
	} else {
		text, _ = ioutil.ReadAll(os.Stdin)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for i := 0; i < 1000; i++ {
		jq.Parse(string(text), program)
	}
	result, _ := jq.Parse(string(text), program)
	for _, r := range result {
		fmt.Println(Dump(r, 0))
	}
}

func indent(depth int) string {
	return strings.Repeat("  ", depth)
}

func DumpString(depth int, v string) string {
	return fmt.Sprintf("%s%s\n", indent(depth), v)
}
func Dump(v gq.Jv, depth int) (out string) {
	fmt.Println(v.RefCount())
	switch v.Kind() {
	case gq.KIND_INVALID:
		out += DumpString(depth, "INVALID")
	case gq.KIND_NULL:
		out += DumpString(depth, "NULL")
	case gq.KIND_TRUE:
		out += DumpString(depth, "TRUE")
	case gq.KIND_FALSE:
		out += DumpString(depth, "FALSE")
	case gq.KIND_NUMBER:
		return fmt.Sprintf("%s%s\n", indent(depth), v.String())
	case gq.KIND_STRING:
		out += DumpString(depth, gq.JvString(v).StringValue())
	case gq.KIND_ARRAY:
		for i, value := range gq.JvArray(v).Array() {
			out += fmt.Sprintf("%s[%d]: ", indent(depth), i)
			if value.Kind() == gq.KIND_OBJECT || value.Kind() == gq.KIND_ARRAY {
				out += "\n" + Dump(value, depth+1)
			} else {
				out += Dump(value, 0)
			}
		}
	case gq.KIND_OBJECT:
		gq.JvObject(v).ForEach(func(key gq.Jv, value gq.Jv) {
			out += fmt.Sprintf("%s%s: ", indent(depth), gq.JvString(key).StringValue())
			if value.Kind() == gq.KIND_OBJECT || value.Kind() == gq.KIND_ARRAY {
				out += "\n" + Dump(value, depth+1)
			} else {
				out += Dump(value, 0)
			}
			value.Free()
		})

	default:
		panic(v)
	}
	return
}
