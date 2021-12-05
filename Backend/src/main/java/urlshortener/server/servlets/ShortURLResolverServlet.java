package urlshortener.server.servlets;

import urlshortener.services.ShortURLService;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class ShortURLResolverServlet extends HttpServlet {
    private final ShortURLService urlService;

    public ShortURLResolverServlet(ShortURLService urlService) {
        super();
        this.urlService = urlService;
    }

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
        response.sendRedirect(longURL);
        response.setStatus(302);
        response.getWriter().println("REDIRECTING TO " + longURL);
    }
}
