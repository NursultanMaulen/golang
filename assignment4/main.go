package main

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type User struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

func setupLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",  // PATH
		MaxSize:    10,              // MAX FILE SIZE IN LOG
		MaxBackups: 3,               // MAX BACKUP SIZE
		MaxAge:     28,              // MAX AGE TO CONTAIN
		Compress:   true,            // COMPRESS OLD LOGS
	})
}

// VALIDATOR
var validate = validator.New()

var logger = logrus.New()

var users = make(map[string]User)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)
var (
	totalRequests  = expvar.NewInt("total_requests") 
	totalErrors    = expvar.NewInt("total_errors")
	requestLatency = expvar.NewFloat("request_latency")
)

func initPrometheus() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func RequestMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		statusCode := wrappedWriter.statusCode
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode)).Observe(duration)
	})
}

func trackRequestMetrics(start time.Time) {
	duration := time.Since(start).Seconds()
	requestLatency.Set(duration)
	totalRequests.Add(1) 
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func initLogger() {
	if os.Getenv("APP_ENV") == "production" {
		// ONLY ERRORS FOR PRODUCTION
		logger.SetLevel(logrus.ErrorLevel)
	} else {
		// SHOWING DETAILS FOR DEVS
		logger.SetLevel(logrus.DebugLevel)
	}
}

func RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// GIVING MANAGEMENT TO THE HTTP
		next.ServeHTTP(wrappedWriter, r)

		// LOGGING REQUEST STATUS
		duration := time.Since(start)
		logger.WithFields(logrus.Fields{
			"method":       r.Method,
			"path":         r.URL.Path,
			"status":       wrappedWriter.statusCode,
			"duration_ms":  duration.Milliseconds(),
			"user_agent":   r.UserAgent(),
			"remote_addr":  r.RemoteAddr,
		}).Info("HTTP request processed")
	})
}

// statusRecorder FOR LOGGING RESPONSE STATUS
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader CATCHES RESPONSE STATUS
func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content-Security-Policy (CSP)
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none';")
		// X-Frame-Options
		w.Header().Set("X-Frame-Options", "DENY")
		// X-XSS-Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Strict-Transport-Security (HSTS)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// X-Content-Type-Options
		w.Header().Set("X-Content-Type-Options", "nosniff")

		next.ServeHTTP(w, r)
	})
}

// RegisterHandler handles post requests for registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer trackRequestMetrics(start)
	var user User

	// DECODE JSON REQ TO USER
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		totalErrors.Add(1)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// STRUCTURE VALIDATION
	err = validate.Struct(user)
	if err != nil {
		// FORM READABLE ERRORS LIST
		totalErrors.Add(1)
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.ActualTag()))
		}
		http.Error(w, fmt.Sprintf("FAIL VALIDATION: %v", validationErrors), http.StatusBadRequest)
		return
	}

	// PASSWORD HASHING
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		http.Error(w, "FAIL HASHING PASSWORD", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// INIT USER TO LOCAL STORAGE
	users[user.Username] = user

	// RESPONSE FOR SUCCESSFUL VALIDATION
	logger.Info(fmt.Sprintf("USER REGISTERED: %v", err))
	fmt.Fprintf(w, "USER REGISTERED SUCCESSFULLY: %+v\n", user)
}

func csrfTokenHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"csrf_token": csrfToken,
	})
	logger.Warn(fmt.Printf("CSRF TOKEN GENERATED!"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		logger.Warn("LOGIN ATTEMPT WITH EMPTY CREDS")
		http.Error(w, "USERNAME AND PASSWORD REQUIRED", http.StatusBadRequest)
		return
	}

	logger.WithFields(logrus.Fields{
		"username": username,
	}).Info("USER LOGIN ATTEMPT")

	// SEARCH USER
	user, exists := users[username]
	if !exists {
		logger.Error(fmt.Printf("USER NOT FOUND ERROR OCC"))
		http.Error(w, "USER NOT FOUND", http.StatusUnauthorized)
		return
	}
	logger.WithFields(logrus.Fields{
		"username": username,
	}).Info("USER LOGIN SUCCESSFULLY")
	fmt.Fprintf(w, "LOGIN SUCCESSFUL")

	// CHECK PASSWORD
	err := CheckPasswordHash(password, user.Password)
	if err != nil {
		logger.Error(fmt.Sprintf("INVALID CREDENTIALS ENTERED: %v", err))
		http.Error(w, "INVALID CREDENTIALS", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user.Username, user.Role)
	if err != nil {
		logger.Error(fmt.Sprintf("TOKEN GENERATION ERROR: %v", err))
		http.Error(w, "FAIL TOKEN GENERATION", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "TOKEN: %s\n", token)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("PROTECTED ROUTE ACCESSED")
	fmt.Fprintln(w, "THIS IS PROTECTED ROUTE!!!")
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ADMIN-ONLY CONTENT!")
}

func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// LOGGING ERROR, WITHOUT SHOWING TRACKS
				logger.Error(fmt.Sprintf("Error occurred: %v", err))

				// SENDING UNITE RESPONSE TO USER
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// SomeHandler EXAMPLE FOR DEMONSTRATION
func someHandler(w http.ResponseWriter, r *http.Request) {
	// EXAMPLE OF ERROR
	err := fmt.Errorf("SIMULATED ERROR TRIGGERED")
	if err != nil {
		panic(err) // ERROR GENERATION
	}
	fmt.Fprintln(w, "REQUEST SUCCESS")
}

func main() {

	initPrometheus()

	setupLogger()

	initLogger()

	r := mux.NewRouter()

	r.Use(errorHandler)
	r.Use(RequestLoggingMiddleware)
	r.Use(RequestMetricsMiddleware)
	r.Use(loggingMiddleware)

	csrfMiddleware := csrf.Protect([]byte("sdfasddasdas"), csrf.Secure(true), csrf.MaxAge(60))

	// OPEN ROUTES
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/register", RegisterHandler)
	r.HandleFunc("/some-endpoint", someHandler).Methods("GET")

	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/debug/vars", http.DefaultServeMux)

	r.Use(SecurityHeadersMiddleware)

	// PROTECTED ROUTE
	r.Handle("/protected", AuthMiddleware(http.HandlerFunc(protectedHandler)))
	r.Handle("/admin", RBACMiddleware("admin")(http.HandlerFunc(adminHandler)))

	r.HandleFunc("/csrf-token", csrfTokenHandler).Methods("GET")

	r.Use(csrfMiddleware)

	certFile := "server.crt"
	keyFile := "server.key"

	logrus.Info("Application started")

	// SERVER START
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServeTLS(":8080", certFile, keyFile, r))
}
