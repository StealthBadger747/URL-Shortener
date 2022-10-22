# This is a URL Shortener
## Architecture
- Angular and Bootstrap frontend.
- Jetty for the webserver (both frontend and api hosted there).
- Redis as a database (optional).
  - It was chosen because of my familiarity however there are a couple downsides.
  - I also use a Bloom Filter in addition to Redis in order to pre-filter calls to it. This saves time and resources, knowing if the key has already been seen or not.
  - It takes a snapshot of the database every 60s which I think is acceptable.
  - It is a in memory database so the amount of KV pairs is limited by RAM which is not ideal for a use case where there might be billions of entries.
  - If I had more time I would have used a database like LevelDB or Cassandra. Or implemented a custom solution using SSTables.
    - SSTables are used in LevelDB, however since I don't need to write over entries in the database ever, and LevelDB has more functionality than I need (and it isn't easily used with Java). SSTables are great for read heavy workloads from a disk.
- Standalone mode w/o Redis available (no persistence).
  - This mode just uses a ConcurrentHashMap to store the key value pairs. It will not survive a restart or crash.
- The short URLs are generated randomly even for the same input. BaseX encoding schemes could have been used, however I noticed sites like tinyurl.com and bit.ly generate random keys for duplicate entries.

## Note
This project is also hosted on my server in my apartment.
I went a bit overboard with this implementation than was probably expected, but I have been wanting to create a URL Shortener for myself for a while now and this was just a good opportunity. My server environment is exclusively in Docker so it didn't take too long to set up.

## How to run/build
There are two Bash scripts.
 - The first one is `run_docker.sh` which builds the docker container and then runs it with `docker-compose`. It contains Redis by default as well and the configuration options can be changed in the `docker-compose.yml`.
   - If the persistency is tested with the Redis database. Please wait at least 60s after creating the short URL to stop the docker services. This allows Redis to make a snapshot.
 - The second one is `run_standalone.sh` which requires npm and maven to be installed on the system. The maven build targets Java 11+. The script has environment variables at the top which can be modified.
   - Java 11 or greater is required

## Sources/Third Party Libraries:
- Angular for frontend
- Bloom Filter from my Distributed File System project.
- pom.xml libraries:
  - MurMur (fast hashing algo) https://sangupta.com/projects/murmur
  - Apache Commons as a helpful library for just about everything.
  - Jedis for a redis client
  - JSON Simple
  - Jetty for the webserver
- For the chain favicon: https://www.favicon-generator.org/search/---/Chain
