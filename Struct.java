package {{.Package}};

public interface Struct<T> {
    int unpack(ByteBuffer src);
    int pack(ByteBuffer dst);
    T newInstance(); 
}
