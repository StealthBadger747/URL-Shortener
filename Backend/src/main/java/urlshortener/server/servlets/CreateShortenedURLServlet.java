package urlshortener.server.servlets;

import org.apache.commons.text.StringEscapeUtils;
import org.apache.commons.validator.routines.UrlValidator;
import org.json.simple.JSONObject;
import urlshortener.services.ShortURLService;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class CreateShortenedURLServlet extends HttpServlet {
    private final ShortURLService urlService;
    private final UrlValidator urlValidator;

    public CreateShortenedURLServlet(ShortURLService urlService) {
        super();
        this.urlService = urlService;
        this.urlValidator = new UrlValidator();
    }

    @Override
    protected void doPost(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {

        response.setContentType("application/json");
        JSONObject responseJSON = new JSONObject();
        String originalURL = request.getParameter("url");

        // Make sure to escape it to prevent XSS attacks
        // (in case it needs to retrieved/displayed later)
        originalURL = StringEscapeUtils.escapeHtml4(originalURL);

        // Error checking the original URL
        if (originalURL == null || originalURL.isBlank()) {
            response.setStatus(400);
            responseJSON.put("status", "400");
            responseJSON.put("status_message", "URL was blank!");
            response.getWriter().println(responseJSON.toJSONString());
            return;
        }

        // If the URL doesn't already have the http protocol prepended, then prepend it.
        if (!originalURL.startsWith("http")) {
            originalURL = "http://" + originalURL;
        }

        // Don't convert invalid URLs
        if (!this.urlValidator.isValid(originalURL)) {
            response.setStatus(400);
            responseJSON.put("status", "400");
            responseJSON.put("status_message", "URL was malformed!");
            response.getWriter().println(responseJSON.toJSONString());
            return;
        }

        // Success
        response.setStatus(200);
        responseJSON.put("status", "200");
        responseJSON.put("status_message", "OK");
        responseJSON.put("short_url", this.urlService.createShortURL(originalURL));
        response.getWriter().println(responseJSON.toJSONString());
    }
}
