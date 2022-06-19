package controllers

import (
	nethttp "net/http"
	"strings"

	"github.com/dewep-online/solocms/internal/pkg/array"

	"github.com/dewep-online/goppy/middlewares"
	"golang.org/x/text/language"
)

func AllowDomainsMiddleware(allow []string) middlewares.Middleware {
	domains := array.StringsToMap(allow)

	return func(call func(nethttp.ResponseWriter, *nethttp.Request)) func(nethttp.ResponseWriter, *nethttp.Request) {
		return func(w nethttp.ResponseWriter, r *nethttp.Request) {
			if _, ok := domains[r.Host]; !ok {
				w.WriteHeader(nethttp.StatusForbidden)
				return
			}
			call(w, r)
		}
	}
}

func LangMiddleware(data []string) middlewares.Middleware {
	langs := array.StringsToMap(data)
	defaultLang := "en"
	if len(data) > 0 {
		defaultLang = data[0]
	}

	return func(call func(nethttp.ResponseWriter, *nethttp.Request)) func(nethttp.ResponseWriter, *nethttp.Request) {
		return func(w nethttp.ResponseWriter, r *nethttp.Request) {

			info := strings.SplitN(r.RequestURI, "/", 3)
			lang, uri := info[1], ""

			if len(info) == 3 {
				lang, uri = info[1], info[2]
			}

			if len(lang) != 2 {
				lang, uri = "", r.RequestURI[1:]
			}

			if _, ok := langs[lang]; !ok {
				lang = defaultLang
				tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
				if err == nil {
					for _, tag := range tags {
						base, _ := tag.Base()
						if _, ok = langs[base.String()]; ok {
							lang = base.String()
							break
						}
					}
				}

				nethttp.Redirect(w, r, "/"+lang+"/"+uri, nethttp.StatusPermanentRedirect)
				return
			}

			call(w, r)
		}
	}
}

func BasicAuthMiddleware(access map[string]string) middlewares.Middleware {
	authCall := func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		w.WriteHeader(nethttp.StatusUnauthorized)
	}

	return func(call func(nethttp.ResponseWriter, *nethttp.Request)) func(nethttp.ResponseWriter, *nethttp.Request) {
		return func(w nethttp.ResponseWriter, r *nethttp.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				authCall(w, r)
				return
			}

			pass, ok := access[username]
			if !ok || pass != password {
				authCall(w, r)
				return
			}

			call(w, r)
		}
	}
}
