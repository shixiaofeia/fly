package main

import (
	"fmt"
	"github.com/kataras/hcaptcha"
	"html/template"
	"net/http"
)

// Get the following values from: https://dashboard.hcaptcha.com
// Also, check: https://docs.hcaptcha.com/#localdev to test on local environment.
var (
	siteKey   = "-"
	secretKey = "-"
)

var (
	client       = hcaptcha.New(secretKey) /* See `Client.FailureHandler` too. */
	registerForm = template.Must(template.ParseFiles("./pkg/example/hcaptcha/index.html"))
)

func main() {
	http.HandleFunc("/", renderForm)
	http.HandleFunc("/page", client.HandlerFunc(page) /* See `Client.SiteVerify` to get rid of a wrapper if necessary */)

	fmt.Printf("SiteKey = %s\tSecretKey = %s\nListening on: http://yourdomain.com\n",
		siteKey, secretKey)

	http.ListenAndServe(":8000", nil)
}

func page(w http.ResponseWriter, r *http.Request) {
	hcaptchaResp, ok := hcaptcha.Get(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Are you a bot?")
		return
	}

	fmt.Fprintf(w, "Page is inspected by a Human.\nResponse value is: %#+v", hcaptchaResp)
}

func renderForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	registerForm.Execute(w, map[string]string{
		"SiteKey": siteKey,
	})
}
