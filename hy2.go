package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type User struct {
	ID    uint16
	Token string
}

type Cfg struct {
	Ports  Ports  `json:"ports"`
	TLS    TLS    `json:"tls"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
	HTTPS  bool   `json:"https"`
}

type Ports struct {
	Sub  string `json:"sub"`
	Auth string `json:"auth"`
}

type TLS struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

type Request struct {
	Addr string `json:"addr"`
	Auth string `json:"auth"`
	Tx   int64  `json:"tx"`
}

type Response struct {
	OK bool   `json:"ok"`
	ID string `json:"id"`
	Tx int64  `json:"tx"`
	Rx int64  `json:"rx"`
}

var (
	tokens map[string]string
	req    Request
	resp   Response
	cfg    Cfg
	links  string
)

func Init() {
	f1, e1 := os.Open("config.json")
	if e1 != nil {
		log.Fatalln("Can't open file:", e1)
	}

	f2, e8 := os.Open("users.json")
	if e8 != nil {
		log.Fatalln("Can't open file:", e8)
	}

	f3, e9 := os.Open("links.txt")
	if e9 != nil {
		log.Fatalln("Can't open file:", e9)
	}

	buf1, e2 := io.ReadAll(f1)
	if e2 != nil {
		log.Fatalln("Can't read file:", e2)
	}

	buf2, e10 := io.ReadAll(f2)
	if e10 != nil {
		log.Fatalln("Can't read file:", e10)
	}

	buf3, e11 := io.ReadAll(f3)
	if e11 != nil {
		log.Fatalln("Can't read file:", e11)
	}

	e3 := json.Unmarshal(buf1, &cfg)
	if e3 != nil {
		log.Fatalln("Can't decode file:", e3)
	}

	tokens = make(map[string]string, len(buf2)/15)
	e12 := json.Unmarshal(buf2, &tokens)
	if e12 != nil {
		log.Fatalln("Can't decode file:", e12)
	}

	links = string(buf3)
}

func Sub(w http.ResponseWriter, r *http.Request) {
	profileName := cfg.Name
	encodedTitle := base64.StdEncoding.EncodeToString([]byte(profileName))

	w.Header().Set("Profile-Title", "base64:"+encodedTitle)
	w.Header().Set("Profile-Update-Interval", "12")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	encode := base64.StdEncoding.EncodeToString([]byte(links))

	if _, err := w.Write([]byte(encode)); err != nil {
		log.Println("Write error:", err)
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	e4 := json.NewDecoder(r.Body).Decode(&req)
	if e4 != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, exist := tokens[req.Auth]

	if exist {
		resp = Response{
			OK: true,
			ID: id,
			Tx: 200 * 1024 * 1024,
			Rx: 200 * 1024 * 1024,
		}

		e5 := json.NewEncoder(w).Encode(resp)
		if e5 != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	} else {
		e5 := json.NewEncoder(w).Encode(map[string]bool{"ok": false})
		if e5 != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	Init()

	authMux := http.NewServeMux()
	authMux.HandleFunc("/auth", Auth)

	subMux := http.NewServeMux()
	subMux.HandleFunc("/sub/"+cfg.Secret, Sub)

	go func() {
		log.Println("Auth server listening on", cfg.Ports.Auth)
		if err := http.ListenAndServe(cfg.Ports.Auth, authMux); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Auth server crashed:", err)
		}
	}()

	go func() {
		log.Println("Sub server listening on", cfg.Ports.Sub, "(HTTPS:", cfg.HTTPS, ")")
		var err error
		if cfg.HTTPS {
			err = http.ListenAndServeTLS(cfg.Ports.Sub, cfg.TLS.Cert, cfg.TLS.Key, subMux)
		} else {
			err = http.ListenAndServe(cfg.Ports.Sub, subMux)
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("Sub server crashed:", err)
		}
	}()

	select {}
}
