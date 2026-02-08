package urlshortener.server;

import org.eclipse.jetty.server.Handler;
import org.eclipse.jetty.server.handler.HandlerList;
import org.eclipse.jetty.server.handler.ResourceHandler;
import org.eclipse.jetty.servlet.ServletHandler;
import org.eclipse.jetty.servlet.ServletHolder;
import urlshortener.server.servlets.CreateShortenedURLServlet;
import urlshortener.server.servlets.ShortURLResolverServlet;
import urlshortener.services.ShortURLService;

import java.nio.file.Path;

/**
 * This Class is in charge of coordinating the servlets. It acts as the Main Server.
 */
public class MainServer {
    /** The port used for the server */
    private final int port;
    private final Path frontendResources;
    private final ShortURLService urlService;

    /**
     * The constructor for the server
     *
     * @param port the port the sever should start on.
     * @param frontendResources the path to the frontend assets.
     */
    public MainServer(int port, Path frontendResources, ShortURLService urlService) {
        this.port = port;
        this.frontendResources = frontendResources;
        this.urlService = urlService;
    }

    /**
     * Starts the server.
     */
    public void start() {
        System.out.println("Starting Server with port " + this.port);
        System.setProperty("org.eclipse.jetty.LEVEL", "DEBUG");
        org.eclipse.jetty.server.Server server = new org.eclipse.jetty.server.Server(this.port);
        server.setHandler(setupHandlers());
        try {
            server.start();
            server.join();
        } catch (Exception e) {
            System.out.println("Something went terribly wrong in creating the server...");
        }
        System.out.println("Server started successfully");
    }

    /**
     * Sets up the servlet handlers.
     *
     * @return a list of handlers
     */
    public HandlerList setupHandlers() {
        HandlerList handlers = new HandlerList();

        // Set up a resource handler for the frontend static files
        ResourceHandler resourceHandler = new ResourceHandler();
        resourceHandler.setDirectoriesListed(false);
        resourceHandler.setResourceBase(this.frontendResources.toString());

        ServletHandler servlets = new ServletHandler();
        servlets.addServletWithMapping(new ServletHolder(new ShortURLResolverServlet(this.urlService)), "/*");
        servlets.addServletWithMapping(new ServletHolder(new CreateShortenedURLServlet(this.urlService)), "/api/shorten_url");

        handlers.setHandlers(new Handler[] { resourceHandler, servlets });

        return handlers;
    }
}
