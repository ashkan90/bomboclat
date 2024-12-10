package generator

import (
	"bytes"
	"fmt"
	"strings"
)

type CodeBuilder struct {
	buf    *bytes.Buffer
	indent int
}

func (b *CodeBuilder) WriteFunc(name string, params []Parameter, returns []Parameter, body func()) {
	b.Printf("func %s(", name)
	b.WriteParams(params)
	b.Printf(") ")
	if len(returns) > 0 {
		b.WriteReturns(returns)
	}
	b.Printf(" {\n")
	b.indent++
	body()
	b.indent--
	b.Printf("}\n\n")
}

func (b *CodeBuilder) Printf(format string, args ...interface{}) {
	fmt.Fprintf(b.buf, strings.Repeat("\t", b.indent)+format, args...)
}

func (b *CodeBuilder) WriteParams(params []Parameter) {
	for i, p := range params {
		if i > 0 {
			b.buf.WriteString(", ")
		}
		if p.Name != "" {
			fmt.Fprintf(b.buf, "%s ", p.Name)
		}
		if p.IsVariadic {
			b.buf.WriteString("...")
		}
		b.buf.WriteString(p.Type)
	}
}

func (b *CodeBuilder) WriteReturns(returns []Parameter) {
	if len(returns) == 0 {
		return
	}

	if len(returns) == 1 && returns[0].Name == "" {
		b.buf.WriteString(returns[0].Type)
		return
	}

	b.buf.WriteString("(")
	for i, r := range returns {
		if i > 0 {
			b.buf.WriteString(", ")
		}
		if r.Name != "" {
			fmt.Fprintf(b.buf, "%s ", r.Name)
		}
		b.buf.WriteString(r.Type)
	}
	b.buf.WriteString(")")
}

func (b *CodeBuilder) WriteType(name string, methods []Method) {
	b.Printf("type %s interface {\n", name)
	b.indent++
	for _, m := range methods {
		b.WriteMethod(m)
	}
	b.indent--
	b.Printf("}\n\n")
}

func (b *CodeBuilder) WriteMethod(m Method) {
	b.Printf("%s(", m.Name)
	b.WriteParams(m.Params)
	b.buf.WriteString(") ")
	b.WriteReturns(m.Returns)
	b.buf.WriteString("\n")
}
