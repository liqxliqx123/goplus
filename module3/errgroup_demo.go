package module3

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func run(ctx context.Context, addr string, mux *http.ServeMux) error {
	s := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		fmt.Println("ready to shutdown")
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}

func sig(ctx context.Context) (err error) {
	q := make(chan os.Signal)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sigInfo := <-q:
		err = fmt.Errorf("%v", sigInfo)
	case <-ctx.Done():
		err = ctx.Err()
	}
	return err
}

//RunServerWithErrGroup run a http server with errgroup
func RunServerWithErrGroup() {
	group, ctx := errgroup.WithContext(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	group.Go(func() error {
		return run(ctx, ":8000", mux)
	})

	group.Go(func() error {
		return sig(ctx)
	})

	log.Println(group.Wait())
}
