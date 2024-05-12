package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"VEEEKTOR_api/internal/service"
)

var apiPrefix = "/api"

func Start() {
	log.Printf("VEEEKTOR_api is starting...")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mux := getMultiplexer()
	server := &http.Server{Addr: ":8080", Handler: mux}
	defer server.Close()

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func getMultiplexer() *http.ServeMux {
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc(apiPrefix+"/auth/refresh", service.UpdateToken)
	mux.HandleFunc(apiPrefix+"/auth/logout", service.Logout)

	// Users
	mux.HandleFunc(apiPrefix+"/users", service.GetUsersHandler)
	mux.HandleFunc(apiPrefix+"/users/signin", service.UsersSignInHandler)
	mux.HandleFunc(apiPrefix+"/users/signup", service.UsersSignUpHandler)

	// Educational envs
	mux.HandleFunc(apiPrefix+"/educational_envs", service.GetEducatinalEnvironmentsHandler)

	// Departments
	mux.HandleFunc(apiPrefix+"/departments", service.GetDepartmentsHandler)

	// Groups
	mux.HandleFunc(apiPrefix+"/groups", service.GetGroupsHandler)

	// Courses
	mux.HandleFunc(apiPrefix+"/courses", service.GetCouresesHandler)
	mux.HandleFunc(apiPrefix+"/courses/infos", service.GetNestedInfosHandler)
	mux.HandleFunc(apiPrefix+"/courses/labs", service.GetNestedLabsHandler)
	mux.HandleFunc(apiPrefix+"/courses/tests", service.GetNestedTestsHandler)

	return mux
}
