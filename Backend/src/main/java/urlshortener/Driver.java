package urlshortener;

import urlshortener.server.MainServer;
import urlshortener.services.ShortURLService;

import java.nio.file.Path;

public class Driver {

    public static void main(String[] args) {
        int serverPort = 8080;

        // Try and see if there is a different port to use
        if (System.getenv("SERVER_PORT") != null && !System.getenv("SERVER_PORT").isBlank()) {
            try {
                serverPort = Integer.parseInt(System.getenv("SERVER_PORT"));
            } catch (NumberFormatException e) {
                System.err.printf("The 'SERVER_PORT' env variable has invalid formatting with a value of '%s'\n", System.getenv("SERVER_PORT"));
                System.out.printf("Starting the server with the default port of %d instead\n", serverPort);
            }
        }

        MainServer server = new MainServer(serverPort, Path.of("/Users/erik/Documents/URL-Shortener-Challenge/Frontend/dist/UrlShortener"), new ShortURLService());
        server.start();

    }
}
