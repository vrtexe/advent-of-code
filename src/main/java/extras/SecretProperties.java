package extras;

public class SecretProperties {
    public static final String session = "53616c7465645f5faf639e33781208a17cbad71b2af3c1649b90c0029bd926dc718001700af29fd4a7ee814f6e7d7c87d98e1193ceed69a6edaf6706b5db921f";

    public static final String BASE_URL = "https://adventofcode.com";

    public static final String DAY_PREFIX = "day";

    public static final String INPUT_PATH = "input";

    public static String getInputUrl(String year, String day) {
        return "%s/%s/%s/%s".formatted(BASE_URL, year, DAY_PREFIX, day, INPUT_PATH);
    }

    // public static final String INPUT_URL = "%s/%s/%s".formatted(BASE_URL,day,
    // input)

}
