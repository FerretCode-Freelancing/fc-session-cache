# fc-session-cache
A simple in-memory cache for session storage

## How it Works
- Your container or environment has an env variable called `FC_SESSION_CACHE_USERNAME` and `FC_SESSION_CACHE_PASSWORD` exposed
- The cache will expose an API to interact with the cache
- Make requests to each endpoint with the format http://username:password@url-to-cache/endpoint

## Endpoints
The cache exposes these API endpoints:
- Get
- Put
- Remove
- Flush
