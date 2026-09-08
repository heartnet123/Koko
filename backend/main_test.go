package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestCacheOperations(t *testing.T) {
	c := NewCache()
	key := "test-key"
	data := []byte(`{"test": true}`)
	contentType := "application/json"

	c.Set(key, data, contentType, 100*time.Millisecond)

	item, found := c.Get(key)
	if !found {
		t.Fatalf("expected key to be found in cache")
	}
	if string(item.Data) != string(data) {
		t.Errorf("expected %s, got %s", string(data), string(item.Data))
	}
	if item.ContentType != contentType {
		t.Errorf("expected %s, got %s", contentType, item.ContentType)
	}

	time.Sleep(150 * time.Millisecond)
	_, foundAfterExpiry := c.Get(key)
	if foundAfterExpiry {
		t.Errorf("expected key to be expired")
	}
}

func TestEpisodesEndpointAndCache(t *testing.T) {
	cache := NewCache()
	router := SetupRouter(cache)

	// Pre-populate cache to simulate cached episode list
	targetURL := "https://api.jikan.moe/v4/anime/20/episodes?page=1"
	mockPayload := map[string]any{
		"pagination": map[string]any{
			"last_visible_page": 3,
			"has_next_page":     true,
		},
		"data": []map[string]any{
			{
				"mal_id": 1,
				"title":  "Enter: Naruto Uzumaki!",
				"aired":  "2002-10-03T00:00:00+00:00",
			},
			{
				"mal_id": 2,
				"title":  "My Name is Konohamaru!",
				"aired":  "2002-10-10T00:00:00+00:00",
			},
		},
	}
	mockBytes, _ := json.Marshal(mockPayload)
	cache.Set(targetURL, mockBytes, "application/json", 5*time.Minute)

	// Request cached endpoint
	req, _ := http.NewRequest("GET", "/api/anime/20/episodes?page=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var respBody map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	dataList, ok := respBody["data"].([]any)
	if !ok || len(dataList) != 2 {
		t.Fatalf("expected 2 episodes in response, got %v", respBody["data"])
	}

	ep1 := dataList[0].(map[string]any)
	if ep1["title"] != "Enter: Naruto Uzumaki!" {
		t.Errorf("unexpected episode 1 title: %v", ep1["title"])
	}
}

func TestFetchAndCachePopulatesCache(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": [{"mal_id": 1, "title": "Episode 1"}]}`))
	}))
	defer mockServer.Close()

	cache := NewCache()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	fetchAndCache(c, cache, mockServer.URL)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	cachedItem, found := cache.Get(mockServer.URL)
	if !found {
		t.Fatalf("expected response to be saved into cache on cache miss")
	}
	if string(cachedItem.Data) != `{"data": [{"mal_id": 1, "title": "Episode 1"}]}` {
		t.Errorf("unexpected cached data: %s", string(cachedItem.Data))
	}
}
