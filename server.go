package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
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
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract form data
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

	// Check if borrower is blacklisted
	if isBlacklisted(name) {
		// http.Error(w, "You are not eligible for a loan", http.StatusForbidden)
		Error = "You are not eligible for a loan"
		return
	}

	// Check if borrower has exceeded the maximum number of applications
	if isOverApplicationLimit(personalID) {
		// http.Error(w, "You have exceeded the maximum number of loan applications", http.StatusTooManyRequests)
		Error = "You have exceeded the maximum number of loan applications"
		return
	}

	// Calculate monthly payment
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
	monthlyInterestRate := 0.05 / 12.0 // 5% yearly interest rate divided by 12 for monthly rate
	power := math.Pow(1+monthlyInterestRate, float64(term))
	monthlyPayment := int64((monthlyInterestRate * float64(amount) * power) / (power - 1))

	return monthlyPayment
}

func isBlacklisted(name string) bool {
	// Read blacklist file
	blacklistData, err := os.ReadFile("blacklist.txt")

	if err != nil {
		return false
	}

	// Check if personal ID is in blacklist
	blacklist := strings.Split(string(blacklistData), "\n")
	for _, id := range blacklist {
		if id == name {
			return true
		}
	}
	return false
}

func isOverApplicationLimit(personalID string) bool {
	applicationInterval := 24 * time.Hour / 5.0
	// Count number of loan applications in the past 24 hours
	applicationCount := 0
	applicationData, err := os.ReadFile("applications.txt")
	if err == nil {
		applicationStrings := strings.Split(string(applicationData), "\n")
		for _, applicationString := range applicationStrings {
			if applicationString == "" {
				continue
			}
			applicationFields := strings.Split(applicationString, ",")
			if applicationFields[2] == personalID {
				applicationTime, err := time.Parse(time.RFC3339, applicationFields[3])
				if err == nil && time.Since(applicationTime) < applicationInterval {
					applicationCount++
				}
			}
		}
	}

	// Check if the number of applications exceeds the limit
	return applicationCount >= 5
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
