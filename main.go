package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"wechat-saml-proxy/service"
	"wechat-saml-proxy/xsession"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("./index.html")
		if err != nil {
			fmt.Fprint(w, "内部错误")
			return
		}
		fmt.Fprint(w, string(b))
	})
	http.HandleFunc("/MP_verify_32iWga2EVle6QTQm.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "32iWga2EVle6QTQm")
	})
	http.HandleFunc("/api/callback", service.LoginHandler)
	http.HandleFunc("/api/saml", func(w http.ResponseWriter, r *http.Request) {
		session, _ := xsession.Store.Get(r, "user")
		fmt.Fprintf(w, "name is :%s", session.Values["openid4wechat"])
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
