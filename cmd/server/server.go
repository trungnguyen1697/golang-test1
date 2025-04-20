package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang-test1/pkg/config"
	"golang-test1/pkg/middleware"
	"golang-test1/pkg/route"
	"golang-test1/platform/database"
	"golang-test1/platform/logger"
	"os"
	"os/signal"
	"syscall"
)

func Serve() {
	appCfg := config.AppCfg()

	logger.SetUpLogger()
	logr := logger.GetLogger()

	// connect to DB
	if err := database.ConnectDB(); err != nil {
		logr.Panicf("failed database setup. error: %v", err)
	}

	// Define Fiber config & app.
	fiberCfg := config.FiberConfig()
	app := fiber.New(fiberCfg)

	// Attach Middlewares.
	middleware.FiberMiddleware(app)

	// Routes.
	route.GeneralRoute(app)
	route.SwaggerRoute(app)
	route.PublicRoutes(app)
	route.PrivateRoutes(app)
	route.NotFoundRoute(app)

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// start shutdown goroutine
	go func() {
		// capture sigterm and other system call here
		<-sigCh
		logr.Infoln("Shutting down server...")
		_ = app.Shutdown()
	}()

	// start http server
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		logr.Errorf("Oops... server is not running! error: %v", err)
	}

}
