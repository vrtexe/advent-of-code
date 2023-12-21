package util;

import com.google.gson.Gson;

public class JsonSerializer {
    private final Gson gson = new Gson();

    public <T> String serialize(T object) {
        return gson.toJson(object);
    }

    public <T> T deserialize(String json, Class<T> object) {
        return gson.fromJson(json, object);
    }
}
