# fc-session-cache
A simple in-memory cache for session storage

## How it Works
- Your container or environment has an env variable called `FC_SESSION_CACHE_USERNAME` and `FC_SESSION_CACHE_PASSWORD` exposed
- The cache will expose an API to interact with the cache
- Make requests to each endpoint with the format http://username:password@url-to-cache/endpoint

## Endpoints
The cache exposes these API endpoints:
- Get: the Get endpoint takes a `cookie` field in the body and returns the corresponding session
- Put: the Put endpoint takes a `cookie` field and a `session` object and pushes it to the cache
- Remove: the Remove endpoint takes a `cookie` field in the body and removes the corresponding session
- Flush: the Flush endpoint removes all sessions from the cache

## Docker
The cache is also provided as a docker container:
- https://hub.docker.com/r/sthanguy/fc-session-cache
