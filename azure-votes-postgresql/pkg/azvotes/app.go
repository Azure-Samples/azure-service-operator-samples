package azvotes

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type VoteServer struct {
	db *Client
}

func NewVoteServer(db *Client) *VoteServer {
	return &VoteServer{db}
}

func (s *VoteServer) Routes() {
	http.HandleFunc("/", s.handleIndex())
	http.HandleFunc("/ping", s.handlePing())
}

func (s *VoteServer) Start() error {
	s.Routes()
	return http.ListenAndServe(":8080", nil)
}

func (s *VoteServer) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := s.db.Ping(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("server is alive"))
	}
}

func (s *VoteServer) handleIndex() http.HandlerFunc {

	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	type Vals struct {
		Button1 string
		Button2 string
		Value1  int
		Value2  int
		Title   string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.New("tmpl").Parse(templ)
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		v := Vals{
			Title:   "Azure Votes",
			Button1: "dogs",
			Button2: "cats",
		}
		log.Println(r.Method)

		if val, err := s.db.CountVote(v.Button1); err == nil {
			v.Value1 = val
		} else {
			log.Println(err)
		}
		if val, err := s.db.CountVote(v.Button2); err == nil {
			v.Value2 = val
		} else {
			log.Println(err)
		}

		if r.Method == "POST" {
			vote := r.FormValue("vote")
			if vote == "" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if vote == "reset" {
				_, err = s.db.DeleteVotes()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				v.Value1 = 0
				v.Value2 = 0
			} else {
				_, err := s.db.CreateVote(vote)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if v.Button1 == vote {
					v.Value1 = v.Value1 + 1
				}
				if v.Button2 == vote {
					v.Value2 = v.Value2 + 1
				}
			}

		}

		err = tpl.Execute(w, v)
		if err != nil {
			panic(err)
		}
		return

	}
}
