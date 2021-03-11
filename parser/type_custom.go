package parser

import (
    "fmt"
    "strconv"
    "io"
)

var broadcast_type_t = &TypeInfo {
    name: "broadcast_type_t",
    size: 4,
    typecheck: broadcast_type_typecheck,
    pack: func(w io.Writer, order ByteOrder, v interface{}) error {
        var err error
        if v, ok := v.(map[string]interface{}); ok {
            if err := Char.Pack(w, order, v["central_module_c"]); err != nil {return err}
            if err := Char.Pack(w, order, v["server_type_c"]); err != nil {return err}
            return Uint16.Pack(w, order, v["transaction_number_n"])
        }
        if s, ok := v.(string); ok {
            v, err = broadcast_type_fromString(s)
            if err != nil {
                return nil
            }
            return nil
        }
        return nil
    },
}

func broadcast_type_fromString(s string) (map[string]interface{},error) {
    if len(s) <= 2 {
        return nil, fmt.Errorf("too short")
    }
    if !('A' <= s[0] && s[0] <= 'Z') {
        return nil, fmt.Errorf("not A-Z")
    }
    if !('A' <= s[1] && s[1] <= 'Z') {
        return nil, fmt.Errorf("not A-Z")
    }
    i, err := strconv.Atoi(s[2:])
    if err != nil {
        return nil, err
    }
    if !(0 <= i && i <= 0xffff) {
        return nil, fmt.Errorf("transaction_number_n 0-65535")
    }
    return map[string]interface{}{
        "central_module_c": s[0:1],
        "server_type_c": s[1:2],
        "transaction_number_n": i,
    }, nil
}

func broadcast_type_typecheck(v interface{}) bool {
    switch v := v.(type) {
    case string:
        if dict, err := broadcast_type_fromString(v); err != nil {
            return false
        } else {
            broadcast_type_typecheck(dict)
        }
    case map[string]interface{}:
        if c, ok := v["central_module_c"]; !ok || !Char.TypeCheck(c) {
            return false
        }
        if s, ok := v["server_type_c"]; !ok || !Char.TypeCheck(s) {
            return false
        }
        if t, ok := v["transaction_number_n"]; !ok || !Uint16.TypeCheck(t) {
            return false
        }
        return true
    }
    return false
}
