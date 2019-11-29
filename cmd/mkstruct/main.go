package main

import (
    "unicode"
    "log"
    "os"
    "io"
    "bytes"
    "io/ioutil"
    "fmt"
    "flag"
    "strings"
    "go/format"
    "github.com/osamu329/tinystruct/parser"
    "github.com/osamu329/tinystruct/syntax"
)

var (
    filler = flag.Bool("filler", true, "filler as private field")
    bytes2string = flag.Bool("bytes_to_string", false, "[n]byte as string")
)

func ToCamel(s string) string {
    if s == "" {
        panic("")
    }
    if len(s) >= 2 && s[len(s)-2] == '_' {
        s = s[:len(s)-2]
    }
    camel := []rune{unicode.ToUpper(rune(s[0]))}
    for i, c := range s[1:] {
        j := i+ 1
        if s[j] == '_' {
            continue
        } else if s[j-1] == '_' {
            c = unicode.ToUpper(c)
        }
        camel = append(camel, c)
    }
    return string(camel)
}

func GoType(t syntax.Type) string {
    switch t := t.(type) {
    case *syntax.Struct:
        if t.Name() == "" {
            var b bytes.Buffer
            writeStruct(&b, t)
            return b.String()
        } else {
            return t.Name()
        }
    case *syntax.ArrayType:
        base := GoType(t.Base)
        if *bytes2string && base == "byte" {
            return "string"
        }
        return fmt.Sprintf("[%d]%s", t.Len, GoType(t.Base))
    case *syntax.PrimitiveType:
        return t.Name()
    case *syntax.Typedef:
        return ToCamel(t.Name())
    }
    panic(fmt.Sprintf("GoType(%T)", t))
}

func BinTag(t syntax.Type) string {
    switch t := t.(type) {
    case *syntax.Struct:
        return ""
    case *syntax.ArrayType:
        return fmt.Sprintf("// binary:`[%d]%s`", t.Len, t.Base.Name())
    case *syntax.PrimitiveType:
        return fmt.Sprintf("// binary:`%s`", t.Name())
    case *syntax.Typedef:
        return "//"// "+ ToCamel(t.Name())
    }
    panic(fmt.Sprintf("GoType(%T)", t))
}

func writeStruct(w io.Writer, st *syntax.Struct) {
    fmt.Fprintf(w,"struct {\n")
    for _, f := range st.Fields {
        if *filler && strings.HasPrefix(f.Name, "filler_") {
            fmt.Fprintf(w, "%s %s %s\n", f.Name, GoType(f.Type), BinTag(f.Type))
        } else {
            fmt.Fprintf(w, "%s %s %s\n", ToCamel(f.Name), GoType(f.Type), BinTag(f.Type))
        }
    }
    fmt.Fprintf(w, "}")
}

type generator struct {
    w bytes.Buffer
}

func (g *generator) printf(format string, args...interface{}) {
    fmt.Fprintf(&g.w, format, args...)
}


func (g *generator) format() ([]byte, error) {
    return format.Source(g.w.Bytes())

}

func (g *generator) addType(ty syntax.Type) {
    var name string
    switch t := ty.(type) {
    case *syntax.Typedef:
        name = ty.Name()
        ty = t.BaseType()
    }

    st, ok := ty.(*syntax.Struct)
    if !ok {
        return
    }

    g.printf("//struct %s\n", name)
    g.printf("//@ReaderFrom\n")
    g.printf("//@WriterTo\n")
    g.printf("type %s ", ToCamel(name))
    writeStruct(&g.w, st)
    g.printf("\n")
    /*
    _, err := g.format()
    if err != nil {
        os.Stdout.Write(g.w.Bytes())
        panic(err)
    }
    */
}

func Run(filename string) int {
    log.Printf("parse")
    f, err :=  parser.ParseFile(filename)
    if err != nil {
        fmt.Printf("parse error %s", err)
        return 1
    }
    log.Printf("parse done")
    g := generator{}
    g.printf("package main\n")
    g.printf("\n")
    for _, t := range f.TypeList {
        g.addType(t)
    }
    log.Printf("format %d", len(g.w.Bytes()))
    src, err := g.format()
    if err != nil {
        os.Stdout.Write(g.w.Bytes())
        panic(err)
    }
    log.Printf("format done")
    ioutil.WriteFile("struct.go", src, 0666)
    log.Printf("write")
    return 0
}

func main() {
    flag.Parse()

    args := flag.Args()
    if len(args) == 0 {
        fmt.Printf("no input source file")
        os.Exit(1)
    }
    fmt.Printf("%v\n", args)
    os.Exit(Run(args[0]))
}
