#!/usr/bin/env python3

import yaml
import os
import sys
import struct

class PrimitiveType:
    def __init__(self, name, struct_symbol):
        self.name = name
        self.struct = struct.Struct(struct_symbol)

    def typecheck(self, value):
        try:
            self.struct.pack(value) 
            return True
        except Exception as e:
            print(f"ERROR value={repr(value)} : {self}({self.struct.format}) {e}")
            return False

    def size(self):
        return self.struct.size

    def pack(self, value):
        return self.struct.pack(value)

    def __repr__(self):
        return self.name


class CharType:
    def __init__(self):
        self.name = "char"
        self.struct = struct.Struct("c")

    def typecheck(self, value):
        return type(value) in (str, bytes) and len(value) == 1

    def size(self):
        return self.struct.size

    def pack(self, value):
        return self.struct.pack(bytes(value, "ascii"))

    def __repr__(self):
        return self.name
            


_defaultTypeDict = {
        name:PrimitiveType(name, sym) for name, sym in [
            ("int8_t", "b"), ("uint8_t", "B"),
            ("int16_t", "h"), ("uint16_t", "H"),
            ("int32_t", "i"), ("uint32_t", "I"),
            ("int64_t", "q"), ("uint64_t", "Q"),
        ]
}
char = CharType()
uint16_t = PrimitiveType("uint16_t", "H")
class BroadcastType:
    def __init__(self):
        self.name = "broadcast_type_t"

    def typecheck(self, value):
        if isinstance(value, str):
            return all([char.typecheck(value[0]),
                char.typecheck(value[1]),
                uint16_t.typecheck(int(value[2:]))])
        if isinstance(value, dict):
            return all((char.typecheck(value["central_module_c"]),
                char.typecheck(value["server_type_c"]),
                uint16_t.typecheck(value["transaction_number_n"])))
        return False

    def pack(self, value):
        if isinstance(value, str):
            return b"".join([char.pack(value[0]),
                char.pack(value[1]),
                uint16_t.pack(int(value[2:]))])
        if isinstance(value, dict):
            return b"".join((char.pack(value["central_module_c"]),
                char.pack(value["server_type_c"]),
                uint16_t.pack(value["transaction_number_n"])))
        
_defaultTypeDict.update({"char":CharType(), "broadcast_type_t": BroadcastType()})


assert(BroadcastType().typecheck("BO16"))

class ArrayType:
    def __init__(self, basetype, size):
        self.name = "%s[%d]" % (basetype.name, size)
        self.basetype = basetype
        self.size = size

    def typecheck(self, value):
        if self.basetype.name == "char":
            if not isinstance(value, str):
                print("value: %s(%s) is not string %s" % (value, type(value), self.name))
                return False
        elif not isinstance(value, list):
            return False
        for e in value:
            self.basetype.typecheck(value)
        if not len(value) < self.size:
            print("value %s (len:%d) exceeds arraysize %s" % (value, len(value), self.name))
            return False
        return True

    def pack(self, value):
        return b"".join(self.basetype.pack(v) for v in value) + (b"\x00" * (self.basetype.size() * (self.size - len(value))))

    def __repr__(self):
        return self.name

class Struct:
    def __init__(self, name, typeobj):
        self.name = name
        self.typeobj = typeobj

    def typecheck(self, obj):
        return all(self.typeobj[field].typecheck(value) for field, value in obj.items())

    def pack(self, obj):
        return b"".join(self.typeobj[field].pack(value) for field, value in obj.items())

    def __repr__(self):
        return "Struct<%s>"%self.name


# TODO: nest type
# TODO: declaration order
def load_types(typedict):
    result = _defaultTypeDict.copy() 
    for typename in typedict:
        #print("load typedef:", typename)
        typeobj = {}
        for fieldname, fieldtype in typedict[typename].items():
            #print("fieldname:", fieldname)
            if fieldtype in result:
                typeobj[fieldname] = result[fieldtype]
                continue
            if fieldtype.count("[") == fieldtype.count("]"):
                basetype, arraysize = fieldtype[:-1].split("[")
                basetype = result[basetype]
                ary = ArrayType(basetype, int(arraysize))
                result[fieldtype] = ary
                typeobj[fieldname] = result[fieldtype]
        result[typename] = Struct(typename, typeobj)
    return result


def load_yaml(path):
    with open(path) as f:
        return yaml.load(f, Loader=yaml.FullLoader)

def load_typedefs(typedefs_path):
    typedefs = load_yaml(typedefs_path)
    typedict = {}
    if "typedefs" not in typedefs:
        print(typedefs.keys())
    for typename, defs in typedefs["typedefs"].items():
        typedict[typename] = defs
    return load_types(typedict)

def load_object(path):
    try:
        objs = load_yaml(path)
        if isinstance(objs, dict):
            objs = objs["data"]
        return objs
    except yaml.parser.ParserError as ex:
        print(ex)
        sys.exit(1)

def typecheck(objs, typedefs):
    for obj in objs:
        print("obj", repr(obj)) # obj : dict
        for typename, value in obj.items():
            if typename not in typedefs:
                print("typename '%s' is not found." % typename)
            if not typedefs[typename].typecheck(value):
                return False

def dump(objs, typedefs):
    l = []
    for obj in objs:
        objtype = typedefs[list(obj.keys())[0]]
        l.append(objtype.pack(list(obj.values())[0]))
    return b"".join(l)

if __name__ == "__main__":
    try:
        types  = load_typedefs("struct.yaml")
        print(types)
        sample = load_object(sys.argv[1])
        typecheck(sample, types)
        print(dump(sample, types))
    except Exception as ex:
        print(ex)
        pass
