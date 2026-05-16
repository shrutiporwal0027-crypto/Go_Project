# 📋 Developer Quick Reference

## Quick Commands

```bash
# Setup
make install-deps      # Install dependencies
make db-setup          # Create PostgreSQL database
make run               # Run application

# Testing
make test              # Run all tests
make test-cov          # Coverage report
make benchmark          # Run benchmarks

# Cleaning
make clean             # Remove build artifacts
make db-drop           # Drop database

# Building
make build             # Build binary
make docker-build      # Build Docker image

# Utilities
make lint              # Run linter
make format            # Format code
make logs              # View logs
```

---

## File Quick Reference

| File | Purpose |
|------|---------|
| main.go | Entry point & route handlers |
| models.go | Data structures |
| database.go | DB setup & migrations |
| auth.go | JWT & authentication |
| validators.go | Input validation |
| api_handlers.go | REST endpoints |
| middleware.go | HTTP middleware |
| email.go | Email sending |
| payment.go | Payment processing |
| analytics.go | Analytics & reports |

---

## Common Code Snippets

### Create Booking
```go
booking := &Booking{
    FirstName: "John",
    LastName: "Doe",
    Email: "john@example.com",
    NumberOfTickets: 2,
    TotalPrice: 199.98,
    Status: "confirmed",
    TicketID: "TKT-123456",
}
```

### Validate Input
```go
errors := ValidateBookingRequest(firstName, lastName, email, tickets)
if len(errors) > 0 {
    // Handle validation errors
}
```

### Send Email
```go
err := SendBookingConfirmation(booking)
if err != nil {
    GetLogger().Error("Email failed: %v", err)
}
```

### Generate QR Code
```go
qrPath, err := GenerateQRCode(ticketID)
if err != nil {
    GetLogger().Error("QR failed: %v", err)
}
```

### Apply Promo Code
```go
finalPrice, discount, err := ValidateAndApplyPromoCode("SAVE20", totalPrice)
if err != nil {
    // Invalid code
}
```

### Get Analytics
```go
report, err := GenerateReport(startDate, endDate)
analytics := report["summary"].(map[string]interface{})
```

---

## API Endpoint Quick Reference

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -d "email=user@example.com&password=Pass123!@&first_name=John&last_name=Doe"
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -d "email=user@example.com&password=Pass123!@"
```

### Create Booking
```bash
curl -X POST http://localhost:8080/api/v1/book \
  -d "firstName=John&lastName=Doe&email=john@example.com&tickets=2"
```

### Get Bookings
```bash
curl http://localhost:8080/api/v1/bookings
```

### Get Seats
```bash
curl http://localhost:8080/api/v1/seats
```

### Get Waitlist
```bash
curl http://localhost:8080/api/v1/waitlist
```

### Get Analytics
```bash
curl "http://localhost:8080/api/v1/analytics?start=2026-06-01&end=2026-06-30"
```

---

## Debugging

### View Logs
```bash
# All logs
tail -f logs/app.log

# Errors only
tail -f logs/error.log

# Audit trail
tail -f logs/audit.log
```

### Database Query
```bash
psql -U postgres -d booking_system

# View bookings
SELECT * FROM bookings;

# View promo codes
SELECT * FROM promo_codes WHERE active = true;

# View waitlist
SELECT * FROM waitlist ORDER BY position;
```

### Test API Response
```bash
curl -i http://localhost:8080/api/v1/bookings | jq .
```

---

## Configuration Settings

### Database
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=booking_system
```

### Email
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=app-password
```

### Security
```env
JWT_SECRET=your-super-secret-key
JWT_EXPIRY_HOURS=24
```

### Business Logic
```env
BASE_TICKET_PRICE=99.99
CONFERENCE_VENUE=Convention Center
CONFERENCE_DATE=2026-06-15
```

---

## Error Messages & Solutions

| Error | Solution |
|-------|----------|
| Database connection refused | Start PostgreSQL, check credentials |
| Port 8080 in use | Change port or kill existing process |
| Email not sending | Check SMTP credentials in .env |
| QR code failed | Create qrcodes/ directory |
| Token invalid | Request new login token |
| Promo code expired | Check expires_at field in database |
| Tickets sold out | Check remaining tickets or add to waitlist |

---

## Performance Tips

1. **Database Queries**: Use WHERE clauses efficiently
2. **Caching**: Implement Redis for frequently accessed data
3. **Indexing**: Already applied to bookings(email), bookings(status)
4. **Connection Pool**: Configure max connections in database
5. **Load Testing**: Use `ab` tool to test capacity

---

## Security Checklist

- [ ] JWT_SECRET is strong
- [ ] Database password is secure
- [ ] Email credentials are correct
- [ ] CORS is restricted to trusted domains
- [ ] HTTPS is enabled in production
- [ ] Sensitive data is logged
- [ ] SQL queries use parameterization
- [ ] Input is always validated

---

## Testing Checklist

- [ ] Run `go test -v ./...`
- [ ] Check code coverage with `go test -cover`
- [ ] Run benchmarks with `go test -bench`
- [ ] Load test with `ab -n 1000 -c 100 http://localhost:8080/...`
- [ ] Manual API testing with `curl`

---

## Deployment Checklist

- [ ] All tests passing
- [ ] Code reviewed
- [ ] Database backed up
- [ ] Environment variables set
- [ ] SSL certificate installed
- [ ] Logs configured
- [ ] Monitoring enabled
- [ ] Backups scheduled
- [ ] Error tracking setup
- [ ] Performance monitoring ready

---

## Useful Go Commands

```bash
# Format code
go fmt ./...

# Run specific test
go test -run TestName -v

# Profile CPU
go test -cpuprofile=cpu.prof ./...

# Profile memory
go test -memprofile=mem.prof ./...

# View profile
go tool pprof cpu.prof

# Get dependencies tree
go mod graph

# Update dependencies
go get -u
```

---

## PostgreSQL Commands

```bash
# Connect to database
psql -U postgres -d booking_system

# List databases
\l

# List tables
\dt

# Describe table
\d bookings

# Execute query
SELECT * FROM bookings LIMIT 10;

# Export data
pg_dump booking_system > backup.sql

# Import data
psql booking_system < backup.sql
```

---

## Docker Commands

```bash
# Build image
docker build -t booking-system:latest .

# Run container
docker run -p 8080:8080 booking-system:latest

# View logs
docker logs <container_id>

# Stop container
docker stop <container_id>

# Remove image
docker rmi booking-system:latest
```

---

## Kubernetes Commands

```bash
# Deploy
kubectl apply -f k8s/deployment.yaml

# Check status
kubectl get pods

# View logs
kubectl logs <pod_name>

# Scale up
kubectl scale deployment booking-system --replicas=3

# Delete deployment
kubectl delete deployment booking-system
```

---

## Important Functions

| Function | Usage |
|----------|-------|
| ValidateEmail() | Check email format |
| ValidatePassword() | Check password strength |
| GenerateToken() | Create JWT token |
| ValidateToken() | Verify JWT token |
| RegisterUser() | Create user account |
| LoginUser() | Authenticate user |
| GenerateQRCode() | Create ticket QR |
| GeneratePDFTicket() | Create PDF ticket |
| ValidateAndApplyPromoCode() | Apply discount |
| SendBookingConfirmation() | Email confirmation |
| AddToWaitlist() | Add to queue |
| CalculateAnalytics() | Generate report |

---

## Frontend HTML Structure

```html
<!-- Main booking form -->
<form action="/book" method="POST">
  <input name="firstName" required>
  <input name="lastName" required>
  <input name="email" required>
  <input name="tickets" type="number" required>
  <button type="submit">Book Now</button>
</form>

<!-- Ticket preview -->
<div id="ticketSection">
  <div class="ticket">
    <h3>🎫 ADMISSION TICKET</h3>
    <p id="ticketName">-</p>
    <p id="ticketEmail">-</p>
    <p id="ticketId">-</p>
  </div>
</div>
```

---

## Troubleshooting Workflow

1. **Check logs**: `tail -f logs/error.log`
2. **Verify database**: `psql booking_system -c "SELECT 1"`
3. **Test API**: `curl http://localhost:8080/api/v1/bookings`
4. **Check config**: Verify .env values
5. **Restart app**: Stop and restart Go application
6. **Review code**: Check recent changes
7. **Google error**: Search for exact error message

---

## Performance Metrics to Monitor

- Request latency (< 100ms)
- Database query time (< 50ms)
- Email delivery time (< 5s)
- Ticket generation time (< 1s)
- PDF generation time (< 2s)
- Error rate (< 0.1%)
- CPU usage (< 50%)
- Memory usage (< 200MB)

---

**Good luck with your booking system! 🎉**

For full documentation, see README.md
