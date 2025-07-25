# ErmesLog

A structured logging library for Go microservices built on top of logrus.

## Installation

```bash
go get github.com/galvushow/loggo
```

## Quick Start

```go
package main

import (
    "github.com/galvushow/loggo"
    "github.com/sirupsen/logrus"
)

func main() {
    // Simple usage
    logger := ermeslog.NewDefault("mycompany", "api-service")
    logger.Info("Service started")
    
    // Advanced configuration
    logger = ermeslog.New(ermeslog.Config{
        Business:    "mycompany",
        Service:     "api-service",
        Version:     "1.0.0",
        Environment: "production",
        Level:       logrus.InfoLevel,
        FileOutput: &ermeslog.FileConfig{
            Filename:   "/var/log/api.log",
            MaxSize:    100,
            MaxBackups: 3,
            MaxAge:     28,
            Compress:   true,
        },
    })
    
    // With fields
    logger.WithFields(map[string]interface{}{
        "port": 8080,
        "tls":  true,
    }).Info("Server listening")
    
    // Error logging
    err := someFunction()
    if err != nil {
        logger.Error(err, "Failed to process")
    }
}
```

## HTTP Middleware

```go
func main() {
    logger := ermeslog.NewDefault("mycompany", "api")
    
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    
    // Add logging middleware
    handler := logger.HTTPMiddleware(mux)
    
    http.ListenAndServe(":8080", handler)
}
```

## Features

- Structured JSON logging
- Automatic request ID generation
- Context propagation
- File rotation support
- HTTP middleware
- Multiple output targets
- Environment-based formatting

## License

MIT