package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/rs/cors"
)

// Request defines the JSON structure of the incoming POST payload.
type Request struct {
	Packs []int `json:"packs"`
	Order int   `json:"order"`
}

// packChange returns a map where keys are pack sizes (as strings)
// and values are the counts of packs used to cover at least 'order'
// while minimizing the total number of items shipped and then the pack count.
// It assumes that validation for positive numbers has been done.
func packChange(packs []int, order int) map[string]int {
	if len(packs) == 0 {
		return nil
	}

	// Find the smallest pack.
	minPack := packs[0]
	for _, pack := range packs {
		if pack < minPack {
			minPack = pack
		}
	}

	// We only need to consider sums up to order + (minPack - 1)
	limit := order + minPack - 1

	// dp[s] will hold the minimum number of packs needed to achieve a total of s items.
	// If s is not achievable, dp[s] will remain math.MaxInt32.
	dp := make([]int, limit+1)
	// packUsed[s] stores the last pack (pack size) used to form s items.
	packUsed := make([]int, limit+1)

	// Initialize dp: 0 items requires 0 packs; all other sums start as unreachable.
	dp[0] = 0
	packUsed[0] = -1
	for s := 1; s <= limit; s++ {
		dp[s] = math.MaxInt32
		packUsed[s] = -1
	}

	// Compute dp for sums 1 ... limit.
	sort.Ints(packs)
	for s := 1; s <= limit; s++ {
		for _, pack := range packs {
			if s >= pack && dp[s-pack] != math.MaxInt32 {
				candidate := dp[s-pack] + 1
				// We want to minimize the number of packs.
				if candidate < dp[s] {
					dp[s] = candidate
					packUsed[s] = pack
				}
			}
		}
	}

	// Find the smallest total S (>= order) that is achievable.
	bestSum := -1
	for s := order; s <= limit; s++ {
		if dp[s] != math.MaxInt32 {
			bestSum = s
			break
		}
	}
	if bestSum == -1 {
		// If no solution can be formed even with the extra items, return nil.
		return nil
	}

	// Backtrack from bestSum to recover the combination.
	result := backtrack(packUsed, bestSum)
	return result
}

// backtrack reconstructs the pack counts used to form a total sum.
func backtrack(packUsed []int, total int) map[string]int {
	res := make(map[string]int)
	for total > 0 {
		pack := packUsed[total]
		if pack == -1 {
			break // safety check; should not happen if total is achievable.
		}
		key := strconv.Itoa(pack)
		res[key]++
		total -= pack
	}
	return res
}

// packChangeHandler processes incoming POST requests.
func packChangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		// The CORS middleware will handle setting the correct headers.
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validate order is positive
	if req.Order <= 0 {
		http.Error(w, "Order must be greater than 0", http.StatusBadRequest)
		return
	}

	// Validate each pack size is a positive integer
	if len(req.Packs) == 0 {
		http.Error(w, "At least one pack size must be provided", http.StatusBadRequest)
		return
	}
	for i, pack := range req.Packs {
		if pack <= 0 {
			http.Error(w, fmt.Sprintf("Pack size at index %d must be greater than 0", i), http.StatusBadRequest)
			return
		}
	}

	result := packChange(req.Packs, req.Order)
	if result == nil {
		http.Error(w, "Cannot form a valid pack combination", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	// Create a custom mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", packChangeHandler)

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Wrap mux with CORS middleware
	handler := c.Handler(mux)

	fmt.Println("Server listening on port 5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}
