package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/digitalCitizenship/lib/storage/mock"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss-clientstate"
	_ "github.com/volatiletech/authboss/auth"
	"github.com/volatiletech/authboss/defaults"
	_ "github.com/volatiletech/authboss/logout"
	_ "github.com/volatiletech/authboss/register"
	"github.com/volatiletech/authboss/remember"
	"go.uber.org/zap"
)

func setupAuthBoss(ab *authboss.Authboss, database authboss.CreatingServerStorer) {

	const sessionCookieName = "ab_blog"

	var (
		sessionStore abclientstate.SessionStorer
		cookieStore  abclientstate.CookieStorer
	)

	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	sessionStore = abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))

	ab.Config.Paths.RootURL = "http://localhost:5000"
	ab.Config.Modules.LogoutMethod = "DELETE"

	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}

	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	defaults.SetCore(&ab.Config, true, true)

	emailRule := defaults.Rules{
		FieldName: "IIN", Required: true,
		MinLength: 12, MaxLength: 12,
		MatchError:      "Must be a valid IIN",
		AllowWhitespace: false,
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength:       10,
		MatchError:      "Must be a len with > 10",
		AllowWhitespace: false,
	}
	nameRule := defaults.Rules{
		FieldName: "phone", Required: true,
		MinLength: 11, MaxLength: 11,
		AllowWhitespace: false,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON:    true,
		UseUsername: true,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule},
			"recover_end": {passwordRule},
		},
	}

	if err := ab.Init(); err != nil {
		panic(err)
	}
}

func main() {

	//TODO: logout не работает
	//TODO: переписать auth service

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, err := cfg.Build()
	if err != nil {
		log.Fatal("error creating log", err)
	}
	l.Info("starting service")

	ab := authboss.New()
	database := mock.New()
	setupAuthBoss(ab, database)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Logger, ab.LoadClientStateMiddleware, remember.Middleware(ab))

	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.Middleware2(ab, authboss.RequireFullAuth, authboss.RespondUnauthorized))
		mux.MethodFunc("POST", "/blogs/new", stubHandler)
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.ModuleListMiddleware(ab))
		mux.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})
	mux.Get("/", stubHandler)

	mux.Get("/infected/list", getInfectedList(database))
	mux.Get("/infected/new", newInfetcted(database))
	mux.Get("/interactions/status", interactedWithInfected(database))
	mux.Get("/interactions/new", newInteraction(database))

	//func newInteraction
	//func interactedWithInfected
	//func getInfectedList
	//func newInfetcted
	//AddInteraction(firstUserID, secondUserID int64, at int64) error
	//InteractedWithInfected(userID int64) (bool, error)
	//GetInfectedList() ([]int64, error)
	//AddInfected(userID int64) error

	err = http.ListenAndServe(":"+port, mux)

	if err != nil {
		l.Error("closing server", zap.Error(err))
	}
}
