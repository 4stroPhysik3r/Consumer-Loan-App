package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type loanApplication struct {
	amount     int64
	term       int64
	name       string
	personalID string
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", pathHandler)
	http.HandleFunc("/result", pathHandler)
	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		resultHandler(w, r)
	} else {
		indexHandler(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseInt(r.Form.Get("amount"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid loan amount", http.StatusBadRequest)
		return
	}
	term, err := strconv.ParseInt(r.Form.Get("term"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid loan term", http.StatusBadRequest)
		return
	}
	name := r.Form.Get("name")
	personalID := r.Form.Get("personalID")
	Error := ""

	if isBlacklisted(name) {
		http.Error(w, "You are not eligible for a loan", http.StatusForbidden)
		return
	}

	if isOverAppLimit(name, 5) {
		http.Error(w, "Too many loan applications within 24h", http.StatusTooManyRequests)
		return
	}

	monthlyPayment := calculateMonthlyPayment(amount, term)

	// Record loan application
	loanApplication := loanApplication{amount: amount, term: term, name: name, personalID: personalID}
	recordLoanApplication(loanApplication)

	// Render loan result
	tmpl := template.Must(template.ParseFiles("static/result.html"))
	data := struct {
		Amount         int64
		Term           int64
		Name           string
		PersonalID     string
		MonthlyPayment int64
		Error          string
	}{amount, term, name, personalID, monthlyPayment, Error}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func calculateMonthlyPayment(amount, term int64) int64 {
	monthlyInterestRate := 0.05 / 12.0 // 5% yearly interest rate
	power := math.Pow(1+monthlyInterestRate, float64(term))
	monthlyPayment := int64((monthlyInterestRate * float64(amount) * power) / (power - 1))

	return monthlyPayment
}

func isBlacklisted(name string) bool {
	// Read blacklist file
	content, err := os.ReadFile("blacklist.txt")
	if err != nil {
		fmt.Printf("Error reading blacklist file")
	}

	// Check if name is in blacklist
	if strings.Contains(string(content), name) {
		return true
	} else {
		return false
	}
}

func isOverAppLimit(name string, limit int) bool {
	file, err := os.Open("applications.txt")
	if err != nil {
		return false
	}
	defer file.Close()

	// Read the file line by line and store the entries for the given name
	var entries []time.Time
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ", ")
		if len(fields) >= 5 {
			n := fields[2]
			if n == name {
				timestampStr := fields[4]
				timestamp, err := time.Parse(time.RFC3339, timestampStr)
				if err != nil {
					return false
				}
				entries = append(entries, timestamp)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return false
	}

	// If there are less than `limit` entries, return false
	if len(entries) < limit {
		return false
	}

	// Sort the entries by timestamp and check if there are more than `limit` entries within the last 24 hours
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Before(entries[j])
	})
	var count int
	var currentDay time.Time
	for _, entry := range entries {
		if !entry.Truncate(24 * time.Hour).Equal(currentDay) {
			currentDay = entry.Truncate(24 * time.Hour)
			count = 0
		}
		count++
		if count > limit {
			// If there are more than `limit` entries on the same day, return true
			return true
		}
	}
	return false
}

func recordLoanApplication(loanApplication loanApplication) {
	// Write loan application details to file
	applicationString := fmt.Sprintf("%d, %d, %s, %s, %s\n", loanApplication.amount, loanApplication.term, loanApplication.name, loanApplication.personalID, time.Now().Format(time.RFC3339))
	f, err := os.OpenFile("applications.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(applicationString); err != nil {
		log.Fatal(err)
	}
}
