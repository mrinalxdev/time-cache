# TimeCache Usage Guide

TimeCache is a powerful, flexible caching solution for Go applications. This guide will walk you through how to use TimeCache effectively in your projects.

## Installation

To get started with TimeCache, install it using `go get`:

```bash
go get github.com/yourusername/timecache
```

## Basic Usage

Here's a quick example of how to use TimeCache:

```go
import (
    "github.com/yourusername/timecache"
    "time"
)

// Create a new cache with a 5-minute TTL
cache := timecache.NewCache(5*time.Minute, nil)

// Set a value
cache.Set("key", "value")

// Get a value
value, found := cache.Get("key")

// Remove a value
cache.Remove("key")
```

## Advanced Usage: E-commerce Product Catalog Example

Let's walk through a more complex example using an e-commerce product catalog scenario.

### Step 1: Initialize TimeCache

```go
productCache := timecache.NewCache(15*time.Minute, func(key string, value interface{}) {
    fmt.Printf("Product %s removed from cache\n", key)
})
```

This creates a cache with a 15-minute TTL and a callback function that logs when items are removed.

### Step 2: Define Your Data Structure

```go
type Product struct {
    ID    string
    Name  string
    Price float64
}
```

### Step 3: Implement Caching Logic

```go
func getProduct(id string) (Product, bool) {
    // Try to get the product from cache
    if cachedProduct, found := productCache.Get(id); found {
        return cachedProduct.(Product), true
    }

    // If not in cache, fetch from database
    if product, exists := fetchProductFromDB(id); exists {
        // Add to cache for future requests
        productCache.Set(id, product)
        return product, true
    }

    return Product{}, false
}
```

### Step 4: Use the Cache in Your Application

```go
productIDs := []string{"P001", "P002", "P003", "P001", "P004"}
for _, id := range productIDs {
    if product, found := getProduct(id); found {
        fmt.Printf("Retrieved product: %+v\n", product)
    } else {
        fmt.Printf("Product %s not found\n", id)
    }
}
```

## Key Features

### 1. Automatic Expiration

TimeCache automatically removes items after their TTL expires. You don't need to manually manage item expiration.

### 2. Thread-Safe Operations

All operations in TimeCache are thread-safe, making it suitable for concurrent access in multi-goroutine environments.

### 3. Customizable TTL

You can set a global TTL when creating the cache, and also update TTL for individual items:

```go
productCache.UpdateTTL("P002", 30*time.Minute)
```

### 4. Cache Statistics

TimeCache provides built-in statistics to help you monitor its performance:

```go
hits, misses, expired := productCache.Stats()
fmt.Printf("Cache stats - Hits: %d, Misses: %d, Expired: %d\n", hits, misses, expired)
```

### 5. Flexible Key-Value Storage

You can store any type of data in TimeCache, as long as it can be cast to and from `interface{}`.

## Best Practices

1. **Choose an Appropriate TTL**: Set a TTL that balances data freshness with performance gains.
2. **Handle Cache Misses Gracefully**: Always have a fallback method to fetch data when it's not in the cache.
3. **Use Cache Statistics**: Regularly monitor cache performance to optimize your caching strategy.
4. **Implement Cache Invalidation**: For data that changes frequently, implement a strategy to invalidate or update cached items.

## Conclusion

TimeCache provides a simple yet powerful way to implement caching in your Go applications. By reducing database load and improving response times, it can significantly enhance your application's performance.

For more detailed information and advanced usage, please refer to our [API documentation](#).

Happy caching!