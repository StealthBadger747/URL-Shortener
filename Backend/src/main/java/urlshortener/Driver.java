package urlshortener;

import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;
import urlshortener.server.MainServer;
import urlshortener.services.ShortURLService;
import urlshortener.services.ShortURLServiceRedis;
import urlshortener.services.ShortURLServiceStandAlone;

import java.nio.file.InvalidPathException;
import java.nio.file.Path;

public class Driver {
    public static final String SERVER_PORT_ENV = "SERVER_PORT";
    public static final String ANGULAR_FRONTEND_DIR_ENV = "ANGULAR_FRONTEND_DIR";
    public static final String USE_REDIS_ENV = "USE_REDIS";
    public static final String REDIS_IP_ENV = "REDIS_IP";
    public static final int REDIS_PORT = 6379;

    public static void main(String[] args) {
        /* Setup defaults */
        int serverPort = 8080;
        ShortURLService shortURLService = new ShortURLServiceStandAlone();
        Path angularFrontendDir;

        // Try and see if there is a different port to use
        if (System.getenv(SERVER_PORT_ENV) != null && !System.getenv(SERVER_PORT_ENV).isBlank()) {
            try {
                serverPort = Integer.parseInt(System.getenv(SERVER_PORT_ENV));
            } catch (NumberFormatException e) {
                System.err.printf("The '"+SERVER_PORT_ENV+"' env variable has invalid formatting with a value of '%s'\n", System.getenv(SERVER_PORT_ENV));
                System.out.printf("Starting the server with the default port of %d instead\n", serverPort);
            }
        }

        // Check which version of the ShortURLService to use (Redis or standalone). Standalone has no persistence.
        if (System.getenv(USE_REDIS_ENV) != null && System.getenv(USE_REDIS_ENV).equals("TRUE")) {
            if (System.getenv(REDIS_IP_ENV) != null && !System.getenv(REDIS_IP_ENV).isBlank()) {
                System.out.printf("Started with redis %s:%s\n", System.getenv(REDIS_IP_ENV), REDIS_PORT);
                JedisPool pool = new JedisPool(new JedisPoolConfig(), System.getenv(REDIS_IP_ENV), REDIS_PORT);
                shortURLService = new ShortURLServiceRedis(pool);
            }
            else {
                System.err.println("Redis IP was not set!");
                System.exit(1);
            }
        }

        // Check if the Path to the Angular Frontend is valid
        if (System.getenv(ANGULAR_FRONTEND_DIR_ENV) != null && !System.getenv(ANGULAR_FRONTEND_DIR_ENV).isBlank()) {
            try {
                angularFrontendDir = Path.of(System.getenv(ANGULAR_FRONTEND_DIR_ENV));
                MainServer server = new MainServer(serverPort, angularFrontendDir, shortURLService);
                server.start();
            } catch (InvalidPathException e) {
                System.err.println("The provided path could not be found!");
                System.exit(1);
            }
        }
        else {
            System.err.println("The front end env variable was blank!");
            System.exit(1);
        }
    }
}
