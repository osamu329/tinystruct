package {{.Package}};

import java.nio.ByteBuffer;

public final Structs {
    private Structs() {
    }

    public final int unpackArray<T extends Struct>(ByteBuffer src, T[] array, int items_n) {
        int total = 0;
        for (int i = 0; i < items_n; i++) {
            array[i] = T.newStruct();
            array[i].unpack(src);
        }
        return total;
    }

    public final int packArray<T extends Struct>(ByteBuffer dst, T[] array, int items_n) {
        int total = 0;
        for (int i = 0; i < items_n; i++) {
            array[i].pack(dest);
        }
        return total;
    }

    public static final String unpackString(ByteBuffer src, int len) {
        byte[] buf = new byte[len];
        src.get(buf);
        return new String(buf);
    }
}
