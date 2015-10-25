package printer

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fatih/hcl/ast"
)

const (
	blank   = byte(' ')
	newline = byte('\n')
	tab     = byte('\t')
)

func (p *printer) printNode(n ast.Node) []byte {
	var buf bytes.Buffer

	switch t := n.(type) {
	case *ast.ObjectList:
		for i, item := range t.Items {
			buf.Write(p.printObjectItem(item))
			if i != len(t.Items)-1 {
				buf.Write([]byte{newline, newline})
			}
		}
	case *ast.ObjectKey:
		buf.WriteString(t.Token.Text)
	case *ast.ObjectItem:
		buf.Write(p.printObjectItem(t))
	case *ast.LiteralType:
		buf.WriteString(t.Token.Text)
	case *ast.ListType:
		buf.Write(p.printList(t))
	case *ast.ObjectType:
		buf.Write(p.printObjectType(t))
	default:
		fmt.Printf(" unknown type: %T\n", n)
	}

	return buf.Bytes()
}

func (p *printer) printObjectItem(o *ast.ObjectItem) []byte {
	var buf bytes.Buffer

	for i, k := range o.Keys {
		buf.WriteString(k.Token.Text)
		buf.WriteByte(blank)

		// reach end of key
		if i == len(o.Keys)-1 {
			buf.WriteString("=")
			buf.WriteByte(blank)
		}
	}

	buf.Write(p.printNode(o.Val))
	return buf.Bytes()
}

func (p *printer) printLiteral(l *ast.LiteralType) []byte {
	return []byte(l.Token.Text)
}

func (p *printer) printObjectType(o *ast.ObjectType) []byte {
	var buf bytes.Buffer
	buf.WriteString("{")
	buf.WriteByte(newline)

	for _, item := range o.List.Items {
		buf.Write(p.indent(p.printObjectItem(item)))
		buf.WriteByte(newline)
	}

	buf.WriteString("}")
	return buf.Bytes()
}

func (p *printer) printList(l *ast.ListType) []byte {
	var buf bytes.Buffer
	buf.WriteString("[")

	for i, item := range l.List {
		if item.Pos().Line != l.Lbrack.Line {
			// not same line
			buf.WriteByte(newline)
		}

		buf.WriteByte(tab)
		buf.Write(p.printNode(item))

		if i != len(l.List)-1 {
			buf.WriteString(",")
			buf.WriteByte(blank)
		} else if item.Pos().Line != l.Lbrack.Line {
			buf.WriteString(",")
			buf.WriteByte(newline)
		}
	}

	buf.WriteString("]")
	return buf.Bytes()
}

func writeBlank(buf io.ByteWriter, indent int) {
	for i := 0; i < indent; i++ {
		buf.WriteByte(blank)
	}
}

func (p *printer) indent(buf []byte) []byte {
	prefix := []byte{tab}
	if p.cfg.SpaceWidth != 0 {
		for i := 0; i < p.cfg.SpaceWidth; i++ {
			prefix = append(prefix, blank)
		}
	}

	var res []byte
	bol := true
	for _, c := range buf {
		if bol && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		bol = c == '\n'
	}
	return res
}
