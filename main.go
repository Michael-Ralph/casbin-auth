package main

import (
	"net/http"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Define a custom middleware for Casbin
func CasbinAuth(e *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role from context or session
			// Just an example, will replace with our auth logic
			user := getUserFromContext(c)

			// Get request path and method
			path := c.Request().URL.Path
			method := c.Request().Method
			c.Echo().Logger.Printf("Entering CasbinAuth: user=%s, path=%s, method=%s", user, path, method)

			// Check permission using Casbin enforcer
			if ok, err := e.Enforce(user, path, method); err != nil {
				c.Echo().Logger.Printf("Casbin error: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Authorization failed: "+err.Error())
			} else if !ok {
				c.Echo().Logger.Printf("Denied: user=%s, path=%s, method=%s", user, path, method)
				return echo.NewHTTPError(http.StatusForbidden, "Access denied")
			} else {
				c.Echo().Logger.Printf("Allowed: user=%s, path=%s, method=%s", user, path, method)
			}
			// If autherized, proceed to the next middleware or handler
			return next(c)
		}
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Mock function to get user from context. Will replace when I put this into our code
func getUserFromContext(c echo.Context) string {
	// We will get this from our session
	// For the example I will use the header

	user := c.Request().Header.Get("X-User")
	if user == "" {
		// Default to guest if not specified
		user = "guest"
	}
	return user
}

func main() {
	// Initialize echo
	e := echo.New()

	e.Logger.SetOutput(os.Stdout) // Force logs to stdout

	// Add basic middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize Casbin enforcer with model and policy files
	enforcer, err := casbin.NewEnforcer("auth_model.conf", "policy.csv")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Public routes (no auth required)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Create a group for protected routes
	admin := e.Group("/admin")

	// Apply Casbin middleware to admin group
	admin.Use(CasbinAuth(enforcer))

	// Admin routes (protected by Casbin)
	admin.GET("/dashboard", func(c echo.Context) error {
		return c.String(http.StatusOK, "Admin Dashboard")
	})

	admin.GET("/posts", func(c echo.Context) error {
		return c.String(http.StatusOK, "View post")
	})

	admin.POST("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "User created")
	})

	admin.POST("/posts", func(c echo.Context) error {
		return c.String(http.StatusOK, "Upload post")
	})

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
