package parser

import (
  log "github.com/sirupsen/logrus"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)
// typedef.yaml のパーサー

type rootNode struct {
  Version int `yaml:"version"`
  Typedefs []*typedefNode `yaml:"typedefs"`
}

type typedefNode struct {
  Struct structNode `yaml:"struct"`
}

type structNode struct {
  Name string `yaml:"name"`
  Fields map[string]string`yaml:"fields"`
  Attributes map[string]interface{} `yaml:"attributes"`
}

type typedefParser struct {
  typemap map[string]*Type
}

func (p *typedefParser) parseTypedef(rootNode *rootNode) error {
  for _, typedef := range rootNode.Typedefs {
    var fields []*Field
    for fieldname, fieldtype := range typedef.Struct.Fields {
      fields = append(fields, &Field{Name:fieldname, Type:p.resolveType(fieldtype)})
    }
    st := NewStruct(typedef.Struct.Name, fields)
    p.typemap[st.Name()] = &Type{Name:st.Name(), Extra: st}
  }
  return nil
}

func (p *typedefParser) resolveType(typename string) *Type {
  if t, ok := p.typemap[typename]; ok {
    return t
  }
  log.WithFields(log.Fields{"typename":typename}).Info("type not resolved")
  t := &Type{Name:typename}
  p.typemap[typename] = t
  return t
}

func ParseTypedef(path string) (*rootNode, error) {
  yaml_data, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  var rootNode rootNode
  err = yaml.Unmarshal(yaml_data, &rootNode)
  if err != nil {
    return nil, err
  }
  p:= &typedefParser{make(map[string]*Type)}
  for _, t := range primitiveTypes {
    p.typemap[t.Name] = t
  }
  p.parseTypedef(&rootNode)
  return &rootNode, err
}


