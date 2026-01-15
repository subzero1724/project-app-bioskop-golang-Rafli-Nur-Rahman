# Quick Setup Guide

This guide will help you set up and run the Cinema Booking System quickly.

## Prerequisites Checklist

- [ ] Go 1.21+ installed (`go version`)
- [ ] PostgreSQL 13+ installed and running
- [ ] Git installed (optional)
- [ ] Postman or similar API testing tool (optional)

## Step-by-Step Setup

### 1. Database Setup

**Create the database:**

```bash
# Login to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE cinema_booking;

# Exit psql
\q
```

**Run migrations:**

```bash
# From project root directory
psql -U postgres -d cinema_booking -f migrations/001_init_schema.sql
```

You should see output indicating tables were created and sample data was inserted.

### 2. Environment Configuration

**Create .env file:**

```bash
cp .env.example .env
```

**Edit .env file with your settings:**

```env
# Application settings
APP_NAME=Cinema Booking System
APP_PORT=8080
APP_ENV=development

# Database settings (UPDATE THESE!)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password  # <-- CHANGE THIS
DB_NAME=cinema_booking
DB_SSLMODE=disable

# JWT settings (CHANGE IN PRODUCTION!)
JWT_SECRET=your-secret-key-change-this-in-production  # <-- CHANGE THIS
JWT_EXPIRATION_HOURS=24

# Logging
LOG_LEVEL=info
LOG_ENCODING=json
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run the Application

**Option A: Using go run**
```bash
go run cmd/main.go
```

**Option B: Using Make**
```bash
make run
```

**Option C: Build and run**
```bash
make build
./bin/cinema-booking-system
```

### 5. Verify Installation

You should see output similar to:

```
{"level":"info","timestamp":"2026-01-14T...","msg":"Starting Cinema Booking System","app_name":"Cinema Booking System","environment":"development","port":"8080"}
{"level":"info","timestamp":"2026-01-14T...","msg":"Database connection established successfully","host":"localhost","database":"cinema_booking"}
{"level":"info","timestamp":"2026-01-14T...","msg":"Server starting","address":":8080"}
```

The server is now running at `http://localhost:8080`

## Testing the API

### Quick Test with cURL

**1. Register a user:**

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'
```

**2. Login:**

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

Copy the `token` from the response.

**3. Get cinemas:**

```bash
curl http://localhost:8080/api/cinemas
```

**4. Check seat availability:**

```bash
curl "http://localhost:8080/api/cinemas/1/seats?date=2026-01-20&time=19:00"
```

**5. Create a booking (replace YOUR_TOKEN):**

```bash
curl -X POST http://localhost:8080/api/booking \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "cinema_id": 1,
    "seat_id": 1,
    "date": "2026-01-20",
    "time": "19:00",
    "payment_method": 1
  }'
```

### Testing with Postman

1. Import `postman_collection.json` into Postman
2. Update the `base_url` variable if needed (default: http://localhost:8080)
3. Run the requests in order:
   - Register User
   - Login (automatically saves token)
   - Protected endpoints will use the saved token

## Sample Data

The database comes pre-populated with:

### Cinemas
1. Cinema XXI Grand Indonesia (Jakarta Pusat) - 150 seats
2. Cinema XXI Pondok Indah (Jakarta Selatan) - 120 seats  
3. CGV Blitz Pacific Place (Jakarta Selatan) - 180 seats
4. Cinepolis Lippo Mall Puri (Jakarta Barat) - 100 seats

### Payment Methods
1. Credit Card
2. Debit Card
3. E-Wallet
4. Bank Transfer
5. Cash

### Seat Types
- **VIP**: Premium seats (front rows in most cinemas)
- **Regular**: Standard seats
- **Junior**: Child-friendly seats (Cinema 4 only)

## Common Issues

### Issue: "Failed to connect to database"

**Solution:**
- Check if PostgreSQL is running: `pg_isready`
- Verify database exists: `psql -U postgres -l | grep cinema_booking`
- Check credentials in `.env` file
- Ensure DB_HOST and DB_PORT are correct

### Issue: "Failed to read config file"

**Solution:**
- Ensure `.env` file exists in project root
- Check file permissions
- Verify all required fields are present

### Issue: "Port already in use"

**Solution:**
- Change APP_PORT in `.env` to a different port (e.g., 8081)
- Or kill the process using port 8080:
  ```bash
  # Find process
  lsof -i :8080
  # Kill it
  kill -9 <PID>
  ```

### Issue: "Token expired"

**Solution:**
- Login again to get a new token
- Tokens expire after 24 hours (configurable in .env)

## Next Steps

1. Read the full [README.md](README.md) for detailed API documentation
2. Explore the codebase to understand the architecture
3. Test all endpoints with Postman
4. Try implementing additional features

## Development Tips

- Use `make fmt` to format code
- Use `make build` to build the binary
- Check logs for detailed error messages
- Use Postman's Console to debug API calls

## Support

For issues or questions:
1. Check the logs for error messages
2. Review the README.md
3. Ensure all prerequisites are met
4. Verify database connection and migrations

---

Happy coding! 🎬🍿
