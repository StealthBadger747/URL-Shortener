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

/**
 * This servlet is responsible for creating a short URL.
 */
public class CreateShortenedURLServlet extends HttpServlet {
    private final ShortURLService urlService;
    private final UrlValidator urlValidator;

    public CreateShortenedURLServlet(ShortURLService urlService) {
        super();
        this.urlService = urlService;
        this.urlValidator = new UrlValidator();
    }

    /**
     * Handles the submission/request to shorten a URL.
     * Returns a JSON response back to the client.
     */
    @Override
    protected void doPost(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {

        boolean isHtmxRequest = "true".equalsIgnoreCase(request.getHeader("HX-Request"));
        if (isHtmxRequest) {
            response.setContentType("text/html");
        } else {
            response.setContentType("application/json");
        }
        JSONObject responseJSON = new JSONObject();
        String originalURL = request.getParameter("url");

        // Make sure to escape it to prevent XSS attacks
        // (in case it needs to retrieved/displayed later)
        originalURL = StringEscapeUtils.escapeHtml4(originalURL);

        // Error checking the original URL
        if (originalURL == null || originalURL.isBlank()) {
            response.setStatus(400);
            if (isHtmxRequest) {
                response.getWriter().println("<div class=\"alert error\">Please enter a URL before shortening.</div>");
            } else {
                responseJSON.put("status", "400");
                responseJSON.put("status_message", "URL was blank!");
                response.getWriter().println(responseJSON.toJSONString());
            }
            return;
        }

        // If the URL doesn't already have the http protocol prepended, then prepend it.
        if (!originalURL.startsWith("http")) {
            originalURL = "http://" + originalURL;
        }

        // Don't convert invalid URLs
        if (!this.urlValidator.isValid(originalURL)) {
            response.setStatus(400);
            if (isHtmxRequest) {
                response.getWriter().println("<div class=\"alert error\">That URL doesn't look valid. Check the format and try again.</div>");
            } else {
                responseJSON.put("status", "400");
                responseJSON.put("status_message", "URL was malformed!");
                response.getWriter().println(responseJSON.toJSONString());
            }
            return;
        }

        // Success
        response.setStatus(200);
        String shortCode = this.urlService.createShortURL(originalURL);
        String baseUrl = request.getRequestURL().toString().replace(request.getRequestURI(), "");
        String shortUrl = StringEscapeUtils.escapeHtml4(baseUrl + "/" + shortCode);
        if (isHtmxRequest) {
            response.getWriter().println(
                "<div class=\"result\">" +
                    "<p class=\"result-label\">Short URL</p>" +
                    "<a class=\"result-link\" href=\"" + shortUrl + "\" target=\"_blank\" rel=\"noopener noreferrer\">" +
                        shortUrl +
                    "</a>" +
                "</div>"
            );
        } else {
            responseJSON.put("status", "200");
            responseJSON.put("status_message", "OK");
            responseJSON.put("short_url", shortUrl);
            response.getWriter().println(responseJSON.toJSONString());
        }
    }
}
