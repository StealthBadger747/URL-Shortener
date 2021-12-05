package urlshortener;

import org.apache.commons.lang3.StringUtils;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPoolConfig;
import urlshortener.server.MainServer;
import urlshortener.services.ShortURLService;
import urlshortener.services.ShortURLServiceRedis;
import urlshortener.services.ShortURLServiceStandAlone;

import java.nio.file.InvalidPathException;
import java.nio.file.Path;

/**
 * Class responsible for running this project.
 */
public class Driver {
    public static final String SERVER_PORT_ENV = "SERVER_PORT";
    public static final String ANGULAR_FRONTEND_DIR_ENV = "ANGULAR_FRONTEND_DIR";
    public static final String USE_REDIS_ENV = "USE_REDIS";
    public static final String REDIS_IP_ENV = "REDIS_IP";
    public static final String REDIS_PORT_ENV = "REDIS_PORT";

    /**
     * Initializes the classes necessary based on the provided environment variables.
     * This includes (but is not limited to) the use of redis and the web server port.
     * @param args not used.
     */
    public static void main(String[] args) {
        /* Setup defaults */
        int serverPort = 8080;
        ShortURLService shortURLService = new ShortURLServiceStandAlone();
        Path angularFrontendDir;

        // Try and see if there is a different port to use

        if (System.getenv(SERVER_PORT_ENV) != null && StringUtils.isNotBlank(System.getenv(SERVER_PORT_ENV))) {
            try {
                serverPort = Integer.parseInt(System.getenv(SERVER_PORT_ENV));
            } catch (NumberFormatException e) {
                System.err.printf("The '%s' env variable has invalid formatting with a value of '%s'\n", SERVER_PORT_ENV, System.getenv(SERVER_PORT_ENV));
                System.out.printf("Starting the server with the default port of %d instead\n", serverPort);
            }
        }

        // Check which version of the ShortURLService to use (Redis or standalone). Standalone has no persistence.
        if (System.getenv(USE_REDIS_ENV) != null && System.getenv(USE_REDIS_ENV).equals("TRUE")) {
            if (System.getenv(REDIS_IP_ENV) != null && StringUtils.isNotBlank(System.getenv(REDIS_IP_ENV)) &&
                    System.getenv(REDIS_PORT_ENV) != null && StringUtils.isNotBlank(System.getenv(REDIS_PORT_ENV))) {

                try {
                    int redisPort = Integer.parseInt(System.getenv(REDIS_PORT_ENV));
                    System.out.printf("Started with redis %s:%s\n", System.getenv(REDIS_IP_ENV), redisPort);
                    JedisPool pool = new JedisPool(new JedisPoolConfig(), System.getenv(REDIS_IP_ENV), redisPort);
                    shortURLService = new ShortURLServiceRedis(pool);
                } catch (NumberFormatException e) {
                    System.err.printf("The '%s' env variable has invalid formatting with a value of '%s'\n", REDIS_PORT_ENV, System.getenv(REDIS_PORT_ENV));
                    System.exit(1);
                }
            }
            else {
                System.err.println("Redis IP was not set!");
                System.exit(1);
            }
        }

        // Check if the Path to the Angular Frontend is valid
        if (System.getenv(ANGULAR_FRONTEND_DIR_ENV) != null && StringUtils.isNotBlank(System.getenv(ANGULAR_FRONTEND_DIR_ENV))) {
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
