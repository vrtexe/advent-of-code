package common;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.stream.Stream;

public class FileLoader {

    private final static ClassLoader classLoader = ClassLoader.getSystemClassLoader();
    private final BufferedReader reader;

    public FileLoader(String name) {
        final var resource = classLoader.getResourceAsStream(name);

        assert resource != null;
        reader = new BufferedReader(new InputStreamReader(resource));
    }

    public Stream<String> loadFile() {
        return reader.lines();
    }

    public void closeStream() {
        try {
            reader.close();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

}
