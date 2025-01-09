package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/google/uuid"
	"wechat-saml-proxy/service"
	"wechat-saml-proxy/xsession"
)

var (
	callback = "https://wechat-saml-proxy-v1-133482-9-1333979547.sh.run.tcloudbase.com/api/callback"
	appid    = "wxbb7b02e8aaffb2e4"

	callback4cas = ""
)

func main() {
	memory, _ := bigcache.New(context.TODO(), bigcache.DefaultConfig(1*time.Minute))
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
	http.HandleFunc("/cas/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Query())
		return
		session, _ := xsession.Store.Get(r, "user")
		if _, ok := session.Values["openid"]; !ok {
			session.Values["redirect"] = "https://wechat-saml-proxy-v1-133482-9-1333979547.sh.run.tcloudbase." +
				"com/cas/login"
			w.Header().Set("Location", fmt.Sprintf("https://open.weixin.qq."+
				"com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo"+
				"#wechat_redirect", appid, url.QueryEscape(callback)))
			w.WriteHeader(301)
			return
		}
		ticket := uuid.NewString()
		memory.Set(ticket, session.Values["openid"].([]byte))
		w.Header().Set("Location", fmt.Sprintf("%s?ticket=%s&state=%s", callback4cas, ticket))
		w.WriteHeader(301)
		return

	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
