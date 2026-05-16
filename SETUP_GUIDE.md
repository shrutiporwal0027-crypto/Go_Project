# Quick Setup Guide - Booking System

## Prerequisites
- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

---

## Step-by-Step Setup

### Step 1: Clone Repository
```bash
git clone <repository-url>
cd booking-system
```

### Step 2: Download Dependencies
```bash
go mod download
```

### Step 3: Create PostgreSQL Database
```bash
# macOS with Homebrew
brew services start postgresql
createdb booking_system

# Linux
sudo systemctl start postgresql
sudo -u postgres createdb booking_system

# Windows (using PowerShell)
.\pg_ctl.exe -D "C:\Program Files\PostgreSQL\data" start
createdb -U postgres booking_system
```

### Step 4: Configure Environment
```bash
# Copy example configuration
cp .env.example .env

# Edit .env with your settings
# Required:
# - DB_PASSWORD
# - JWT_SECRET
# - SMTP credentials (for email)

# Using nano editor
nano .env

# Or your preferred editor
code .env
```

### Step 5: Run Application
```bash
# Start the server
go run main.go *.go

# Or build and run
go build -o booking-system
./booking-system
```

### Step 6: Verify Installation
```bash
# Open browser and navigate to
http://localhost:8080

# Test API
curl http://localhost:8080/api/v1/bookings

# Should see JSON response
{"bookings": [], "total": 0, "remaining": 50}
```

---

## Common Issues & Solutions

### Issue: "Database connection refused"
```
Error: failed to connect to database

Solution:
1. Verify PostgreSQL is running
   - macOS: brew services list | grep postgres
   - Linux: sudo systemctl status postgresql
   - Windows: services.msc

2. Check credentials in .env
   - Default user is "postgres"
   - Default password is empty or set during install

3. Verify database exists
   - psql -U postgres -l
   - Should show "booking_system" in list
```

### Issue: "Port 8080 already in use"
```
Error: listen tcp :8080: bind: address already in use

Solution:
1. Find process using port 8080
   - Linux/macOS: lsof -i :8080
   - Windows: netstat -ano | findstr :8080

2. Kill the process
   - Linux/macOS: kill -9 <PID>
   - Windows: taskkill /PID <PID> /F

3. Or change port in code/env (search for :8080)
```

### Issue: "Email not sending"
```
Error: failed to send email

Solution:
1. If using Gmail:
   - Enable 2-factor authentication
   - Generate app-specific password
   - Use app password in .env, NOT account password
   
2. Check SMTP credentials
   - SMTP_HOST should be: smtp.gmail.com
   - SMTP_PORT should be: 587
   - SMTP_EMAIL should be: your-email@gmail.com

3. Temporarily disable email for testing
   - Comment out SendEmail() calls
```

### Issue: "QR code generation failed"
```
Error: failed to generate QR code

Solution:
1. Create qrcodes directory
   - mkdir qrcodes

2. Check directory permissions
   - Linux/macOS: chmod 755 qrcodes
   - Windows: Right-click → Properties → Security
```

---

## Testing

### Run Unit Tests
```bash
go test -v ./...
```

### Run Specific Test
```bash
go test -v -run TestValidateEmail
```

### Generate Coverage Report
```bash
go test -cover ./...
```

### Benchmark Tests
```bash
go test -bench=. ./...
```

---

## Development Tips

### File Editing
- **main.go** - Application entry point and routes
- **models.go** - Data structures
- **database.go** - Database setup and migrations
- **API_DOCUMENTATION.md** - API reference

### Common Tasks

#### Add New API Endpoint
1. Create handler function in `api_handlers.go`
2. Add route in `main.go`
3. Add test case in `main_test.go`
4. Update `API_DOCUMENTATION.md`

#### Add Promo Code
```go
promo, err := CreatePromoCode(
    "SAVE20",           // code
    "20% off",          // description
    "percentage",       // type
    20,                 // value
    100,                // max uses
    time.Now().AddDate(0, 1, 0), // expires
)
```

#### Test Email
```go
SendBookingConfirmation(&Booking{
    FirstName: "Test",
    LastName: "User",
    Email: "test@example.com",
    // ... other fields
})
```

---

## Production Deployment

### Pre-deployment Checklist
- [ ] All tests passing
- [ ] Update JWT_SECRET in .env
- [ ] Configure production database
- [ ] Setup SMTP credentials
- [ ] Enable HTTPS/SSL
- [ ] Configure CORS for production domain
- [ ] Setup monitoring and logging
- [ ] Create database backups

### Docker Deployment
```bash
# Build image
docker build -t booking-system:1.0 .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_USER=postgres \
  -e DB_PASSWORD=secure_password \
  booking-system:1.0
```

### Kubernetes Deployment
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

---

## Useful Commands

### Database
```bash
# Connect to database
psql -U postgres -d booking_system

# View all tables
\dt

# View schema of table
\d bookings

# View all indexes
\di

# Backup database
pg_dump booking_system > backup.sql

# Restore database
psql booking_system < backup.sql
```

### Go
```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# View test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Profile application
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
go tool pprof cpu.prof
```

### Logs
```bash
# View application logs
tail -f logs/app.log

# View error logs
tail -f logs/error.log

# View audit logs
tail -f logs/audit.log

# Search logs
grep "error" logs/*.log
```

---

## Support

If you encounter issues:

1. Check [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
2. Review [FEATURES.md](FEATURES.md)
3. Check logs in `logs/` directory
4. Search existing GitHub issues
5. Create new issue with:
   - Error message
   - Steps to reproduce
   - Operating system
   - Go version

---

## Next Steps

1. ✅ Complete installation
2. ✅ Run tests to verify setup
3. ✅ Create test bookings
4. ✅ Review API documentation
5. ✅ Configure email (optional)
6. ✅ Setup payment processing (optional)
7. ✅ Deploy to production

---

**You're all set! 🎉**

Start booking: `http://localhost:8080`
