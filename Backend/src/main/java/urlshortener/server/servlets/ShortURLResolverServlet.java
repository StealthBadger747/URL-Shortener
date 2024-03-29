package urlshortener.server.servlets;

import urlshortener.services.ShortURLService;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

/**
 * This Servlet is responsible for resolving the short URL into the full
 * original URL. It sends a 301 redirect back to the client after resolving it.
 */
public class ShortURLResolverServlet extends HttpServlet {
    private final ShortURLService urlService;

    public ShortURLResolverServlet(ShortURLService urlService) {
        super();
        this.urlService = urlService;
    }

    /**
     * Handles resolving the short URL into the full original URL.
     */
    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {

        response.setContentType("text/html");
        String shortcutURL = request.getPathInfo();

        // Try to resolve the URL
        String longURL = this.urlService.resolveShortURL(shortcutURL.substring(1));
        if (longURL == null) {
            response.setStatus(404);
            response.getWriter().println("404 NOT FOUND!");
            return;
        }

        // Set the redirect to the resolved URL
        response.setStatus(301);
        response.setHeader("Location", longURL);
        response.getWriter().println("REDIRECTING TO " + longURL);
    }
}
