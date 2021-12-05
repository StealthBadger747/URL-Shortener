package urlshortener.services;

public abstract class ShortURLService {

    /**
     * Converts the original (long) URL into a small six character string.
     * @param originalURL the original (long) URL.
     * @return a six character string that is a key to the original URL.
     */
    abstract public String createShortURL(String originalURL);

    /**
     * Tries to resolve the shortcut URL to the original URL.
     * @param shortcutURL the shortcut URL to resolve without the leading /.
     * @return the original URL or null if not found.
     */
    abstract public String resolveShortURL(String shortcutURL);

}
