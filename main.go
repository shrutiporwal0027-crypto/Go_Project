package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const conferenceTickets = 50

var (
	conferenceName   = "AI Technology Summit 2026"
	remainingTickets = 50
	bookings         []UserData
	bookingCounter   = 0
	mu               sync.Mutex
)

type UserData struct {
	ID              int
	FirstName       string
	LastName        string
	Email           string
	NumberOfTickets int
	BookedAt        string
	BookingRef      string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/confirmation", confirmationHandler)
	http.HandleFunc("/bookings", bookingsHandler)
	http.HandleFunc("/cancel", cancelHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/export", exportHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/schedule", scheduleHandler)
	http.HandleFunc("/contact", contactHandler)
	http.Handle("/style.css", http.FileServer(http.Dir(".")))

	fmt.Println("🚀 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	cancelled := r.URL.Query().Get("cancelled") == "true"
	success := r.URL.Query().Get("success") == "true"

	tmpl, _ := template.ParseFiles("index.html")
	data := map[string]interface{}{
		"ConferenceName":   conferenceName,
		"RemainingTickets": remainingTickets,
		"TotalTickets":     conferenceTickets,
		"BookingCount":     len(bookings),
		"Cancelled":        cancelled,
		"Success":          success,
	}
	tmpl.Execute(w, data)
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("book.html")
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	firstName := sanitizeInput(r.FormValue("firstName"))
	lastName := sanitizeInput(r.FormValue("lastName"))
	email := strings.TrimSpace(r.FormValue("email"))
	ticketsStr := r.FormValue("tickets")

	// Validation
	errors := validateBooking(firstName, lastName, email, ticketsStr)
	if len(errors) > 0 {
		tmpl, _ := template.ParseFiles("book.html")
		tmpl.Execute(w, map[string]interface{}{"Errors": errors})
		return
	}

	tickets, _ := strconv.Atoi(ticketsStr)

	mu.Lock()

	if tickets > remainingTickets {
		mu.Unlock()
		http.Error(w, fmt.Sprintf("Only %d tickets available", remainingTickets), http.StatusBadRequest)
		return
	}

	bookingCounter++
	bookingRef := fmt.Sprintf("BK%d%03d", time.Now().Unix()%1000, bookingCounter)

	booking := UserData{
		ID:              bookingCounter,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		NumberOfTickets: tickets,
		BookedAt:        time.Now().Format("2006-01-02 15:04:05"),
		BookingRef:      bookingRef,
	}

	bookings = append(bookings, booking)
	remainingTickets -= tickets
	bookingID := booking.ID
	mu.Unlock()
	http.Redirect(w, r, fmt.Sprintf("/confirmation?id=%d", bookingID), http.StatusSeeOther)
}

func cancelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, _ := strconv.Atoi(idStr)

	mu.Lock()
	defer mu.Unlock()

	for i, booking := range bookings {
		if booking.ID == id {
			remainingTickets += booking.NumberOfTickets
			bookings = append(bookings[:i], bookings[i+1:]...)

			http.Redirect(w, r, "/?cancelled=true", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Booking not found", http.StatusNotFound)
}

func confirmationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	mu.Lock()
	defer mu.Unlock()

	var booking *UserData
	for i := range bookings {
		if bookings[i].ID == id {
			booking = &bookings[i]
			break
		}
	}

	if booking == nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	const confirmationTmpl = `
	<html>
	<head><meta charset="UTF-8"><title>Booking Confirmation</title>
	<link rel="stylesheet" href="style.css"></head>
	<body>
	<nav style="background: linear-gradient(90deg, #1a237e, #283593); padding: 15px; position: sticky; top: 0; z-index: 100;">
		<div class="container" style="display: flex; gap: 30px;">
			<a href="/" style="color: white; text-decoration: none; font-size: 18px; font-weight: bold;">🏠 Home</a>
			<a href="/bookings" style="color: white; text-decoration: none; font-size: 18px;">📋 All Bookings</a>
			<a href="/book" style="color: white; text-decoration: none; font-size: 18px;">🎫 Book Now</a>
		</div>
	</nav>
	<div class="container">
	<div style="background: white; padding: 30px; border-radius: 15px; box-shadow: 0 8px 32px rgba(0,0,0,0.2); text-align: center; max-width: 600px; margin: 50px auto;">
		<h1 style="color: #28a745; margin: 0; font-size: 36px;">✅ Booking Confirmed!</h1>
		<p style="color: #666; margin: 10px 0 30px 0; font-size: 18px;">Thank you for booking with us!</p>

	<div style="background: hsla(0, 11%, 95%, 0.86); padding: 20px; border-radius: 10px; text-align: left; margin-bottom: 20px;">

	<p style="color:#1e3a8a; font-size:16px;">
		<strong style="color:#111827;">Booking Reference:</strong> {{.BookingRef}}
	</p>

	<p style="color:#7c3aed; font-size:16px;">
		<strong style="color:#111827;">Booking ID:</strong> #{{.ID}}
	</p>

	<p style="color:#059669; font-size:16px;">
		<strong style="color:#111827;">Name:</strong> {{.FirstName}} {{.LastName}}
	</p>

	<p style="color:#dc2626; font-size:16px;">
		<strong style="color:#111827;">Email:</strong> {{.Email}}
	</p>

	<p style="color:#ea580c; font-size:16px;">
		<strong style="color:#111827;">Tickets:</strong> {{.NumberOfTickets}}
	</p>

	<p style="color:#0891b2; font-size:16px;">
		<strong style="color:#111827;">Booked At:</strong> {{.BookedAt}}
	</p>

</div>

		<div style="background: #b3f0b3; padding: 15px; border-radius: 8px; margin-bottom: 20px;">
			<p style="color: #2e7d32; margin: 0;"><strong>✅ Your booking is confirmed.</strong></p>
		</div>

		<a href="/" style="background: linear-gradient(135deg, #667eea, #764ba2); color: white; padding: 12px 30px; text-decoration: none; border-radius: 8px; display: inline-block; font-weight: bold;">← Back to Home</a>
	</div>
	</div>
	</body>
	</html>
	`

	t := template.Must(template.New("confirmation").Parse(confirmationTmpl))
	if err := t.Execute(w, booking); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func bookingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle cancellation
		idStr := r.FormValue("id")
		id, _ := strconv.Atoi(idStr)

		mu.Lock()
		for i, booking := range bookings {
			if booking.ID == id {
				remainingTickets += booking.NumberOfTickets
				bookings = append(bookings[:i], bookings[i+1:]...)
				mu.Unlock()
				http.Redirect(w, r, "/bookings?cancelled=true", http.StatusSeeOther)
				return
			}
		}
		mu.Unlock()
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
	<html>
	<head><meta charset="UTF-8"><title>All Bookings</title>
	<link rel="stylesheet" href="style.css"></head>
	<body>
	<nav style="background: linear-gradient(90deg, #1a237e, #283593); padding: 15px; position: sticky; top: 0; z-index: 100;">
		<div class="container" style="display: flex; gap: 30px;">
			<a href="/" style="color: white; text-decoration: none; font-size: 18px; font-weight: bold;">🏠 Home</a>
			<a href="/bookings" style="color: white; text-decoration: none; font-size: 18px;">📋 All Bookings</a>
			<a href="/book" style="color: white; text-decoration: none; font-size: 18px;">🎫 Book Now</a>
			<a href="/search" style="color: white; text-decoration: none; font-size: 18px;">🔍 Search</a>
			<a href="/stats" style="color: white; text-decoration: none; font-size: 18px;">📊 Stats</a>
			<a href="/export" style="color: white; text-decoration: none; font-size: 18px;">📥 Export</a>
		</div>
	</nav>
	<div class="container">
	<h1>📋 All Bookings</h1>
	`)

	if len(bookings) == 0 {
		fmt.Fprintf(w, `<p>No bookings yet</p>`)
	} else {
		fmt.Fprintf(w, `<table style="width:100%%; border-collapse:collapse;">
		<tr style="background:#667eea; color:white;">
		<th style="padding:10px; border:1px solid #ddd;">Ref</th>
		<th style="padding:10px; border:1px solid #ddd;">ID</th>
		<th style="padding:10px; border:1px solid #ddd;">Name</th>
		<th style="padding:10px; border:1px solid #ddd;">Email</th>
		<th style="padding:10px; border:1px solid #ddd;">Tickets</th>
		<th style="padding:10px; border:1px solid #ddd;">Date</th>
		<th style="padding:10px; border:1px solid #ddd;">Action</th>
		</tr>`)

		for _, booking := range bookings {
			fmt.Fprintf(w, `<tr style="border:1px solid #ddd;">
			<td style="padding:10px; font-family:monospace;">%s</td>
			<td style="padding:10px;">#%d</td>
			<td style="padding:10px;">%s %s</td>
			<td style="padding:10px;">%s</td>
			<td style="padding:10px;">%d</td>
			<td style="padding:10px;">%s</td>
			<td style="padding:10px;">
			<form method="POST" action="/bookings" style="display:inline;">
			<input type="hidden" name="id" value="%d">
			<button type="submit" style="background:#d9534f; color:white; padding:5px 10px; border:none; cursor:pointer; border-radius:5px;">Cancel</button>
			</form>
			</td>
			</tr>`, booking.BookingRef, booking.ID, booking.FirstName, booking.LastName, booking.Email, booking.NumberOfTickets, booking.BookedAt, booking.ID)
		}

		fmt.Fprintf(w, `</table>`)
	}

	fmt.Fprintf(w, `
	<p style="margin-top: 20px;"><a href="/" style="color:#667eea;">← Back to Home</a></p>
	</div>
	</body>
	</html>
	`)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	tmpl, _ := template.ParseFiles("about.html")
	tmpl.Execute(w, nil)
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	tmpl, _ := template.ParseFiles("schedule.html")
	tmpl.Execute(w, nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	tmpl, _ := template.ParseFiles("contact.html")
	tmpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
		<html>
		<head><meta charset="UTF-8"><title>Search Bookings</title>
		<link rel="stylesheet" href="style.css"></head>
		<body>
		<nav style="background: linear-gradient(90deg, #1a237e, #283593); padding: 15px; position: sticky; top: 0; z-index: 100;">
			<div class="container" style="display: flex; gap: 30px;">
				<a href="/" style="color: white; text-decoration: none; font-size: 18px; font-weight: bold;">🏠 Home</a>
				<a href="/bookings" style="color: white; text-decoration: none; font-size: 18px;">📋 All Bookings</a>
				<a href="/search" style="color: white; text-decoration: none; font-size: 18px;">🔍 Search</a>
			</div>
		</nav>
		<div class="container">
		<h1>🔍 Search Bookings</h1>
		<form method="POST" style="background: white; padding: 30px; border-radius: 15px; box-shadow: 0 8px 32px rgba(0,0,0,0.2);">
			<div style="margin-bottom: 20px;">
				<label style="display: block; font-weight: bold; margin-bottom: 8px;">Search by Name or Email:</label>
				<input type="text" name="query" required placeholder="Enter name or email" style="width: 100%%; padding: 12px; border: 2px solid #e0e0e0; border-radius: 8px; font-size: 16px;">
			</div>
			<button type="submit" style="background: linear-gradient(135deg, #667eea, #764ba2); color: white; padding: 12px 30px; border: none; border-radius: 8px; cursor: pointer; font-weight: bold;">🔍 Search</button>
		</form>
		<p style="margin-top: 20px;"><a href="/" style="color:#667eea;">← Back to Home</a></p>
		</div>
		</body>
		</html>
		`)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	query := strings.ToLower(strings.TrimSpace(r.FormValue("query")))

	mu.Lock()
	var results []UserData
	for _, booking := range bookings {
		if strings.Contains(strings.ToLower(booking.FirstName+" "+booking.LastName), query) ||
			strings.Contains(strings.ToLower(booking.Email), query) ||
			strings.Contains(strings.ToLower(booking.BookingRef), query) {
			results = append(results, booking)
		}
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
	<html>
	<head><meta charset="UTF-8"><title>Search Results</title>
	<link rel="stylesheet" href="style.css"></head>
	<body>
	<nav style="background: linear-gradient(90deg, #1a237e, #283593); padding: 15px; position: sticky; top: 0; z-index: 100;">
		<div class="container" style="display: flex; gap: 30px;">
			<a href="/" style="color: white; text-decoration: none; font-size: 18px; font-weight: bold;">🏠 Home</a>
			<a href="/bookings" style="color: white; text-decoration: none; font-size: 18px;">📋 All Bookings</a>
			<a href="/search" style="color: white; text-decoration: none; font-size: 18px;">🔍 Search</a>
		</div>
	</nav>
	<div class="container">
	<h1>🔍 Search Results for "%s"</h1>
	`, query)

	if len(results) == 0 {
		fmt.Fprintf(w, `<p>No bookings found matching your search.</p>`)
	} else {
		fmt.Fprintf(w, `<p>Found %d booking(s):</p><table style="width:100%%; border-collapse:collapse;">
		<tr style="background:#667eea; color:white;">
		<th style="padding:10px; border:1px solid #ddd;">Ref</th>
		<th style="padding:10px; border:1px solid #ddd;">Name</th>
		<th style="padding:10px; border:1px solid #ddd;">Email</th>
		<th style="padding:10px; border:1px solid #ddd;">Tickets</th>
		<th style="padding:10px; border:1px solid #ddd;">Date</th>
		</tr>`, len(results))

		for _, booking := range results {
			fmt.Fprintf(w, `<tr style="border:1px solid #ddd;">
			<td style="padding:10px; font-family:monospace;">%s</td>
			<td style="padding:10px;">%s %s</td>
			<td style="padding:10px;">%s</td>
			<td style="padding:10px;">%d</td>
			<td style="padding:10px;">%s</td>
			</tr>`, booking.BookingRef, booking.FirstName, booking.LastName, booking.Email, booking.NumberOfTickets, booking.BookedAt)
		}

		fmt.Fprintf(w, `</table>`)
	}

	fmt.Fprintf(w, `
	<p style="margin-top: 20px;"><a href="/search" style="color:#667eea;">← Search Again</a> | <a href="/" style="color:#667eea;">Home</a></p>
	</div>
	</body>
	</html>
	`)
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=bookings.csv")

	fmt.Fprintf(w, "Booking Reference,Booking ID,First Name,Last Name,Email,Tickets,Booked At\n")

	mu.Lock()
	defer mu.Unlock()

	for _, booking := range bookings {
		fmt.Fprintf(w, "%s,%d,%s,%s,%s,%d,%s\n",
			booking.BookingRef,
			booking.ID,
			booking.FirstName,
			booking.LastName,
			booking.Email,
			booking.NumberOfTickets,
			booking.BookedAt)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	totalBookings := len(bookings)
	totalTickets := 0
	emailDomains := make(map[string]int)

	for _, booking := range bookings {
		totalTickets += booking.NumberOfTickets
		parts := strings.Split(booking.Email, "@")
		if len(parts) == 2 {
			emailDomains[parts[1]]++
		}
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
	<html>
	<head><meta charset="UTF-8"><title>Booking Statistics</title>
	<link rel="stylesheet" href="style.css"></head>
	<body>
	<nav style="background: linear-gradient(90deg, #1a237e, #283593); padding: 15px; position: sticky; top: 0; z-index: 100;">
		<div class="container" style="display: flex; gap: 30px;">
			<a href="/" style="color: white; text-decoration: none; font-size: 18px; font-weight: bold;">🏠 Home</a>
			<a href="/bookings" style="color: white; text-decoration: none; font-size: 18px;">📋 All Bookings</a>
			<a href="/stats" style="color: white; text-decoration: none; font-size: 18px;">📊 Stats</a>
		</div>
	</nav>
	<div class="container">
	<h1>📊 Booking Statistics</h1>

	<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin-bottom: 40px;">
		<div style="background: white; padding: 30px; border-radius: 15px; box-shadow: 0 8px 32px rgba(0,0,0,0.2); text-align: center;">
			<h2 style="color: #667eea; margin: 0; font-size: 48px;">%d</h2>
			<p style="color: #666; margin: 10px 0 0 0; font-size: 18px;">Total Bookings</p>
		</div>

		<div style="background: white; padding: 30px; border-radius: 15px; box-shadow: 0 8px 32px rgba(0,0,0,0.2); text-align: center;">
			<h2 style="color: #764ba2; margin: 0; font-size: 48px;">%d</h2>
			<p style="color: #666; margin: 10px 0 0 0; font-size: 18px;">Tickets Sold</p>
		</div>

		<div style="background: white; padding: 30px; border-radius: 15px; box-shadow: 0 8px 32px rgba(0,0,0,0.2); text-align: center;">
			<h2 style="color: #ff6b6b; margin: 0; font-size: 48px;">%d</h2>
			<p style="color: #666; margin: 10px 0 0 0; font-size: 18px;">Tickets Available</p>
		</div>

		<div style="background: white; padding: 30px; border-radius: 15px; box-shadow: 0 8px 32px rgba(0,0,0,0.2); text-align: center;">
			<h2 style="color: #4caf50; margin: 0; font-size: 48px;">%.1f%%</h2>
			<p style="color: #666; margin: 10px 0 0 0; font-size: 18px;">Capacity Used</p>
		</div>
	</div>

	<div style="background: white; padding: 30px; border-radius: 15px; box-shadow: 0 8px 32px rgba(0,0,0,0.2); margin-bottom: 40px;">
		<h2>Email Domain Distribution</h2>
	`, totalBookings, totalTickets, remainingTickets, float64(totalTickets)/float64(conferenceTickets)*100)

	if len(emailDomains) > 0 {
		fmt.Fprintf(w, `<table style="width:100%%; border-collapse:collapse;">
		<tr style="background:#f5f5f5;">
		<th style="padding:10px; border:1px solid #ddd; text-align:left;">Domain</th>
		<th style="padding:10px; border:1px solid #ddd; text-align:left;">Count</th>
		</tr>`)

		for domain, count := range emailDomains {
			fmt.Fprintf(w, `<tr style="border:1px solid #ddd;">
			<td style="padding:10px;">%s</td>
			<td style="padding:10px;">%d</td>
			</tr>`, domain, count)
		}

		fmt.Fprintf(w, `</table>`)
	} else {
		fmt.Fprintf(w, `<p>No email data available</p>`)
	}

	fmt.Fprintf(w, `
	</div>

	<p style="margin-top: 20px;"><a href="/" style="color:#667eea;">← Back to Home</a></p>
	</div>
	</body>
	</html>
	`)
}

func validateBooking(firstName, lastName, email, ticketsStr string) []string {
	var errors []string

	if firstName == "" || len(firstName) < 2 {
		errors = append(errors, "First name must be at least 2 characters")
	}

	if lastName == "" || len(lastName) < 2 {
		errors = append(errors, "Last name must be at least 2 characters")
	}

	if !isValidEmail(email) {
		errors = append(errors, "Please enter a valid email address")
	}

	tickets, err := strconv.Atoi(ticketsStr)
	if err != nil || tickets <= 0 || tickets > 20 {
		errors = append(errors, "Tickets must be between 1 and 20")
	}

	return errors
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func sanitizeInput(input string) string {
	// Basic sanitization - remove potentially dangerous characters
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&apos;")
	return strings.TrimSpace(input)
}
