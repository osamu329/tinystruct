package {{.Package}};

import java.nio.ByteBuffer;

{{ with .Typedef -}}
public class {{.Name}} implements Struct {
{{- with .BaseType | toStruct }}
{{- $struct := . }}
{{- range .Fields }}
    {{- if .Type.IsArray }}
    {{- $array := toArray .Type }}
    public {{$array.Base.Name}}[] {{.Name}} = new {{$array.Base.Name}}[{{$array.Len}}];
    {{- else }}
    public {{.Type.Name}} {{.Name}}{{if .Type | isStruct }} = new {{.Type.Name}}(){{else}} = 0{{end}};
    {{- end }}
{{- end}}

    @Override
    public void unpack(ByteBuffer src) {
        {{- range .Fields }}
        {{ unpack $struct . "src"}};
        {{- end }}
    }

    @Override
    public void pack(ByteBuffer dest) {
    }
    {{- end}}

}
{{end}}
