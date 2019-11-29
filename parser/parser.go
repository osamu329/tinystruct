package parser

import (
	"github.com/osamu329/tinystruct/syntax"
	"io"
	"log"
	"os"
	"strconv"
	"text/scanner"
)

func ParseFile(filename string) (*syntax.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func Parse(src io.Reader) (*syntax.File, error) {
	p := &parser{}
	p.file = syntax.NewFile()
	p.s.Init(src)
	p.s.Whitespace ^= 1 << '\n'
	p.parse()
	return p.file, nil
}

type parser struct {
	s    scanner.Scanner
	tok  Token
	lit  string
	pos  scanner.Position
	file *syntax.File
}

func (p *parser) parse() {
	tok := p.scan()
	for p.tok != EOF {
		switch tok {
		case TYPEDEF:
			ty := p.typedef()
			p.file.Types[ty.Name] = ty
			continue
		case COMMENT:
		case EOF:
			log.Printf("EOF")
			return
		default:
			log.Printf("%s %s %s", p.pos, tok, p.lit)
			return
		}
		tok = p.scan()
	}
}

func (p *parser) typedef() *syntax.Typedef {
	p.expect(TYPEDEF)
	var typ syntax.Type
	switch p.tok {
	case STRUCT:
		typ = p.structType()
	}
	name := p.ident()
	//log.Printf("typdef: %s", name)
	p.expect(SEMICOLON)
	return &syntax.Typedef{Name: name, Type: typ}
}

func (p *parser) number() int {
	lit := p.lit
	p.expect(INT)
	i, err := strconv.Atoi(lit)
	if err != nil {
		panic(err)
	}
	return i
}

func (p *parser) field() *syntax.Field {
	ty := p.typename()
	name := p.ident()
	for p.got(LBRACK) {
		i := p.number()
		ty = &syntax.ArrayType{Base: ty, Len: i}
		p.expect(RBRACK)
	}
	return &syntax.Field{Name: name, Type: ty}
}

func (p *parser) typename() syntax.Type {
	t, ok := p.findType(p.lit)
	if !ok {
		log.Printf("%s type %s not found", p.pos, p.lit)
		os.Exit(1)
	}
	p.scan()
	return t
}

func (p *parser) structType() *syntax.Struct {
	p.expect(STRUCT)
	name := ""
	if p.tok == IDENT {
		name = p.ident()
	}
	p.expect(LBRACE)
	var fields []*syntax.Field
	for p.tok != RBRACE && p.tok != EOF {
		if p.tok == STRUCT {
			var ty syntax.Type = p.structType()
			name := p.ident()
			for p.got(LBRACK) {
				i := p.number()
				ty = &syntax.ArrayType{Base: ty, Len: i}
				p.expect(RBRACK)
			}
			fields = append(fields, &syntax.Field{Name: name, Type: ty})
		} else {
			fields = append(fields, p.field())
		}
		p.expect(SEMICOLON)
		//log.Printf("%s end field %s", p.pos, p.tok)
	}
	p.expect(RBRACE)
	return &syntax.Struct{Name: name, Fields: fields}
}

func (p *parser) ident() string {
	name := p.lit
	p.expect(IDENT)
	return name
}

func (p *parser) scan() Token {
scan:
	r := p.s.Scan()
	p.lit = p.s.TokenText()
	p.pos = p.s.Position
	//fmt.Printf("%s\n", p.s.Position)
	switch r {
	case scanner.Ident:
		if t, ok := keywords[p.lit]; ok {
			p.tok = t
		} else {
			p.tok = IDENT
		}
	case scanner.Comment:
		p.tok = COMMENT
	case scanner.Int:
		p.tok = INT
	case ';':
		p.tok = SEMICOLON
	case '[':
		p.tok = LBRACK
	case ']':
		p.tok = RBRACK
	case '{':
		p.tok = LBRACE
	case '}':
		p.tok = RBRACE
	case '#':
		p.skipLine()
		goto scan
	case '\n':
		p.tok = LF
		goto scan
	case scanner.EOF:
		p.tok = EOF
	default:
		p.tok = RUNE
	}

	return p.tok
}

func (p *parser) expect(t Token) {
	if p.tok != t {
		log.Printf("%s expect %s but %s", p.pos, t, p.tok)
		os.Exit(1)
	}
	p.scan()
}

func (p *parser) findType(typename string) (syntax.Type, bool) {
	t, ok := p.file.Types[typename]
	return t, ok
}

func (p *parser) got(t Token) bool {
	if p.tok == t {
		p.scan()
		return true
	}
	return false
}

func (p *parser) skipLine() {
	var r rune
	for r = p.s.Peek(); r != '\n' && r != scanner.EOF; r = p.s.Scan() {
	}
	//fmt.Printf("skip %s\t0x%02x\n", p.s.Position, r)
}
