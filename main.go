package main

import (
	"bufio"
	fmt "github.com/jhunt/go-ansi"
	"net/http"
	"encoding/json"
	"os"
	"regexp"
	"os/exec"
	"strings"
	"io"
	"io/ioutil"

	"github.com/jhunt/go-log"
)

type Song struct {
	Active int `json:"active"`
	File   string `json:"file"`
}

func ReadM3u(path string) []Song {
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("unable to open %s for reading: %s", path, err)
		return nil
	}
	defer f.Close()

	songs := make([]Song, 0)
	pat := regexp.MustCompile(`^(#\s*)?(.+)$`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		m := pat.FindStringSubmatch(s)
		if m != nil {
			active := 1
			if m[1] != "" {
				active = 0
			}
			songs = append(songs, Song{
				Active: active,
				File:   m[2],
			})
		}
	}

	log.Infof("read a playlist of %d tracks", len(songs))
	return songs
}

func WriteM3u(out io.Writer, songs []Song) {
	for _, song := range songs {
		if song.Active == 1 {
			fmt.Fprintf(out, "%s\n", song.File)
		} else {
			fmt.Fprintf(out, "#%s\n", song.File)
		}
	}

	log.Infof("wrote a playlist of %d tracks", len(songs))
}

func main() {
	log.SetupLogging(log.LogConfig{
		Type: "console",
		Level: "info",
	})
	log.Infof("mixbooth starting up...")

	Playlist := fmt.Sprintf("%s/playlist.m3u", os.Getenv("RADIO_ROOT"))

	http.HandleFunc("/playlist", func (w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if req.Method == "GET" {
			var out struct {
				Playlist []Song `json:"playlist"`
			}

			out.Playlist = ReadM3u(Playlist)
			b, err := json.Marshal(out)
			if err != nil {
				log.Errorf("unable to marshal json response to GET /playlist: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"internal error"}`)
				return
			}
			w.WriteHeader(200)
			w.Write(b)

		} else if req.Method == "PUT" {
			var in struct {
				Playlist []Song `json:"playlist"`
			}
			defer req.Body.Close()

			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Errorf("unable to read body paylod in PUT /playlist: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"internal error"}`)
				return
			}
			if err := json.Unmarshal(b, &in); err != nil {
				log.Errorf("unable to unmarshal json body paylod in PUT /playlist: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"internal error"}`)
				return
			}

			f, err := os.OpenFile(Playlist+"~", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Errorf("unable to open temporary playlist %s~ for writing: %s", Playlist, err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"internal error"}`)
				return
			}
			defer f.Close()

			WriteM3u(f, in.Playlist)
			w.WriteHeader(200)
			fmt.Fprintf(w, `{"ok":"updated"}`)

		} else {
			w.WriteHeader(404)
			fmt.Fprintf(w, `{"error":"not found"}`)
		}
	})

	http.HandleFunc("/upload", func (w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if req.Method == "POST" {
			var in struct {
				URL string `json:"url"`
			}

			defer req.Body.Close()

			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Errorf("unable to read body paylod in POST /upload: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"internal error"}`)
				return
			}
			if err := json.Unmarshal(b, &in); err != nil {
				log.Errorf("unable to unmarshal json body paylod in POST /upload: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"internal error"}`)
				return
			}

			log.Infof("ingesting %s from upstream", in.URL)

			cmd := exec.Command("ingest", in.URL)
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Errorf("unable to run ingest process: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"internal error"}`)
				return
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				log.Errorf("unable to run ingest process: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"internal error"}`)
				return
			}

			if err := cmd.Start(); err != nil {
				log.Errorf("unable to run ingest process: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"ingestion failed"}`)
				return
			}

			go io.Copy(os.Stdout, stdout)
			go io.Copy(os.Stderr, stderr)

			if err := cmd.Wait(); err != nil {
				log.Errorf("unable to run ingest process: %s", err)
				w.WriteHeader(500)
				fmt.Fprintf(w, `{"error":"ingestion failed"}`)
				return
			}

			w.WriteHeader(200)
			fmt.Fprintf(w, `{"ok":"ingested"}`)

		} else {
			w.WriteHeader(404)
			fmt.Fprintf(w, `{"error":"not found"}`)
		}
	})

	http.Handle("/", http.FileServer(http.Dir(os.Getenv("HTDOCS_ROOT"))))

	log.Infof("mixbooth listening on :5000")
	http.ListenAndServe(":5000", nil)
	log.Infof("mixbooth shut down")
}
