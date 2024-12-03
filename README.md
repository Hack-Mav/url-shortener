# URL Shortener

A simple and scalable URL shortener service built with **Go**, **Google Cloud Datastore**, and **Redis** for caching. This service allows users to create shortened URLs, set expiration dates, and track URL usage. It is deployed on **Google App Engine** for reliability and scalability.

---

## Features

- Shorten long URLs with unique short IDs.
- Set expiration dates for short URLs.
- Cache frequently accessed URLs using Redis for faster redirections.
- Deployed on Google App Engine for serverless scalability.

---

## Tech Stack

- **Programming Language**: Go (Golang)
- **Database**: Google Cloud Datastore
- **Cache**: Redis
- **Deployment**: Google App Engine
- **Framework**: Gin (web framework)

---

## Endpoints

### 1. Shorten URL

- **Endpoint**: `/shorten`
- **Method**: `POST`
- **Description**: Shortens a long URL and optionally sets an expiration date.
- **Request Body**:
  ```json
  {
    "long_url": "https://example.com",
    "expiry_date": "2024-12-31" // Optional (YYYY-MM-DD). Defaults to 30 days from creation.
  }
  ```

- **Response Body**:
  ```json
  {
    "short_url": "abc123",
    "expiry_date": "2024-12-31"
  }
  ```

### 2. Redirect URL

- **Endpoint**: `/:short_id`
- **Method**: `GET`
- **Description**: Redirects to the original long URL if the short URL is valid and not expired.
- **Response**:
    - **301 Redirect**: If the short URL is valid.
    - **410 Gone**: If the URL has expired.
    - **404 Not Found**: If the URL does not exist.


