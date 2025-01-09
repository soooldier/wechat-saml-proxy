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

	callback4cas = "https://idaas-saas-idp.eco.teems.com.cn/cidp/login/ai-b41f38383fcb411fb5f0ff0ec3166152"
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
		session, _ := xsession.Store.Get(r, "user")
		if _, ok := session.Values["openid"]; !ok {
			session.Values["redirect"] = "https://wechat-saml-proxy-v1-133482-9-1333979547.sh.run.tcloudbase." +
				"com/cas/login?service=" + r.URL.Query().Get("service")
			w.Header().Set("Location", fmt.Sprintf("https://open.weixin.qq."+
				"com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo"+
				"#wechat_redirect", appid, url.QueryEscape(callback)))
			w.WriteHeader(301)
			return
		}
		parsed, _ := url.Parse(r.URL.Query().Get("service"))
		ticket := uuid.NewString()
		memory.Set(ticket, []byte(session.Values["openid"].(string)))
		w.Header().Set("Location", fmt.Sprintf("%s?ticket=%s&state=%s", callback4cas, ticket, parsed.Query().Get("state")))
		w.WriteHeader(301)
		return

	})
	http.HandleFunc("/cas/check", func(w http.ResponseWriter, r *http.Request) {
		ticket := r.URL.Query().Get("ticket")
		openid, _ := memory.Get(ticket)
		defer func() {
			memory.Delete(ticket)
		}()
		openid4string := string(openid)
		if openid4string == "oNEbn6637Lh18k3ZAN7mkRGq-U2U" {
			openid4string = "yitttang"
		} else {
			openid4string = "unknown"
		}
		xml := "<cas:serviceResponse xmlns:cas=\"http://www.yale.edu/tp/cas\">\n   <cas:authenticationSuccess>\n " +
			"  <cas:user>" + openid4string + "</cas:user>\n   " +
			"  <cas:attributes>\n   " +
			"  <cas:user>" + openid4string + "</cas:user>\n   " +
			"  <cas:userSourceId></cas:userSourceId>\n  " +
			"  <cas:mail></cas:mail>\n" +
			"  <cas:userId>" + openid4string + "</cas:userId>\n" +
			"  </cas:attributes>\n " +
			"  <cas:proxyGrantingTicket>PGTIOU-84678-8a9d...</cas:proxyGrantingTicket>\n " +
			"  </cas:authenticationSuccess>\n </cas:serviceResponse>"
		fmt.Fprint(w, xml)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
