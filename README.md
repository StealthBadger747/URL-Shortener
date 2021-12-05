# This is a URL Shortener
## Architecture
- Angular and Bootstrap frontend.
- Jetty for the webserver (both frontend and api hosted there).
- Redis as a database (optional).
  - It was chosen because of my familiarity however there are a couple downsides.
  - It takes a snapshot of the database every 60s which I think is acceptable.
  - It is a in memory database so the amount of KV pairs is limited by RAM which is not ideal for a use case where there might be billions of entries.
  - If I had more time I would have used a database like LevelDB or Cassandra. Or implemented a custom solution using SSTables.
    - SSTables are used in LevelDB, however since I don't need to write over entries in the database ever LevelDB as more functionality than I need (and it isn't easily used with Java). SSTables are great for read heavy workloads from a disk.
- Standalone mode w/o Redis available (no persistence).
  - This mode just uses a ConcurrentHashMap to store the key value pairs. It will not survive a restart or crash.
- The short URLs are generated randomly even for the same input. BaseX encoding schemes could have been used, however I noticed sites like tinyurl.com and bit.ly generate random keys for duplicate entries.

## How to run/build

There are two Bash scripts.
 - The first one is `run_docker.sh` which builds the docker container and then runs it with `docker-compose`. It contains Redis by default as well and the configuration options can be changed in the `docker-compose.yml`.
 - The second one is `run_standalone.sh` which requires npm and maven to be installed on the system. The maven build targets Java 8. The script has environment variables at the top which can be modified.

## Sources/Third Party Libraries:
- Angular for frontend
- pom.xml libraries:
  - MurMur (fast hashing algo) https://sangupta.com/projects/murmur
  - Apache Commons as a helpful library for just about everything.
  - Jedis for a redis client
  - JSON Simple
  - Jetty for the webserver
- For the chain favicon: https://www.favicon-generator.org/search/---/Chain