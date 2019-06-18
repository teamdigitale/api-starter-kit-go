package api

import (
	"fmt"
	"net/http"
	"time"
	"math/rand"

	"github.com/labstack/echo/v4"
)


// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendPetstoreError(ctx echo.Context, code int, message string) error {
	err := ctx.JSON(code, message)
	return err
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

//
// Implement the methods declared in the generated interface.
//


func (app *MyApplication) GetEcho(ctx echo.Context) error {

    var ts = time.Now()
	result := Timestamps{
	    Ts: ts,
	}
	return ctx.JSON(http.StatusOK, result)
}

func (app *MyApplication) GetStatus(ctx echo.Context) error {

	ctx.Response().Header().Set("Cache-Control", "no-store")

    var f = app.r.Int() % 5
	fmt.Println("string %d", f)
    if (f < 3) {
        var result = Problem {
            Status: int32(200),
            Title: "ok",
        }
        return ctx.JSON(http.StatusOK, result)
	}
	var result = Problem {
		Status: int32(503),
		Title: "ko",
	}

	return ctx.JSON(http.StatusServiceUnavailable, result)
}