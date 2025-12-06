# gradle-resolver-redirector
Run this app on your machine and add following piece in repository section of build.gradle file :
```
maven{
    url 'http://your-machine-ip:10010'
}
```

It will listen to 10010 port by default.

This app will check following repositories for libraries:

> https://dl.google.com/dl/android/maven2

> https://jcenter.bintray.com/

> https://repo1.maven.org/maven2

## Docker

### Build and run with Docker
```bash
docker build -t gradle-resolver-redirector .
docker run -d -p 10010:10010 --name gradle-resolver-redirector gradle-resolver-redirector
```

### Using Docker Compose
```bash
docker-compose up -d
```

To stop the container:
```bash
docker-compose down
```

## Security Features

The application includes several security measures to prevent abuse and attacks:

- **Rate Limiting**: 100 requests per second per IP address (burst of 200)
- **Request Timeouts**: 60-second timeout for request processing
- **Path Validation**: Blocks path traversal attacks (`..`, `//`) and limits path length
- **Method Whitelisting**: Only allows GET and HEAD methods
- **Response Size Limits**: Maximum 50MB response size to prevent memory exhaustion
- **HTTP Client Timeouts**: 30-second timeout for outgoing requests to Maven repositories
- **Connection Limits**: Limits idle connections and connection pool size
- **Security Headers**: Adds X-Frame-Options, X-Content-Type-Options, X-XSS-Protection, and CSP headers
- **Input Validation**: Validates and sanitizes all incoming requests
- **Error Handling**: Prevents information leakage through error messages

