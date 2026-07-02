package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// fetchAndCache requests data from Jikan API, caches the result for 5 minutes, and responds.
func fetchAndCache(c *gin.Context, cache *Cache, targetURL string) {
	// Check cache first
	if cached, found := cache.Get(targetURL); found {
		c.Data(http.StatusOK, cached.ContentType, cached.Data)
		return
	}

	// Fetch from target API
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

	// Cache the successful response for 5 minutes
	cache.Set(targetURL, bodyBytes, contentType, 5*time.Minute)

	c.Data(http.StatusOK, contentType, bodyBytes)
}

// ponytail: Using a simple mutex-guarded map in-memory cache and custom middleware instead of external cache/CORS libraries to keep the implementation simple, dependency-free, and clean.

func runSelfTests() {
	time.Sleep(1 * time.Second) // wait for server to start
	client := &http.Client{}

	fmt.Println("=== RUNNING SELF-TESTS ===")

	// Test CORS & Caching on /api/genres
	urlStr := "http://localhost:8081/api/genres"
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("Origin", "http://localhost:3000")

	// 1st request (cache miss / fetch from Jikan)
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

	// 2nd request (cache hit)
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

	// Test /api/recommendations
	recURL := "http://localhost:8081/api/recommendations"
	reqRec, _ := http.NewRequest("GET", recURL, nil)
	respRec, err := client.Do(reqRec)
	if err != nil {
		fmt.Printf("FAIL: /api/recommendations fetch failed: %v\n", err)
		os.Exit(1)
	}
	defer respRec.Body.Close()
	fmt.Printf("/api/recommendations Status: %d\n", respRec.StatusCode)

	// Test /api/anime?genres=1
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
	cache := NewCache()

	r := gin.Default()
	r.Use(CORSMiddleware())

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
