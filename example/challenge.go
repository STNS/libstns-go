package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/STNS/libstns-go/libstns"
)

var stns *libstns.STNS

func challengeCode(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	code, err := stns.CreateUserChallengeCode(r.FormValue("user"))
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, string(code))
}

func verify(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	err := stns.VerifyWithUser(r.FormValue("user"), []byte(r.FormValue("code")), []byte(r.FormValue("signature")))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	code, err := stns.GetUserChallengeCode(r.FormValue("user"))
	if err != nil {
		panic(err)
	}
	if string(code) == r.FormValue("code") {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}
func main() {
	s, err := libstns.NewSTNS("https://stns.lolipop.io/v1/", nil)
	if err != nil {
		panic(err)
	}
	stns = s

	go func() {
		http.HandleFunc("/challenge", challengeCode)
		http.HandleFunc("/verify", verify)
		if err := http.ListenAndServe("127.0.0.1:18000", nil); err != nil {
			panic(err)
		}
	}()

	// sorry...
	time.Sleep(1 * time.Second)
	u := "http://127.0.0.1:18000/challenge?user=pyama"

	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	sig, err := stns.Sign(code)
	if err != nil {
		panic(err)
	}
	values := url.Values{}
	values.Set("user", "pyama")
	values.Set("signature", string(sig))
	values.Add("code", string(code))

	req, err := http.NewRequest(
		"POST",
		"http://127.0.0.1:18000/verify",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("verify ok")
	} else {
		fmt.Println("verify failed")
	}
}
