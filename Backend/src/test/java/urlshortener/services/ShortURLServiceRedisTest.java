package urlshortener.services;

import org.junit.jupiter.api.Test;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPool;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class ShortURLServiceRedisTest {

    @Test
    void createShortUrlPersistsUsingRedis() {
        FakeJedisPool pool = new FakeJedisPool();

        ShortURLServiceRedis service = new ShortURLServiceRedis(pool);
        String shortCode = service.createShortURL("https://example.com");

        assertNotNull(shortCode);
        assertTrue(shortCode.length() == ShortURLServiceRedis.SHORT_URL_LEN);
        assertEquals("https://example.com", pool.getStore().get(shortCode));
    }

    @Test
    void resolveShortUrlReturnsStoredValue() {
        FakeJedisPool pool = new FakeJedisPool();
        pool.getStore().put("abc123", "https://example.com");

        ShortURLServiceRedis service = new ShortURLServiceRedis(pool);
        service.bloomFilter.put("abc123".getBytes());

        String resolved = service.resolveShortURL("abc123");

        assertEquals("https://example.com", resolved);
    }

    @Test
    void resolveShortUrlReturnsNullForNil() {
        FakeJedisPool pool = new FakeJedisPool();

        ShortURLServiceRedis service = new ShortURLServiceRedis(pool);
        service.bloomFilter.put("deadbe".getBytes());

        String resolved = service.resolveShortURL("deadbe");

        assertNull(resolved);
    }

    private static final class FakeJedisPool extends JedisPool {
        private final Map<String, String> store = new ConcurrentHashMap<>();

        @Override
        public Jedis getResource() {
            return new FakeJedis(store);
        }

        public Map<String, String> getStore() {
            return store;
        }
    }

    private static final class FakeJedis extends Jedis {
        private final Map<String, String> store;

        private FakeJedis(Map<String, String> store) {
            this.store = store;
        }

        @Override
        public String set(String key, String value) {
            store.put(key, value);
            return "OK";
        }

        @Override
        public String get(String key) {
            return store.getOrDefault(key, "nil");
        }

        @Override
        public void close() {
            // No-op for fake implementation.
        }
    }
}
