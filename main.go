package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type OTPEntry struct {
	Code      string
	ExpiresAt time.Time
}

var otpStore = make(map[string]OTPEntry)

var mutex sync.Mutex

func generateOTPHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiry := time.Now().Add(5 * time.Minute)

	mutex.Lock()
	otpStore[email] = OTPEntry{Code: code, ExpiresAt: expiry}
	mutex.Unlock()

	fmt.Fprintf(w, "OTP for %s is %s and expires in 5 mins\n", email, code)
}

func validateOTPHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	code := r.URL.Query().Get("code")
	if email == "" || code == "" {
		http.Error(w, "Email and code required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	entry, exists := otpStore[email]
	mutex.Unlock()

	if !exists {
		http.Error(w, "No active OTP found for this email", http.StatusNotFound)
		return
	}
	if time.Now().After(entry.ExpiresAt) {
		http.Error(w, "OTP expired", http.StatusUnauthorized)
		return
	}
	if entry.Code != code {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	// invalidate otp after use
	mutex.Lock()
	delete(otpStore, email)
	mutex.Unlock()

	fmt.Fprintln(w, "OTP verified successfully")
}

func main() {

	http.HandleFunc("/generate", generateOTPHandler)
	http.HandleFunc("/validate", validateOTPHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Simple OTP Service\n\nUse /generate?email=your_email to generate an OTP\nUse /validate?email=your_email&code=your_otp to validate it")
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
