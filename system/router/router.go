package router

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/labstack/echo"
)

// Router handles all the routes of the API
type Router struct {
	Echo   *echo.Echo
	DBConn *sql.DB
}

//NewRouter initialises the router
func NewRouter(db *sql.DB) Router {
	return Router{
		Echo: echo.New(),
		DBConn: db
	}
}

// Start will start the http router
func (r *Router) Start() {
	r.registerRoutes()
	r.Echo.Logger.Fatal(r.Echo.Start(":1323"))
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
