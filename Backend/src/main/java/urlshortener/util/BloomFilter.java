package urlshortener.util;

import com.sangupta.murmur.Murmur3;

import java.util.ArrayList;
import java.util.BitSet;
import java.util.Random;

/**
 * This is useful for determining whether the program has seen something before.
 * It is much more efficient than manually looking at all stored entries.
 */
public class BloomFilter {
    private long numItems;
    private final int numFilterBits;
    private final int numHashFunctions;
    private final long seedNumber;
    private final BitSet filter;

    /**
     * Constructor for the Bloom Filter.
     * @param m the number of bits in the filter.
     * @param k the number of hash functions.
     */
    public BloomFilter(int m, int k) {
        this.numFilterBits = m;
        this.numHashFunctions = k;
        this.filter = new BitSet(m);
        this.seedNumber = new Random().nextLong();
    }

    /**
     * Computes hashes given data.
     * @param data the data to hash.
     * @return An ArrayList of the hashes.
     */
    private ArrayList<Long> computeHashes(byte[] data) {
        ArrayList<Long> hashes = new ArrayList<>();

        long hash1 = Murmur3.hash_x86_32(data, data.length, seedNumber);
        long hash2 = Murmur3.hash_x86_32(data, data.length, hash1);
        hashes.add(hash1);
        hashes.add(hash2);

        for(int i = 2; i < this.numHashFunctions; i++) {
            long derivedHash = (hash1 + i) * hash2;
            hashes.add(derivedHash);
        }

        return hashes;
    }

    /**
     * Calculates the bits to put into the filter.
     * @param data the data we want to "see".
     */
    public void put(byte[] data) {
        BitSet tempSet = new BitSet(this.numFilterBits);
        ArrayList<Long> hashes = computeHashes(data);

        for(long hash : hashes) {
            int index = Math.abs((int) hash % numFilterBits);
            tempSet.set(index, true);
        }
        filter.or(tempSet);
        this.numItems++;
    }

    /**
     * Returns if the given data has been seen before.
     * @param data the data we want to check.
     * @return true or false.
     */
    public boolean get(byte[] data) {
        ArrayList<Long> hashes = computeHashes(data);

        for(long hash : hashes) {
            int index = Math.abs((int) hash % numFilterBits);
            if(this.filter.get(index) == false) {
                return false;
            }
        }
        return true;
    }

    /**
     * Returns the false positive probability for the filter given its current number of elements
     * @return the probability represented as a float.
     */
    public float falsePositiveProb() {
        return (float) Math.pow(1 - Math.exp((this.numHashFunctions * -1) / ((float) this.numFilterBits / this.numItems)),
                this.numHashFunctions);
    }
}