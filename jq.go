package gq

// #cgo LDFLAGS: -ljq
// #include "jq.h"
import "C"

import "errors"

type JQ struct {
	state *C.struct_jq_state
}

func NewJQ() JQ {
	jq := JQ{}
	jq.Init()
	return jq
}

func (jq *JQ) Init() {
	jq.state = C.jq_init()
}

func (jq *JQ) Compile(program string) error {
	result := C.jq_compile(jq.state, C.CString(program))
	if result == 0 {
		return errors.New("failed to compile")
	} else {
		return nil
	}
}

func (jq *JQ) CompileArgs(program string, jv Jv) error {
	result := C.jq_compile_args(jq.state, C.CString(program), jv.value())
	if result == 0 {
		return errors.New("failed to compile")
	} else {
		return nil
	}
}

func (jq *JQ) Start(jv Jv) {
	C.jq_start(jq.state, jv.v, 0)
}

func (jq *JQ) Iter() chan Jv {
	ch := make(chan Jv)
	go func() {
		for {
			value := C.jq_next(jq.state)
			if C.jv_is_valid(value) == 1 {
				ch <- NewJv(value)
			} else {
				close(ch)
				break
			}
		}
	}()
	return ch
}
func (jq *JQ) Parse(buffer string, program string) (result []Jv, err error) {
	if err = jq.Compile(program); err != nil {
		return
	}

	parser := NewJVParser(0)
	parser.SetBuffer(buffer)

	for jv := range parser.Iter() {
		jq.Start(jv)

		for v := range jq.Iter() {
			result = append(result, v)
		}
	}

	return
}
