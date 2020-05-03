package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss-clientstate"
	_ "github.com/volatiletech/authboss/auth"
	"github.com/volatiletech/authboss/defaults"
	_ "github.com/volatiletech/authboss/register"
	"github.com/volatiletech/authboss/remember"
	"go.uber.org/zap"
)

var (
	ab           = authboss.New()
	database     = NewMemStorer()
	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

func setupAuthboss() {
	ab.Config.Paths.RootURL = "http://localhost:5000"
	ab.Config.Modules.LogoutMethod = "GET"

	// Set up our server, session and cookie storage mechanisms.
	// These are all from this package since the burden is on the
	// implementer for these.
	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	// Another piece that we're responsible for: Rendering views.
	// Though note that we're using the authboss-renderer package
	// that makes the normal thing a bit easier.
	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}

	// We render mail with the authboss-renderer but we use a LogMailer
	// which simply sends the e-mail to stdout.
	//ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "ab_views")

	// The preserve fields are things we don't want to
	// lose when we're doing user registration (prevents having
	// to type them again)
	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	// TOTP2FAIssuer is the name of the issuer we use for totp 2fa
	//ab.Config.Modules.TOTP2FAIssuer = "ABBlog"
	//ab.Config.Modules.ResponseOnUnauthed = authboss.RespondRedirect

	// Turn on e-mail authentication required
	//ab.Config.Modules.TwoFactorEmailAuthRequired = true

	// This instantiates and uses every default implementation
	// in the Config.Core area that exist in the defaults package.
	// Just a convenient helper if you don't want to do anything fancy.
	defaults.SetCore(&ab.Config, true, false)

	// Here we initialize the bodyreader as something customized in order to accept a name
	// parameter for our user as well as the standard e-mail and password.
	//
	// We also change the validation for these fields
	// to be something less secure so that we can use test data easier.
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
		//Confirms: map[string][]string{
		//	"register":    {"password", authboss.ConfirmPrefix + "password"},
		//	"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		//},
		//Whitelist: map[string][]string{
		//	"register": {"IIN", "phone", "password"},
		//},
	}

	//oauthcreds := struct {
	//	ClientID     string `toml:"client_id"`
	//	ClientSecret string `toml:"client_secret"`
	//}{}

	//// Set up 2fa
	//twofaRecovery := &twofactor.Recovery{Authboss: ab}
	//if err := twofaRecovery.Setup(); err != nil {
	//	panic(err)
	//}

	//totp := &totp2fa.TOTP{Authboss: ab}
	//if err := totp.Setup(); err != nil {
	//	panic(err)
	//}

	//sms := &sms2fa.SMS{Authboss: ab, Sender: smsLogSender{}}
	//if err := sms.Setup(); err != nil {
	//	panic(err)
	//}

	// Set up Google OAuth2 if we have credentials in the
	// file oauth2.toml for it.
	//_, err := toml.DecodeFile("oauth2.toml", &oauthcreds)
	//if err == nil && len(oauthcreds.ClientID) != 0 && len(oauthcreds.ClientSecret) != 0 {
	//	fmt.Println("oauth2.toml exists, configuring google oauth2")
	//	ab.Config.Modules.OAuth2Providers = map[string]authboss.OAuth2Provider{
	//		"google": {
	//			OAuth2Config: &oauth2.Config{
	//				ClientID:     oauthcreds.ClientID,
	//				ClientSecret: oauthcreds.ClientSecret,
	//				Scopes:       []string{`profile`, `email`},
	//				Endpoint:     google.Endpoint,
	//			},
	//			FindUserDetails: aboauth.GoogleUserDetails,
	//		},
	//	}
	//} else if os.IsNotExist(err) {
	//	fmt.Println("oauth2.toml doesn't exist, not registering oauth2 handling")
	//} else {
	//	fmt.Println("error loading oauth2.toml:", err)
	//}

	//err := auth.Auth{}.Init(ab)
	//if err != nil {
	//	panic(err)
	//}

	// Initialize authboss (instantiate modules etc.)
	if err := ab.Init(); err != nil {
		panic(err)
	}
}

const (
	sessionCookieName = "ab_blog"
)

func main() {

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, err := cfg.Build()
	if err != nil {
		log.Fatal("error creating log", err)
	}
	l.Info("starting service")

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

	// Initialize authboss
	setupAuthboss()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	//
	//router := gin.New()
	//router.Use(gin.Logger(), ab.LoadClientStateMiddleware, remember.Middleware(ab))

	mux := chi.NewRouter()
	// The middlewares we're using:
	// - logger just does basic logging of requests and debug info
	// - nosurfing is a more verbose wrapper around csrf handling
	// - LoadClientStateMiddleware is required for session/cookie stuff
	// - remember middleware logs users in if they have a remember token
	// - dataInjector is for putting data into the request context we need for our template layout
	mux.Use(middleware.Logger, ab.LoadClientStateMiddleware, remember.Middleware(ab))

	// Authed routes
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized)) //, lock.Middleware(ab), confirm.Middleware(ab))
		mux.MethodFunc("GET", "/blogs/new", stubHandler)
		mux.MethodFunc("GET", "/blogs/{id}/edit", stubHandler)
		mux.MethodFunc("POST", "/blogs/{id}/edit", stubHandler)
		mux.MethodFunc("POST", "/blogs/new", stubHandler)
		// This should actually be a DELETE but can't be bothered to make a proper
		// destroy link using javascript atm. See where AB allows you to configure
		// the logout HTTP method.
		mux.MethodFunc("GET", "/blogs/{id}/destroy", stubHandler)
	})

	// Routes
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.ModuleListMiddleware(ab))
		mux.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})
	mux.Get("/blogs", stubHandler)
	mux.Get("/", stubHandler)

	err = http.ListenAndServe(":"+port, mux)

	//router.GET("/", func(c *gin.Context) {
	//	c.String(200, "hello world")
	//})
	//router.POST("/register", func(c *gin.Context) {
	//	c.String(200, "register")
	//})
	//router.GET("/login", func(c *gin.Context) {
	//	c.String(200, "login")
	//})
	//router.GET("/logout", func(c *gin.Context) {
	//	c.String(200, "logout")
	//})
	//
	//err = router.Run(":" + port)
	if err != nil {
		l.Error("closing server", zap.Error(err))
	}
}

func stubHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("/stubHandler"))
}
