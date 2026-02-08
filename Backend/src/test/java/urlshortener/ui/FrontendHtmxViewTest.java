package urlshortener.ui;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.junit.jupiter.api.Test;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class FrontendHtmxViewTest {

    @Test
    void frontendUsesHtmxAndShortenEndpoint() throws IOException {
        Path frontendPath = Path.of("..", "FrontendHtmx", "index.html");
        Document document = Jsoup.parse(new File(frontendPath.toString()), "UTF-8");

        Element htmxScript = document.selectFirst("script[src*=\"htmx.org@2.0.8\"]");
        assertNotNull(htmxScript);

        Element form = document.selectFirst("form[hx-post=\"/api/shorten_url\"]");
        assertNotNull(form);
        assertEquals("#result", form.attr("hx-target"));

        Element resultPanel = document.selectFirst("#result");
        assertNotNull(resultPanel);
        assertTrue(resultPanel.hasAttr("aria-live"));
    }
}
