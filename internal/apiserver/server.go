package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/assanoff/image-resizer/internal/imageresizer"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
}

func newServer() *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/", s.handlerRoot()).Methods("GET")
	s.router.HandleFunc("/resize", s.handleImageResize()).Methods("POST")
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) handlerRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{"ok": "I'm working"})
	}
}

func (s *server) handleImageResize() http.HandlerFunc {
	type message struct {
		ID     string `json:"id"`
		Size   uint   `json:"size"`
		Body   string `json:"body"`
		Format string `json:"format"`
	}
	req := &message{}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// resBody, err := imageresizer.DecodeImageBase64(req.Body, req.Size, req.Size, req.Format)
		decodedImage, err := imageresizer.DecodeImageBase64(req.Body)
		if err != nil {
			fmt.Errorf("could not decode image form base64", err)
		}

		thumbnail, err := imageresizer.ResizeImage(decodedImage, req.Size, req.Size, req.Format)
		if err != nil {
			fmt.Errorf("could not resize image", err)
		}

		encodedImage, err := imageresizer.EncodeImageBase64(thumbnail)
		if err != nil {
			fmt.Errorf("could not encode image to base64", err)
		}

		resp := &message{}
		resp.ID = req.ID
		resp.Size = req.Size
		resp.Format = req.Format
		resp.Body = encodedImage

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, resp)

	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
