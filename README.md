# **GO - Microservice**

This repository contains the essential code to build a minimal microservice in **go** <br>
It comes from the Linkedin Certification Course **Build a Microservice with Go** <br>
Thanks to @fpmoles alias Frank P Moley.

## **Prerequisites**

For optimal experience, I recommend installed on your machine (\*NIX flavored OS for sure):

- Go version 1.20 or greater
- PostgreSql
- Docker
- Curl works just fine, but if you look for luxury -> go for **httpie**

## **Get started**

1. Start database - ./postgres_start.sh
   - note: you can make this persistent by specifying volumes in the script such as adding:
   ```
   -e PGDATA=/var/lib/postgresql/data/pgdata \
   -v /home/your_username/directory_of_your_choice_to_persist_data:/var/lib/postgresql/data \
   ```
2. Exec into docker container

   ```
   docker exec -it local-pg /bin/bash
   ```

3. Launch psql from inside the docker container
   ```
   psql -U postgres
   ```
4. Copy/paste the schema file in the **data/** directory directly in the psql console to generate data to work on

5. Start the echo server with **go run main.go** from within the root folder

6. In another terminal, try out the different routes and the five standard REST methods implemented in this repository
