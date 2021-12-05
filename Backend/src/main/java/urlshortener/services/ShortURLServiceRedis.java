package urlshortener.services;

import org.apache.commons.lang3.RandomStringUtils;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;
import urlshortener.util.BloomFilter;

public class ShortURLServiceRedis extends ShortURLService {
    /** Variables for the bloom filter
     * The bloom filter provides concurrency protection to the redis database,
     * assuming only one client is being used at a time.
     *
     * The probability of false positive with 9.2M keys in the filter is 0.02% */
    private static final int BLOOM_M = 80000000;
    private static final int BLOOM_K = 10;
    public final BloomFilter bloomFilter;
    private final Object bloomLock;
    /** Determines how long the generated URL path can be */
    public static final int SHORT_URL_LEN = 6;
    /** Connection to the redis database */
    private JedisPool redisPool;

    public ShortURLServiceRedis(JedisPool redisPool) {
        this.bloomFilter = new BloomFilter(BLOOM_M, BLOOM_K);
        this.bloomLock = new Object();
        this.redisPool = redisPool;
    }

    /**
     * Saves the original and shortcut (k,v pair) to a database/data structure.
     * @param originalURL the original URL that gets saved.
     * @param generatedShortcut the unique shortcut generated for that URL.
     */
    private void saveShortURL(String originalURL, String generatedShortcut) {
        // Add it to the bloomFilter
        this.bloomFilter.put(generatedShortcut.getBytes());
        // Add it to redis
        try (Jedis jedis = this.redisPool.getResource()) {
            jedis.set(generatedShortcut, originalURL);
        }
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

        return generatedShortcut;
    }

    /**
     * Tries to resolve the shortcut URL to the original URL.
     * @param shortcutURL the shortcut URL to resolve without the leading /.
     * @return the original URL or null if not found.
     */
    public String resolveShortURL(String shortcutURL) {
        // Save time by checking bloom filter first
        if (shortcutURL.length() == SHORT_URL_LEN + 1 || !this.bloomFilter.get(shortcutURL.getBytes())) {
            return null;
        }

        // Retrieve the original URL from Redis
        String originalURL;
        try (Jedis jedis = this.redisPool.getResource()) {
            originalURL = jedis.get(shortcutURL);
        }
        if (originalURL.equals("nil")) {
            return null;
        }

        return originalURL;
    }
}
