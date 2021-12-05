package urlshortener.services;

import org.apache.commons.lang3.RandomStringUtils;
import urlshortener.util.BloomFilter;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentLinkedDeque;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ShortURLService {
    private static final int BLOOM_M = 8000000;
    private static final int BLOOM_K = 10;
    public static final int SHORT_URL_LEN = 6;
    public final BloomFilter bloomFilter;
    private final Object bloomLock;
    public final ConcurrentHashMap<String, String> shortenedURLsMap;
    private final Pattern shortURLPattern;

    public ShortURLService() {
        this.bloomFilter = new BloomFilter(BLOOM_M, BLOOM_K);
        this.bloomLock = new Object();
        this.shortenedURLsMap = new ConcurrentHashMap<>();
        this.shortURLPattern = Pattern.compile("^/[A-Za-z0-9]{"+SHORT_URL_LEN+"}$");
    }

    /**
     * Saves the original and shortcut (k,v pair) to a database/data structure.
     * @param originalURL the original URL that gets saved.
     * @param generatedShortcut the unique shortcut generated for that URL.
     */
    private void saveShortURL(String originalURL, String generatedShortcut) {
        // Add it to the bloomFilter
        this.bloomFilter.put(generatedShortcut.getBytes());
        // Add it to the Map
        this.shortenedURLsMap.put(generatedShortcut, originalURL);
    }

    /**
     * Converts the original (long) URL into a small six character string.
     * @param originalURL the original (long) URL.
     * @return a six character string that is a key to the original URL.
     */
    public String createShortURL(String originalURL) {
        // Generate random Alphanumeric string of SHORT_URL_LEN
        String generatedShortcut = RandomStringUtils.randomAlphanumeric(SHORT_URL_LEN);

        // If it has been seen before by the bloom filter, then regenerate
        synchronized (this.bloomLock) {
            while (this.bloomFilter.get(generatedShortcut.getBytes())) {
                generatedShortcut = RandomStringUtils.randomAlphanumeric(SHORT_URL_LEN);
            }

            // Save the URL
            this.saveShortURL(originalURL, generatedShortcut);
        }

        //ConcurrentLinkedDeque<String> concurrentLinkedDeque = new ConcurrentLinkedDeque<>();
        //concurrentLinkedDeque.

        return generatedShortcut;
    }

    /**
     * Tries to resolve the shortcut URL to the original URL.
     * @param shortcutURL the shortcut URL to resolve.
     * @return the original URL or null if not found.
     */
    public String resolveShortURL(String shortcutURL) {
        Matcher matcher = shortURLPattern.matcher(shortcutURL);
        if (!matcher.find() || !this.bloomFilter.get(shortcutURL.getBytes())) {
            return null;
        }

        return this.shortenedURLsMap.get(shortcutURL);
    }
}
