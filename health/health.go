package health

import (
	"context"
	"go-crud/config"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthStatus string

const (
	Healthy   HealthStatus = "healthy"
	Unhealthy HealthStatus = "unhealthy"
	Degraded  HealthStatus = "degraded"
)

type HealthCheck struct {
	Status    HealthStatus         `json:"status"`
	Timestamp string              `json:"timestamp"`
	Version   string              `json:"version"`
	Uptime    string              `json:"uptime"`
	Checks    map[string]Check    `json:"checks"`
}

type Check struct {
	Status    HealthStatus `json:"status"`
	Message   string       `json:"message,omitempty"`
	Duration  string       `json:"duration,omitempty"`
	Error     string       `json:"error,omitempty"`
}

type HealthChecker struct {
	version string
	startTime time.Time
	checks  map[string]func() Check
}

func NewHealthChecker(version string) *HealthChecker {
	return &HealthChecker{
		version:   version,
		startTime: time.Now(),
		checks:    make(map[string]func() Check),
	}
}

func (hc *HealthChecker) AddCheck(name string, checkFunc func() Check) {
	hc.checks[name] = checkFunc
}

func (hc *HealthChecker) GetHealth() HealthCheck {
	start := time.Now()
	checks := make(map[string]Check)
	overallStatus := Healthy
	
	for name, checkFunc := range hc.checks {
		checkStart := time.Now()
		check := checkFunc()
		check.Duration = time.Since(checkStart).String()
		checks[name] = check
		
		if check.Status == Unhealthy {
			overallStatus = Unhealthy
		} else if check.Status == Degraded && overallStatus != Unhealthy {
			overallStatus = Degraded
		}
	}
	
	uptime := time.Since(hc.startTime)
	
	return HealthCheck{
		Status:    overallStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   hc.version,
		Uptime:    uptime.String(),
		Checks:    checks,
	}
}

func (hc *HealthChecker) ServeHTTP(c *gin.Context) {
	health := hc.GetHealth()
	
	statusCode := 200
	if health.Status == Unhealthy {
		statusCode = 503
	} else if health.Status == Degraded {
		statusCode = 200
	}
	
	c.JSON(statusCode, health)
}

// Database health check
func DatabaseHealthCheck() Check {
	start := time.Now()
	
	if config.DB == nil {
		return Check{
			Status:  Unhealthy,
			Message: "Database connection not initialized",
			Error:   "Database not initialized",
		}
	}
	
	// Test database connection
	sqlDB, err := config.DB.DB()
	if err != nil {
		return Check{
			Status:  Unhealthy,
			Message: "Failed to get database connection",
			Error:   err.Error(),
		}
	}
	
	// Ping database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := sqlDB.PingContext(ctx); err != nil {
		return Check{
			Status:  Unhealthy,
			Message: "Database ping failed",
			Error:   err.Error(),
		}
	}
	
	// Check connection pool stats
	stats := sqlDB.Stats()
	if stats.OpenConnections > stats.MaxOpenConnections*0.9 {
		return Check{
			Status:  Degraded,
			Message: "Database connection pool is nearly full",
		}
	}
	
	duration := time.Since(start)
	
	return Check{
		Status:  Healthy,
		Message: "Database connection is healthy",
		Duration: duration.String(),
	}
}

// Memory health check
func MemoryHealthCheck() Check {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Check if memory usage is too high
	memUsagePercent := float64(m.Alloc) / float64(m.Sys) * 100
	
	if memUsagePercent > 90 {
		return Check{
			Status:  Unhealthy,
			Message: "Memory usage is critically high",
		}
	} else if memUsagePercent > 80 {
		return Check{
			Status:  Degraded,
			Message: "Memory usage is high",
		}
	}
	
	return Check{
		Status:  Healthy,
		Message: "Memory usage is normal",
	}
}

// Application health check
func ApplicationHealthCheck() Check {
	// Check if the application is responsive
	// This could include checking critical services, configurations, etc.
	
	return Check{
		Status:  Healthy,
		Message: "Application is running normally",
	}
}

// Health check endpoints
func SetupHealthRoutes(r *gin.Engine, version string) {
	healthChecker := NewHealthChecker(version)
	
	// Add health checks
	healthChecker.AddCheck("database", DatabaseHealthCheck)
	healthChecker.AddCheck("memory", MemoryHealthCheck)
	healthChecker.AddCheck("application", ApplicationHealthCheck)
	
	// Health check endpoint
	r.GET("/health", healthChecker.ServeHTTP)
	
	// Liveness probe (simple check)
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "alive",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	
	// Readiness probe (detailed check)
	r.GET("/health/ready", healthChecker.ServeHTTP)
}
