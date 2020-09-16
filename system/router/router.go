package router

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/labstack/echo"
)

// Router handles all the routes of the API
type Router struct {
	Echo *echo.Echo
}

//NewRouter initialises the router
func NewRouter() Router {
	return Router{
		Echo: echo.New(),
	}
}

// Start will start the http router
func (r *Router) Start() {
	r.registerRoutes()
	r.Echo.Logger.Fatal(r.Echo.Start(":1323"))
}

func (r *Router) registerRoutes() {
	fooType := reflect.TypeOf(r)
	fooVal := reflect.ValueOf(r)
	fmt.Println("--- Registering Routes ---")
	for i := 0; i < fooType.NumMethod(); i++ {
		method := fooType.Method(i)
		if strings.Contains(method.Name, "RegisterRouteFor") {
			fmt.Println(strings.Replace(method.Name, "RegisterRouteFor", "", 1))
			method.Func.Call([]reflect.Value{fooVal})
		}
	}
}
