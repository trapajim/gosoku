package router

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// Router handles all the routes of the API
type Router struct {
	Echo       *echo.Echo
	AdminGroup *echo.Group
	DBConn     *sql.DB
	Timeout    time.Duration
}

//NewRouter initialises the router
func NewRouter(db *sql.DB, duration time.Duration) Router {
	echo := echo.New()
	group := echo.Group("/admin")
	return Router{
		Echo:       echo,
		AdminGroup: group,
		DBConn:     db,
		Timeout:    duration,
	}
}

// Start will start the http router
func (r *Router) Start() {
	r.registerRoutes()
	r.Echo.Logger.SetLevel(log.INFO)
	go func() {
		if err := r.Echo.Start(":1323"); err != nil {
			r.Echo.Logger.Info("shutting down the server")
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Echo.Shutdown(ctx); err != nil {
		r.Echo.Logger.Fatal(err)
	}
}

//UseLogger is a short cut to add the logger middleware
func (r *Router) UseLogger() {
	r.Echo.Use(middleware.Logger())
}

//UseGzip is a short cut to add the gzip middleware
func (r *Router) UseGzip() {
	r.Echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
}

func (r *Router) registerRoutes() {
	method := reflect.TypeOf(r)
	val := reflect.ValueOf(r)
	fmt.Println("--- Registering Routes ---")
	for i := 0; i < method.NumMethod(); i++ {
		method := method.Method(i)
		if strings.Contains(method.Name, "RegisterRouteFor") {
			fmt.Println(strings.Replace(method.Name, "RegisterRouteFor", "", 1))
			method.Func.Call([]reflect.Value{val})
		}
	}
}
