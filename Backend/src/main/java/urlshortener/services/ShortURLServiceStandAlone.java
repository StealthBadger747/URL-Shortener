package urlshortener.services;

import org.apache.commons.lang3.RandomStringUtils;

import java.util.concurrent.ConcurrentHashMap;

public class ShortURLServiceStandAlone extends ShortURLService {
    public static final int SHORT_URL_LEN = 6;
    public final ConcurrentHashMap<String, String> shortenedURLsMap;

    public ShortURLServiceStandAlone() {
        this.shortenedURLsMap = new ConcurrentHashMap<>();
    }

    /**
     * Converts the original (long) URL into a small six character string.
     * @param originalURL the original (long) URL.
     * @return a six character string that is a key to the original URL.
     */
    public String createShortURL(String originalURL) {
        // Generate random Alphanumeric string of SHORT_URL_LEN
        String generatedShortcut = RandomStringUtils.randomAlphanumeric(SHORT_URL_LEN);

        // Try and assign the generatedShortcut to the map, but if it exists, regenerate and retry.
        while (shortenedURLsMap.putIfAbsent(generatedShortcut, originalURL) != null) {
            generatedShortcut = RandomStringUtils.randomAlphanumeric(SHORT_URL_LEN);
        }

        return generatedShortcut;
    }

    /**
     * Tries to resolve the shortcut URL to the original URL.
     * @param shortcutURL the shortcut URL to resolve without the leading /.
     * @return the original URL or null if not found.
     */
    public String resolveShortURL(String shortcutURL) {
        if (shortcutURL.length() == SHORT_URL_LEN + 1) {
            return null;
        }

        return this.shortenedURLsMap.get(shortcutURL);
    }
}
