package server

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"github.com/lukecarr/trophies/frontend"
	sql "github.com/lukecarr/trophies/internal/db"
	"github.com/lukecarr/trophies/internal/env"
	"github.com/lukecarr/trophies/internal/routes"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Router *httprouter.Router
	Env    *env.Env
}

func New(dsn, npsso, rawg string) *Server {
	db, err := sql.New(dsn)

	if err != nil {
		log.Fatalln("Failed to open SQLite connection", err)
	} else {
		log.Println("Established connection with SQLite!")
	}

	// In-memory mode
	if dsn == sql.MemoryDsn {
		if _, ok := os.LookupEnv("DISABLE_IN_MEMORY_WARN"); !ok {
			log.Println("WARNING: Launching in in-memory mode as 'DSN' environment variable wasn't set. Data will be lost on shutdown!")
		}

		n, migrateErr := sql.Migrate(db.Sql)

		if err != nil {
			log.Fatalln("Failed to perform migrations on in-memory database", migrateErr)
		} else {
			log.Printf("Applied %d migration(s) successfully!\n", n)
		}
	}

	router := httprouter.New()

	srv := &Server{
		Router: router,
		Env:    env.New(db),
	}
	srv.Env.Services = env.NewServices(db, rawg)

	srv.Router.GET("/api/version", routes.Version())
	routes.Game(srv.Env, srv.Router)

	frontendFS, err := fs.Sub(frontend.Static, "dist")
	if err != nil {
		panic(err)
	}
	fsServer := http.FileServer(http.FS(frontendFS))

	// Read index.html into memory
	indexFile, err := frontendFS.Open("index.html")
	if err != nil {
		panic(err)
	}
	indexContent, err := io.ReadAll(indexFile)
	if err != nil {
		panic(err)
	}
	indexFile.Close()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := frontendFS.Open(r.URL.Path[1:])

		if err == nil {
			f.Close()
			fsServer.ServeHTTP(w, r)
			return
		}

		http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(indexContent))
	})

	return srv
}

func (s *Server) Listen(addr string) {
	done := make(chan bool)

	go func() {
		log.Fatal(http.ListenAndServe(addr, s.Router))
	}()
	log.Printf("Listening on %v\n", addr)
	<-done
}
