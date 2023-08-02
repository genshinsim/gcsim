package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

type UserStore interface {
	Create(id string, name string, ctx context.Context) error //create a new user; name is the discord tag
	Read(id string, ctx context.Context) ([]byte, error)      //"user" should be set in ctx for auth purpose
	UpdateData(data []byte, ctx context.Context) error        //"user" should be set in ctx for auth purpose
	Has(id string, ctx context.Context) (bool, error)         //return true if user exists
}

type RoleChecker interface {
	UserHasDBTagRole(userid string, tag model.DBTag) bool
}

var ErrUserAlreadyExists = errors.New("user already exist")
var ErrUserNotFound = errors.New("user not found")
var ErrAccessDenied = errors.New("access denied")
var ErrInvalidRequest = errors.New("invalid request")
var ErrServerError = errors.New("unexpected server error")

type discordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

type claim struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func (s *Server) tokenCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil && err != http.ErrNoCookie {
			s.Log.Infow("error reading cookie", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err == http.ErrNoCookie {
			s.Log.Info("no token cookie sent; skipping")
			next.ServeHTTP(w, r)
			return
		}
		s.Log.Infow("parsing token", "token", c.Value)
		var cl claim
		token, err := jwt.ParseWithClaims(c.Value, &cl, func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.Discord.JWTKey), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				s.Log.Infow("token is not valid ", "err", err)
			} else {
				s.Log.Infow("error parsing token from cookie", "err", err)
			}
			next.ServeHTTP(w, r)
			return
		}
		if token.Valid {
			r = r.WithContext(context.WithValue(r.Context(), UserContextKey, cl.User))
			s.Log.Infow("token ok, user context set", "user", cl.User)
		}
		//do stuff here
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Log.Info("starting login request")
		code := r.Header.Get("x-discord-code")
		if code == "" {
			s.Log.Infow("forbidden, no discord code provided")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return
		}
		redirect := r.Header.Get("x-discord-redirect")
		if redirect == "" {
			s.Log.Infow("bad request, no redirect url provided")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		s.Log.Infow("requesting discord token", "code", code)

		conf := &oauth2.Config{
			RedirectURL:  redirect,
			ClientID:     s.cfg.Discord.ClientID,
			ClientSecret: s.cfg.Discord.ClientSecret,
			Scopes: []string{
				"identity",
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://discord.com/oauth2/authorize",
				TokenURL:  "https://discord.com/api/oauth2/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		}

		s.Log.Infow("starting token exchange")

		token, err := conf.Exchange(r.Context(), code)
		if err != nil {
			s.Log.Errorw("unexpected error exchanging token", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.Log.Infow("discord token received sucessfully", "token", token)

		res, err := conf.Client(r.Context(), token).Get("https://discord.com/api/users/@me")
		if err != nil || res.StatusCode != 200 {
			s.Log.Errorw("unexpected error getting user info", "err", err)
			w.WriteHeader(res.StatusCode)
			return
		}

		// d, err := httputil.DumpResponse(res, true)
		s.Log.Infow("@me requested from discord sucessfully")

		defer res.Body.Close()

		var du discordUser

		err = json.NewDecoder(res.Body).Decode(&du)

		if err != nil {
			s.Log.Errorw("unexpected error decoding user json", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// check if user exists; create if not
		ok, err := s.cfg.UserStore.Has(du.ID, r.Context())
		if err != nil {
			s.Log.Errorw("unexpected error encountered checking if user exists", "err", err, "user", du)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ok {
			//create
			err := s.cfg.UserStore.Create(du.ID, fmt.Sprintf("%v#%v", du.Username, du.Discriminator), r.Context())
			if err != nil {
				s.Log.Errorw("unexpected error encountered reading user", "err", err, "user", du)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		u, err := s.cfg.UserStore.Read(du.ID, context.WithValue(r.Context(), UserContextKey, du.ID))
		if err != nil {
			s.Log.Errorw("unexpected error encountered reading user", "err", err, "user", du)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//create a JWT with user's id
		expirationTime := time.Now().Add(14 * 24 * time.Hour)
		jwtToken := jwt.New(jwt.SigningMethodHS256)
		claims := jwtToken.Claims.(jwt.MapClaims)
		claims["user"] = du.ID
		tokenString, err := jwtToken.SignedString([]byte(s.cfg.Discord.JWTKey))
		if err != nil {
			s.Log.Errorw("unexpected error signing jwt", "err", err, "user", du)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// response should include user data + a JWT token set in cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
		w.WriteHeader(200)
		w.Write(u)

	}
}

func (s *Server) UserSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		d, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Log.Infow("bad request saving user data; body cannot be read", "err", err)
			return
		}
		s.Log.Infow("received request to save user data", "data", string(d), "user", r.Context().Value("user"))
		err = s.cfg.UserStore.UpdateData(d, r.Context())
		switch err {
		case nil:
			w.WriteHeader(http.StatusOK)
		case ErrUserNotFound:
			w.WriteHeader(http.StatusBadRequest)
			s.Log.Warnw("unexpected update from a user that does not exist", "user", r.Context().Value("user"))
		case ErrInvalidRequest:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			s.Log.Errorw("unexpected error updating data", "err", err)
		}
	}
}
