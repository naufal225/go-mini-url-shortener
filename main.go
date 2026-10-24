package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strings"
)

type ShortenerRequest struct {
	Url string
}

type ShortenerResponse struct {
	Code string
}

type GetLinkResponse struct {
	Code string
	Link string
}

var StoredValue map[string]string = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwzyx1234567890"

func main() {
	http.HandleFunc("POST /shorten", HandlerShortener)
	http.HandleFunc("GET /{code}", HandlerGetLinkRedirect)
	http.HandleFunc("GET /info/{code}", HandlerGetInfo)

	log.Println("Server is listening on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error listen server")
	}

}

func generateCode() string {
	var b = make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func HandlerShortener(w http.ResponseWriter, r *http.Request) {
	// ambil request masuk
	decoder := json.NewDecoder(r.Body)
	// ambil link
	var req ShortenerRequest

	if err := decoder.Decode(&req); err != nil {
		fmt.Println(err)
		return
	}

	var link = req.Url
	// cek apakah link ada prefix http:// atau https://
	if !(strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")) {
		var message map[string]string = map[string]string{
			"message": "link harus diawali 'http://' atau 'https://'",
		}

		jsonInBytes, err := json.Marshal(&message)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonInBytes)
		return
	} else {
		// pecah http:// https:// dan link jadi masing masing
		_, domain, _ := strings.Cut(link, "//")

		// kalau prefix ada tapi link kosong, bad request
		if domain == "" {
			var message map[string]string = map[string]string{
				"message": "domain url tidak ada",
			}

			jsonInBytes, err := json.Marshal(&message)
			if err != nil {
				fmt.Println(err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonInBytes)
			return
		}
	}

	// cek apakah link udah ada
	var ada bool = false
	for code, val := range StoredValue {
		if val == link {
			// kalau ada balikkin itu aja
			var res ShortenerResponse = ShortenerResponse{Code: code}

			jsonInBytes, err := json.Marshal(&res)
			if err != nil {
				fmt.Println(err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(jsonInBytes)
			ada = true
			break
		}
	}

	if ada {
		return
	}

	// ubah jadi pasangan kode dan nilai, simpan ke StoredValue
	var code = generateCode()
	StoredValue[code] = link

	// return response
	var res ShortenerResponse = ShortenerResponse{Code: code}

	jsonInBytes, err := json.Marshal(&res)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
}

func HandlerGetInfo(w http.ResponseWriter, r *http.Request) {
	// ambil request
	var codeReq string = r.PathValue("code")
	// cari di StoredValue apakah ada
	var ada bool = false
	for code, val := range StoredValue {
		if code == codeReq {
			// kalau ada balikkin itu
			var res GetLinkResponse = GetLinkResponse{Code: code, Link: val}

			jsonInBytes, err := json.Marshal(&res)
			if err != nil {
				fmt.Println(err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusFound)
			w.Write(jsonInBytes)
			ada = true
			break
		}
	}
	// kalau ga ada return 404
	if !ada {
		var message map[string]string = map[string]string{
			"message": "link tidak ditermukan",
		}

		jsonInBytes, err := json.Marshal(&message)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonInBytes)
	}
}

func HandlerGetLinkRedirect(w http.ResponseWriter, r *http.Request) {
	// ambil request
	var codeReq string = r.PathValue("code")
	// cari di StoredValue apakah ada
	var ada bool = false
	var foundedLink string = ""
	for code, link := range StoredValue {
		if code == codeReq {

			ada = true
			foundedLink = link
			break
		}
	}

	if ada {
		http.Redirect(w, r, foundedLink, http.StatusMovedPermanently)
	}

	// kalau ga ada return 404
	if !ada {
		var message map[string]string = map[string]string{
			"message": "link tidak ditermukan",
		}

		jsonInBytes, err := json.Marshal(&message)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonInBytes)
	}
}
