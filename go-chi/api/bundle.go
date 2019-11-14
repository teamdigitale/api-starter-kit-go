package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

)

type Filter = func(http.Handler) http.Handler


// Handler creates http.Handler with routing matching OpenAPI spec.
func HandlerCustom(si ServerInterface, middlewares []Filter ) http.Handler {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,	
	})
	fmt.Println(cors)
	// r.Use(cors.Handler)
	
	r.Group(func(r chi.Router) {
		r.Use(GetEchoCtx)
		r.Get("/echo", si.GetEcho)
	})
	r.Group(func(r chi.Router) {
		r.Use(OptionsEchoCtx)
		r.Options("/echo", si.OptionsEcho)
	})
	r.Group(func(r chi.Router) {
		r.Use(GetStatusCtx)
		r.Get("/status", si.GetStatus)
	})

	return r
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendPetstoreError(w http.ResponseWriter, code int, message string) {
	problem := Problem{
		Status: int32(code),
		Title:  "Internal Server Error",
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(problem)
}

/**
 * This structure contains all shared data for this app,
 * in our case it's a random number generator.
 */
type MyApplication struct {
	r *rand.Rand
}

func CreateApplication() *MyApplication {
	return &MyApplication{
		r: rand.New(rand.NewSource(99)),
	}
}

// Define a cors handler.
func CORSFilter() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	})
}

//
// Implement the methods declared in the generated interface.
//

func (app *MyApplication) OptionsEcho(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	w.WriteHeader(http.StatusNoContent)
}

func (app *MyApplication) GetEcho(w http.ResponseWriter, r *http.Request) {

	var ts = time.Now()
	result := Timestamps{
		Ts: ts,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (app *MyApplication) GetStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-store")

	var f = app.r.Int() % 5
	fmt.Println("string %d", f)
	if f < 3 {
		var result = Problem{
			Status: int32(200),
			Title:  "ok",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	}
	var result = Problem{
		Status: int32(503),
		Title:  "ko",
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(result)
}

/*


func ProblemErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

	httpError, ok := err.(*echo.HTTPError)
	fmt.Println("problemErrorHandler: ", err, ok, httpError)
	if !ok {
		return
	}

	var problem Problem
	problem.Status = int32(httpError.Code)
	problem.Title, ok = httpError.Message.(string)
	if !ok {
		c.Logger().Error("Message is not a string", httpError)
		return
	}

	if httpError.Message == "Path was not found" {
		problem.Status = 404
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(httpError.Code)
		} else {
			c.Response().Header().Set("Content-Type", "application/problem+json")
			err = c.JSON(int(problem.Status), problem)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}

}
*/
