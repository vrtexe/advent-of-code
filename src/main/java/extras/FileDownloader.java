package extras;

import java.io.IOException;
import java.net.CookieHandler;
import java.net.CookieManager;
import java.net.HttpCookie;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;

public class FileDownloader {

    public final HttpClient webClient;
    private final Path filePath;
    private String year;
    private String day;

    public FileDownloader(String file, String year, String day) {
        this.year = file;
        this.day = day;

        createCookie();

        filePath = Path.of(URI.create(file).getPath());
        webClient = HttpClient.newBuilder().cookieHandler(CookieHandler.getDefault())
                .connectTimeout(Duration.ofSeconds(10)).build();

    }

    public final void downloadIfNotExists() throws IOException, InterruptedException {
        if (!Files.exists(filePath)) {
            downloadFile();
        }
    }

    public void downloadFile() throws IOException, InterruptedException {
        var file = webClient
                .send(HttpRequest.newBuilder().uri(URI.create(SecretProperties.getInputUrl(year, day)))
                        .GET().build(), HttpResponse.BodyHandlers.ofByteArray())
                .body();

        Files.write(filePath, file);
    }

    private void createCookie() {
        CookieHandler.setDefault(new CookieManager());

        var cookie = new HttpCookie("session", SecretProperties.session);
        cookie.setPath("/");
        cookie.setVersion(0);

        try {
            if (CookieHandler.getDefault() instanceof CookieManager cookieManager) {
                cookieManager.getCookieStore().add(new URI(SecretProperties.BASE_URL), cookie);
            }
        } catch (Exception e) {
            System.out.println(e);
        }
    }
}
