# Shorty
A URL shortener written in Go

# Routes

## POST shorty/new

```json
{
  "user_id": 1,
  "url": "https://example.com"
}```

## GET shorty/shorturl

Redirects to the long URL
