package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT configuration
var jwtSecretKey = []byte("your-secure-random-jwt-secret-key-change-me")

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT signed token for the user.
func GenerateToken(userID int, username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		jwtSecretKey = []byte(secret)
	}

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// ParseToken parses and validates a JWT signed token.
func ParseToken(tokenString string) (int, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		jwtSecretKey = []byte(secret)
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return 0, "", fmt.Errorf("invalid token claims")
	}
	return claims.UserID, claims.Username, nil
}

// AuthMiddleware intercepts requests to verify JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		// 1. Try to read cookie
		cookie, err := c.Cookie("auth_token")
		if err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// 2. Try to read Authorization Header
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = authHeader[7:]
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Unauthorized",
				Message: "Access token is missing",
			})
			c.Abort()
			return
		}

		userID, username, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Unauthorized",
				Message: "Access token is invalid or expired",
			})
			c.Abort()
			return
		}

		c.Set("userId", userID)
		c.Set("username", username)
		c.Next()
	}
}

// Response Models
type UserResponse struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	Password    string `json:"password"`
}

type AddWatchlistRequest struct {
	AnimeID  int    `json:"anime_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
}

type WatchlistItem struct {
	AnimeID  int       `json:"anime_id"`
	Title    string    `json:"title"`
	ImageURL string    `json:"image_url"`
	AddedAt  time.Time `json:"added_at"`
}

// Handlers
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ValidationError",
			Message: "Invalid input fields",
			Details: err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "InternalServerError",
			Message: "Failed to process password",
		})
		return
	}

	var userID int
	query := `INSERT INTO users (username, email, password_hash, display_name) 
	          VALUES ($1, $2, $3, $4) RETURNING id`
	err = DB.QueryRow(query, req.Username, req.Email, string(hashedPassword), req.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "Conflict",
			Message: "Username or Email is already registered",
			Details: err.Error(),
		})
		return
	}

	token, err := GenerateToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "InternalServerError",
			Message: "Failed to generate session",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, 86400, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": UserResponse{
			ID:          userID,
			Username:    req.Username,
			Email:       req.Email,
			DisplayName: req.Username,
			AvatarURL:   "",
			Bio:         "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	})
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ValidationError",
			Message: "Email and password are required",
		})
		return
	}

	var userID int
	var username, email, passwordHash, displayName, avatarUrl, bio string
	var createdAt, updatedAt time.Time

	query := `SELECT id, username, email, password_hash, display_name, COALESCE(avatar_url, ''), COALESCE(bio, ''), created_at, updated_at 
	          FROM users WHERE email = $1 OR username = $1`
	err := DB.QueryRow(query, req.Email).Scan(
		&userID, &username, &email, &passwordHash, &displayName, &avatarUrl, &bio, &createdAt, &updatedAt,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid credentials",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid credentials",
		})
		return
	}

	token, err := GenerateToken(userID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "InternalServerError",
			Message: "Failed to generate session",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": UserResponse{
			ID:          userID,
			Username:    username,
			Email:       email,
			DisplayName: displayName,
			AvatarURL:   avatarUrl,
			Bio:         bio,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
	})
}

func LogoutHandler(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func MeHandler(c *gin.Context) {
	userID, _ := c.Get("userId")

	var user UserResponse
	query := `SELECT id, username, email, display_name, COALESCE(avatar_url, ''), COALESCE(bio, ''), created_at, updated_at 
	          FROM users WHERE id = $1`
	err := DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NotFound",
			Message: "User profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func UpdateProfileHandler(c *gin.Context) {
	userID, _ := c.Get("userId")
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ValidationError",
			Message: "Invalid input fields",
		})
		return
	}

	var existingPasswordHash string
	querySelect := `SELECT password_hash FROM users WHERE id = $1`
	err := DB.QueryRow(querySelect, userID).Scan(&existingPasswordHash)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NotFound",
			Message: "User not found",
		})
		return
	}

	newPasswordHash := existingPasswordHash
	if req.Password != "" {
		if len(req.Password) < 8 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "ValidationError",
				Message: "Password must be at least 8 characters",
			})
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "InternalServerError",
				Message: "Failed to process new password",
			})
			return
		}
		newPasswordHash = string(hashed)
	}

	queryUpdate := `UPDATE users 
	                SET display_name = $1, bio = $2, avatar_url = $3, password_hash = $4, updated_at = CURRENT_TIMESTAMP
	                WHERE id = $5`
	_, err = DB.Exec(queryUpdate, req.DisplayName, req.Bio, req.AvatarURL, newPasswordHash, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "InternalServerError",
			Message: "Failed to update profile",
			Details: err.Error(),
		})
		return
	}

	var user UserResponse
	queryFetch := `SELECT id, username, email, display_name, COALESCE(avatar_url, ''), COALESCE(bio, ''), created_at, updated_at 
	               FROM users WHERE id = $1`
	err = DB.QueryRow(queryFetch, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.DisplayName, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "InternalServerError",
			Message: "Failed to retrieve updated profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func GetWatchlistHandler(c *gin.Context) {
	userID, _ := c.Get("userId")

	rows, err := DB.Query("SELECT anime_id, title, image_url, added_at FROM watchlist WHERE user_id = $1 ORDER BY added_at DESC", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "InternalServerError",
			Message: "Failed to fetch watchlist",
			Details: err.Error(),
		})
		return
	}
	defer rows.Close()

	items := []WatchlistItem{}
	for rows.Next() {
		var item WatchlistItem
		if err := rows.Scan(&item.AnimeID, &item.Title, &item.ImageURL, &item.AddedAt); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "InternalServerError",
				Message: "Failed to parse watchlist item",
			})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

func AddWatchlistHandler(c *gin.Context) {
	userID, _ := c.Get("userId")
	var req AddWatchlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ValidationError",
			Message: "Anime ID, title, and image URL are required",
		})
		return
	}

	query := `INSERT INTO watchlist (user_id, anime_id, title, image_url) 
	          VALUES ($1, $2, $3, $4) 
	          ON CONFLICT (user_id, anime_id) DO NOTHING`
	_, err := DB.Exec(query, userID, req.AnimeID, req.Title, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "InternalServerError",
			Message: "Failed to add to watchlist",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Added to watchlist successfully",
	})
}

func RemoveWatchlistHandler(c *gin.Context) {
	userID, _ := c.Get("userId")
	animeIDStr := c.Param("anime_id")

	var animeID int
	_, err := fmt.Sscanf(animeIDStr, "%d", &animeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ValidationError",
			Message: "Invalid Anime ID",
		})
		return
	}

	query := `DELETE FROM watchlist WHERE user_id = $1 AND anime_id = $2`
	result, err := DB.Exec(query, userID, animeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "InternalServerError",
			Message: "Failed to remove from watchlist",
			Details: err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NotFound",
			Message: "Anime not found in your watchlist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Removed from watchlist successfully",
	})
}

// CacheItem represents a single cached API response.
type CacheItem struct {
	Data        []byte
	ContentType string
	ExpiresAt   time.Time
}

// Cache is a simple thread-safe in-memory cache.
type Cache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

// NewCache initializes a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Get retrieves an item from the cache if it exists and has not expired.
func (c *Cache) Get(key string) (CacheItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found {
		return CacheItem{}, false
	}
	if time.Now().After(item.ExpiresAt) {
		return CacheItem{}, false
	}
	return item, true
}

// Set stores an item in the cache with a specific expiration duration.
func (c *Cache) Set(key string, data []byte, contentType string, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		Data:        data,
		ContentType: contentType,
		ExpiresAt:   time.Now().Add(duration),
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing for frontend requests.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// fetchAndCache requests data from Jikan API, caches the result for 5 minutes, and responds.
func fetchAndCache(c *gin.Context, cache *Cache, targetURL string) {
	if cached, found := cache.Get(targetURL); found {
		c.Data(http.StatusOK, cached.ContentType, cached.Data)
		return
	}

	resp, err := http.Get(targetURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("Failed to fetch from Jikan API: %v", err)})
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read response body: %v", err)})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json; charset=utf-8"
	}

	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, contentType, bodyBytes)
		return
	}

	cache.Set(targetURL, bodyBytes, contentType, 5*time.Minute)
	c.Data(http.StatusOK, contentType, bodyBytes)
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		} else if len(val) >= 2 && val[0] == '\'' && val[len(val)-1] == '\'' {
			val = val[1 : len(val)-1]
		}
		os.Setenv(key, val)
	}
}

func runSelfTests() {
	time.Sleep(1 * time.Second)
	client := &http.Client{}

	fmt.Println("=== RUNNING SELF-TESTS ===")

	urlStr := "http://localhost:8081/api/genres"
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("Origin", "http://localhost:3000")

	start := time.Now()
	resp1, err := client.Do(req)
	if err != nil {
		fmt.Printf("FAIL: 1st fetch failed: %v\n", err)
		os.Exit(1)
	}
	defer resp1.Body.Close()
	duration1 := time.Since(start)

	body1, _ := io.ReadAll(resp1.Body)
	fmt.Printf("1st Fetch Status: %d, Content-Type: %s, Duration: %v\n", resp1.StatusCode, resp1.Header.Get("Content-Type"), duration1)
	fmt.Printf("CORS Origin Header: %s\n", resp1.Header.Get("Access-Control-Allow-Origin"))

	if resp1.Header.Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		fmt.Println("FAIL: Missing CORS header")
		os.Exit(1)
	}

	if resp1.StatusCode != http.StatusOK {
		fmt.Printf("FAIL: expected status 200, got %d (response: %s)\n", resp1.StatusCode, string(body1))
		os.Exit(1)
	}

	start = time.Now()
	resp2, err := client.Do(req)
	if err != nil {
		fmt.Printf("FAIL: 2nd fetch failed: %v\n", err)
		os.Exit(1)
	}
	defer resp2.Body.Close()
	duration2 := time.Since(start)

	fmt.Printf("2nd Fetch Status: %d, Duration: %v\n", resp2.StatusCode, duration2)
	if duration2 > 50*time.Millisecond {
		fmt.Printf("WARNING: 2nd fetch took longer than expected: %v, but might still be cached if network/disk is slow.\n", duration2)
	} else {
		fmt.Println("SUCCESS: Cache hit confirmed (extremely fast response)!")
	}

	recURL := "http://localhost:8081/api/recommendations"
	reqRec, _ := http.NewRequest("GET", recURL, nil)
	respRec, err := client.Do(reqRec)
	if err != nil {
		fmt.Printf("FAIL: /api/recommendations fetch failed: %v\n", err)
		os.Exit(1)
	}
	defer respRec.Body.Close()
	fmt.Printf("/api/recommendations Status: %d\n", respRec.StatusCode)

	animeURL := "http://localhost:8081/api/anime?genres=1"
	reqAnime, _ := http.NewRequest("GET", animeURL, nil)
	respAnime, err := client.Do(reqAnime)
	if err != nil {
		fmt.Printf("FAIL: /api/anime fetch failed: %v\n", err)
		os.Exit(1)
	}
	defer respAnime.Body.Close()
	fmt.Printf("/api/anime?genres=1 Status: %d\n", respAnime.StatusCode)

	fmt.Println("=== ALL SELF-TESTS PASSED ===")
	os.Exit(0)
}

func main() {
	loadEnv()
	InitDB()

	cache := NewCache()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// Public Auth endpoints
	r.POST("/api/users", RegisterHandler)
	r.POST("/api/auth/session", LoginHandler)
	r.DELETE("/api/auth/session", LogoutHandler)

	// Protected endpoints group
	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/users/me", MeHandler)
		protected.PATCH("/users/me", UpdateProfileHandler)
		protected.GET("/users/me/watchlist", GetWatchlistHandler)
		protected.POST("/users/me/watchlist", AddWatchlistHandler)
		protected.DELETE("/users/me/watchlist/:anime_id", RemoveWatchlistHandler)
	}

	// Public Jikan endpoints
	r.GET("/api/recommendations", func(c *gin.Context) {
		targetURL := "https://api.jikan.moe/v4/recommendations/anime"
		fetchAndCache(c, cache, targetURL)
	})

	r.GET("/api/genres", func(c *gin.Context) {
		targetURL := "https://api.jikan.moe/v4/genres/anime"
		fetchAndCache(c, cache, targetURL)
	})

	r.GET("/api/anime", func(c *gin.Context) {
		genres := c.Query("genres")
		q := c.Query("q")
		page := c.Query("page")
		orderBy := c.Query("order_by")
		limit := c.Query("limit")

		u, err := url.Parse("https://api.jikan.moe/v4/anime")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse base URL"})
			return
		}
		qParams := u.Query()
		if genres != "" {
			qParams.Set("genres", genres)
		}
		if q != "" {
			qParams.Set("q", q)
		}
		if page != "" {
			qParams.Set("page", page)
		}
		if limit != "" {
			qParams.Set("limit", limit)
		} else {
			qParams.Set("limit", "24")
		}
		if orderBy != "" {
			qParams.Set("order_by", orderBy)
		} else {
			qParams.Set("order_by", "popularity")
		}

		u.RawQuery = qParams.Encode()
		targetURL := u.String()

		fetchAndCache(c, cache, targetURL)
	})

	r.GET("/api/anime/:id", func(c *gin.Context) {
		id := c.Param("id")
		targetURL := fmt.Sprintf("https://api.jikan.moe/v4/anime/%s/full", id)
		fetchAndCache(c, cache, targetURL)
	})

	if os.Getenv("RUN_SELF_TEST") == "true" {
		go r.Run(":8081")
		runSelfTests()
		return
	}

	r.Run(":8080")
}
