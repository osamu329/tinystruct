package main

import (
	"fmt"
	"text/template"
	"log"
    "github.com/osamu329/tinystruct/parser"
    "github.com/osamu329/tinystruct/syntax"
	"os"
)

func Run(filename string) int {
	file, err := parser.ParseFile(filename)
    if err != nil {
        log.Printf("parse error %s", err)
        return 1
    }

	funcMap := template.FuncMap{
		"isStruct" : func (t syntax.Type) bool {
			if _, ok := t.(*syntax.Struct); ok {
				return ok
			}
			if _, ok := t.(*syntax.Typedef); ok {
				return ok
			}
			return false
		},
		"toStruct" : func (t syntax.Type) *syntax.Struct {
			st, _ := t.(*syntax.Struct)
			return st
		},
		"toArray" : func (t syntax.Type) *syntax.ArrayType{
			st, _ := t.(*syntax.ArrayType)
			return st
		},
		"unpack" : func (st *syntax.Struct, f *syntax.Field, bufname string) string {
			switch ft := f.Type.(type) {
			case *syntax.PrimitiveType:
				switch ft.Name() {
				case "int":
					return fmt.Sprintf("this.%s = %s.getInt()", f.Name, bufname)
				case "char":
					return fmt.Sprintf("this.%s = (char)%s.get()", f.Name, bufname)
				case "uint8_t":
					return fmt.Sprintf("this.%s = (short)%s.get()", f.Name, bufname)
				case "int8_t":
					return fmt.Sprintf("this.%s = %s.get()", f.Name, bufname)
				case "uint16_t":
					return fmt.Sprintf("this.%s = (int)%s.getShort()", f.Name, bufname)
				case "int16_t":
					return fmt.Sprintf("this.%s = %s.getShort()", f.Name, bufname)
				case "uint32_t":
					return fmt.Sprintf("this.%s = (long)%s.getInt()", f.Name, bufname)
				case "int32_t":
					return fmt.Sprintf("this.%s = %s.getShort()", f.Name, bufname)
				}
				return "unpack " + f.Name
			case *syntax.Typedef:
				return fmt.Sprintf("this.%s.unpack(%s)", f.Name, bufname)
			case *syntax.Struct:
				return "struct " + ft.Name()
			case *syntax.ArrayType:
				if ft.IsString() {
					return fmt.Sprintf("this.%s = Structs.unpackString(%s, %d);", f.Name, bufname, ft.Len)
				}
				for _, f0 := range st.Fields {
					if f0.Name == "items_n" || f0.Name == "items_c" {
						return fmt.Sprintf("Structs.unpackArray(%s, this.%s, %s)", bufname, f.Name, f0.Name);
					}
				}
			}
			return "Unknown"
		},
	}
	t, err := template.New("struct").Funcs(funcMap).ParseFiles("struct.java");
	if err != nil {
        log.Printf("load template error %s", err)
        return 1
	}
	for _, ty := range file.TypeList {
		data := struct{
			File *syntax.File
			Typedef *syntax.Typedef
			Package string
		}{file, ty.(*syntax.Typedef), "oapi"}
		fmt.Printf("=== %s ===\n", ty.Name())
		if err = t.ExecuteTemplate(os.Stdout, "struct.java", &data); err != nil {
			log.Printf("execute template error %s", err)
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(Run("struct.h"))
}
