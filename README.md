# Go URL Shortener

<https://go-shorturl.elghaz.my.id>

A minimalist URL shortener built with Go, featuring Swiss design aesthetics and hexagonal architecture.

## Features

- **URL Shortening**: Convert long URLs into short, shareable links
- **Click Tracking**: Real-time analytics with auto-updating statistics
- **Redis Caching**: Fast URL resolution with Redis cache layer
- **Swiss Design UI**: Clean, minimalist interface with Tailwind CSS
- **HTMX Integration**: Dynamic updates without page reloads

## Tech Stack

- **Backend**: Go 1.25+
- **Database**: PostgreSQL 18
- **Cache**: Redis 7.4
- **Frontend**: HTMX + Tailwind CSS
- **Architecture**: Hexagonal (Ports & Adapters)

## Architecture

```
internal/
├── adapters/          # External interfaces
│   ├── cache/        # Redis implementation
│   ├── handlers/     # HTTP handlers
│   └── repositories/ # Database implementation
├── core/             # Business logic
│   ├── domain/       # Domain models
│   ├── ports/        # Interfaces
│   └── services/     # Business services
└── config/           # Configuration
```

## Prerequisites

- Go 1.25 or higher
- Docker & Docker Compose

## Quick Start

1. **Clone the repository**

```bash
git clone https://github.com/elghazx/go-shorturl.git
cd go-shorturl
```

2. **Set up environment variables**

```bash
cp .env.example .env
```

3. **Start database services**

```bash
docker-compose up -d
```

4. **Install dependencies**

```bash
go mod download
```

5. **Run the application**

```bash
go run cmd/server/main.go
```

6. **Access the application**

```
http://localhost:8080
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Home page |
| POST | `/shorten` | Create short URL |
| GET | `/{shortCode}` | Redirect to original URL |
| GET | `/stats` | Statistics page |
| GET | `/api/stats` | Statistics API (JSON) |

## Development

### Project Structure

- `cmd/server/` - Application entry point
- `internal/` - Internal packages
- `templates/` - HTML templates
- `docker-compose.yml` - Docker services configuration

Services:

- **postgres**: PostgreSQL database
- **redis**: Redis cache

## Features in Detail

### URL Shortening

- Generates 8-character short codes using base64 encoding
- Automatic HTTPS prefix for URLs without protocol
- Duplicate URL handling

### Click Tracking

- Asynchronous click counting with goroutines
- Background context to prevent cancellation
- Real-time statistics updates every 3 seconds

### Caching Strategy

- Cache-first approach for URL resolution
- Automatic cache population on miss
- Redis for high-performance caching

### Swiss Design UI

- Monospace typography (Courier New)
- Black and white color scheme
- Bold, uppercase headings
- Minimal borders and clean layout
- Hover effects for interactive elements

## Author

[elghazx](https://github.com/elghazx)
