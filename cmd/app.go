package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/pflag"
	"github.com/stevenxie/merlin/internal/info"
	"github.com/stevenxie/merlin/pkg/api"
	ess "github.com/unixpickle/essentials"
)

// opts are a set of program options.
var opts struct {
	ShowVersion bool
	ShowHelp    bool
	Port        int
}

// Define CLI flags, initialize program.
func init() {
	pflag.BoolVarP(&opts.ShowHelp, "help", "h", false, "Show help (usage).")
	pflag.BoolVarP(&opts.ShowVersion, "version", "v", false, "Show version.")
	pflag.IntVarP(&opts.Port, "port", "p", 3000, "Port to listen on.")

	loadEnv() // load .env variables
	pflag.Parse()
}

// Exec is the entrypoint to command rgv.
func Exec() {
	if opts.ShowHelp {
		showHelp()
		os.Exit(0)
	}
	if opts.ShowVersion {
		fmt.Println(info.Version)
		os.Exit(0)
	}

	// Create program logger.
	logger, err := buildLogger()
	if err != nil {
		ess.Die("Error while building zap.SugaredLogger:", err)
	}

	// Create and run server.
	server, err := api.NewServer(logger)
	if err != nil {
		ess.Die("Error while building server:", err)
	}

	addr := fmt.Sprintf(":%d", opts.Port)
	fmt.Printf("Listening on address '%s'...\n", addr)

	go shutdownUponInterrupt(server)
	err = server.ListenAndServe(addr)
	if (err != nil) && (err != http.ErrServerClosed) {
		ess.Die("Error while starting server:", err)
	}
}

func shutdownUponInterrupt(s *api.Server) {
	const timeout = 1 * time.Second

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)

	<-ch // wait for a signal
	fmt.Printf("Shutting down server gracefully (timeout: %s)...\n", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		ess.Die("Error during server shutdown:", err)
	}
}
