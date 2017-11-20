package form

import (
	//iogos-replace
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cheikhshift/db"
	netform "github.com/cheikhshift/form"
	"github.com/cheikhshift/gos/core"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/fatih/color"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

var store = sessions.NewCookieStore([]byte("a very very very very secret key"))

type NoStruct struct {
	/* emptystruct */
}

func NetsessionGet(key string, s *sessions.Session) string {
	return s.Values[key].(string)
}

func UrlAtZ(url, base string) (isURL bool) {
	isURL = strings.Index(url, base) == 0
	return
}

func NetsessionDelete(s *sessions.Session) string {
	//keys := make([]string, len(s.Values))

	//i := 0
	for k := range s.Values {
		// keys[i] = k.(string)
		NetsessionRemove(k.(string), s)
		//i++
	}

	return ""
}

func NetsessionRemove(key string, s *sessions.Session) string {
	delete(s.Values, key)
	return ""
}
func NetsessionKey(key string, s *sessions.Session) bool {
	if _, ok := s.Values[key]; ok {
		//do something here
		return true
	}

	return false
}

func Netadd(x, v float64) float64 {
	return v + x
}

func Netsubs(x, v float64) float64 {
	return v - x
}

func Netmultiply(x, v float64) float64 {
	return v * x
}

func Netdivided(x, v float64) float64 {
	return v / x
}

func NetsessionGetInt(key string, s *sessions.Session) interface{} {
	return s.Values[key]
}

func NetsessionSet(key string, value string, s *sessions.Session) string {
	s.Values[key] = value
	return ""
}
func NetsessionSetInt(key string, value interface{}, s *sessions.Session) string {
	s.Values[key] = value
	return ""
}

func dbDummy() {
	smap := db.O{}
	smap["key"] = "set"
	log.Println(smap)
}

func Netimportcss(s string) string {
	return fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\" /> ", s)
}

func Netimportjs(s string) string {
	return fmt.Sprintf("<script type=\"text/javascript\" src=\"%s\" ></script> ", s)
}

func formval(s string, r *http.Request) string {
	return r.FormValue(s)
}

func renderTemplate(w http.ResponseWriter, p *Page) bool {
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path : web%s.tmpl reason : %s", p.R.URL.Path, n))

			DebugTemplate(w, p.R, fmt.Sprintf("web%s", p.R.URL.Path))
			w.WriteHeader(http.StatusInternalServerError)

			pag, err := loadPage("/your-500-page")

			if err != nil {
				log.Println(err.Error())
				return
			}

			if pag.isResource {
				w.Write(pag.Body)
			} else {
				pag.R = p.R
				pag.Session = p.Session
				renderTemplate(w, pag) ///your-500-page"

			}
		}
	}()

	t := template.New("PageWrapper")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(p.Body))
	outp := new(bytes.Buffer)
	err := t.Execute(outp, p)
	if err != nil {
		log.Println(err.Error())
		DebugTemplate(w, p.R, fmt.Sprintf("web%s", p.R.URL.Path))
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		pag, err := loadPage("/your-500-page")

		if err != nil {
			log.Println(err.Error())
			return false
		}
		pag.R = p.R
		pag.Session = p.Session
		p = nil
		if pag.isResource {
			w.Write(pag.Body)
		} else {
			renderTemplate(w, pag) // "/your-500-page"

		}
		return false
	}

	p.Session.Save(p.R, w)

	fmt.Fprintf(w, html.UnescapeString(outp.String()))
	p.Session = nil
	p.Body = nil
	p.R = nil
	p = nil
	return true

}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string, *sessions.Session)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var session *sessions.Session
		var er error
		if session, er = store.Get(r, "session-"); er != nil {
			session, _ = store.New(r, "session-")
		}
		if attmpt := apiAttempt(w, r, session); !attmpt {
			fn(w, r, "", session)
		} else {
			context.Clear(r)
		}

	}
}

func mResponse(v interface{}) string {
	data, _ := json.Marshal(&v)
	return string(data)
}
func apiAttempt(w http.ResponseWriter, r *http.Request, session *sessions.Session) (callmet bool) {
	var response string
	response = ""

	if strings.Contains(r.URL.Path, "") {

		if _, ok := session.Values["formtoken"]; !ok {
			session.Values["formtoken"] = core.NewLen(10)
			session.Save(r, w)
		}

		if r.ContentLength > 0 {

			if !netform.ValidateRequest(r, session.Values["formtoken"].(string)) || r.ContentLength > int64(netform.MaxSize*netform.MB) {

				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "text/xml")

				w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><Error>Invalid request sent</Error>"))
				return true
			}
		}

	}

	if isURL := (r.URL.Path == "/test/form-build" && r.Method == strings.ToUpper("GET")); !callmet && isURL {

		w.Header().Set("Content-Type", "text/html")
		SampleForm := SampleForm{Text: "Sample", Created: "2017-05-05", Emal: "sample", FieldF: "orange", Count: 500, Terms: true}
		w.Write([]byte(NetBuild(&SampleForm, "/target/url", "POST", "Update", session)))

		callmet = true
	}

	if isURL := (r.URL.Path == "/target/url" && r.Method == strings.ToUpper("POST")); !callmet && isURL {

		var sampleform SampleForm
		if err := netform.Form(r, &sampleform); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/xml")
			w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?><Error>%s Value :%v</Error>", err.Error(), sampleform)))
			return true
		}
		response = mResponse(sampleform)
		//http.ServeFile(w, r, netform.Path(sampleform.Photo))

		callmet = true
	}

	if callmet {
		session.Save(r, w)
		if response != "" {
			//Unmarshal json
			//w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(response))
		}
		return
	}
	return
}
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}
func DebugTemplate(w http.ResponseWriter, r *http.Request, tmpl string) {
	lastline := 0
	linestring := ""
	defer func() {
		if n := recover(); n != nil {
			log.Println()
			// log.Println(n)
			log.Println("Error on line :", lastline+1, ":"+strings.TrimSpace(linestring))
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()

	p, err := loadPage(r.URL.Path)
	filename := tmpl + ".tmpl"
	body, err := Asset(filename)
	session, er := store.Get(r, "session-")

	if er != nil {
		session, er = store.New(r, "session-")
	}
	p.Session = session
	p.R = r
	if err != nil {
		log.Print(err)

	} else {

		lines := strings.Split(string(body), "\n")
		// log.Println( lines )
		linebuffer := ""
		waitend := false
		open := 0
		for i, line := range lines {

			processd := false

			if strings.Contains(line, "{{with") || strings.Contains(line, "{{ with") || strings.Contains(line, "with}}") || strings.Contains(line, "with }}") || strings.Contains(line, "{{range") || strings.Contains(line, "{{ range") || strings.Contains(line, "range }}") || strings.Contains(line, "range}}") || strings.Contains(line, "{{if") || strings.Contains(line, "{{ if") || strings.Contains(line, "if }}") || strings.Contains(line, "if}}") || strings.Contains(line, "{{block") || strings.Contains(line, "{{ block") || strings.Contains(line, "block }}") || strings.Contains(line, "block}}") {
				linebuffer += line
				waitend = true

				endstr := ""
				processd = true
				if !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}")) {

					open++

				}
				for i := 0; i < open; i++ {
					endstr += "\n{{end}}"
				}
				//exec
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
				t, _ = t.Parse(ReadyTemplate(body))
				lastline = i
				linestring = line
				erro := t.Execute(outp, p)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}
			}

			if waitend && !processd && !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end")) {
				linebuffer += line

				endstr := ""
				for i := 0; i < open; i++ {
					endstr += "\n{{end}}"
				}
				//exec
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
				t, _ = t.Parse(ReadyTemplate(body))
				lastline = i
				linestring = line
				erro := t.Execute(outp, p)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}

			}

			if !waitend && !processd {
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
				t, _ = t.Parse(ReadyTemplate(body))
				lastline = i
				linestring = line
				erro := t.Execute(outp, p)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}
			}

			if !processd && (strings.Contains(line, "{{end") || strings.Contains(line, "{{ end")) {
				open--

				if open == 0 {
					waitend = false

				}
			}
		}

	}

}

func DebugTemplatePath(tmpl string, intrf interface{}) {
	lastline := 0
	linestring := ""
	defer func() {
		if n := recover(); n != nil {

			log.Println("Error on line :", lastline+1, ":"+strings.TrimSpace(linestring))
			log.Println(n)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()

	filename := tmpl
	body, err := Asset(filename)

	if err != nil {
		log.Print(err)

	} else {

		lines := strings.Split(string(body), "\n")
		// log.Println( lines )
		linebuffer := ""
		waitend := false
		open := 0
		for i, line := range lines {

			processd := false

			if strings.Contains(line, "{{with") || strings.Contains(line, "{{ with") || strings.Contains(line, "with}}") || strings.Contains(line, "with }}") || strings.Contains(line, "{{range") || strings.Contains(line, "{{ range") || strings.Contains(line, "range }}") || strings.Contains(line, "range}}") || strings.Contains(line, "{{if") || strings.Contains(line, "{{ if") || strings.Contains(line, "if }}") || strings.Contains(line, "if}}") || strings.Contains(line, "{{block") || strings.Contains(line, "{{ block") || strings.Contains(line, "block }}") || strings.Contains(line, "block}}") {
				linebuffer += line
				waitend = true

				endstr := ""
				if !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}")) {

					open++

				}

				for i := 0; i < open; i++ {
					endstr += "\n{{end}}"
				}
				//exec

				processd = true
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
				t, _ = t.Parse(ReadyTemplate([]byte(fmt.Sprintf("%s%s", linebuffer, endstr))))
				lastline = i
				linestring = line
				erro := t.Execute(outp, intrf)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}
			}

			if waitend && !processd && !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}")) {
				linebuffer += line

				endstr := ""
				for i := 0; i < open; i++ {
					endstr += "\n{{end}}"
				}
				//exec
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
				t, _ = t.Parse(ReadyTemplate([]byte(fmt.Sprintf("%s%s", linebuffer, endstr))))
				lastline = i
				linestring = line
				erro := t.Execute(outp, intrf)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}

			}

			if !waitend && !processd {
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
				t, _ = t.Parse(ReadyTemplate([]byte(fmt.Sprintf("%s%s", linebuffer))))
				lastline = i
				linestring = line
				erro := t.Execute(outp, intrf)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}
			}

			if !processd && (strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}")) {
				open--

				if open == 0 {
					waitend = false

				}
			}
		}

	}

}
func Handler(w http.ResponseWriter, r *http.Request, contxt string, session *sessions.Session) {
	var p *Page
	p, err := loadPage(r.URL.Path)

	if err != nil {
		log.Println(err.Error())

		w.WriteHeader(http.StatusNotFound)

		pag, err := loadPage("/your-404-page")

		if err != nil {
			log.Println(err.Error())
			//context.Clear(r)
			return
		}
		pag.R = r
		pag.Session = session
		if p != nil {
			p.Session = nil
			p.Body = nil
			p.R = nil
			p = nil
		}
		if pag.isResource {
			w.Write(pag.Body)
		} else {
			renderTemplate(w, pag) //"/your-500-page"
		}
		context.Clear(r)
		return
	}

	if !p.isResource {
		w.Header().Set("Content-Type", "text/html")
		p.Session = session
		p.R = r
		renderTemplate(w, p) //fmt.Sprintf("web%s", r.URL.Path)

		// log.Println(w)
	} else {
		w.Header().Set("Cache-Control", "public")
		if strings.Contains(r.URL.Path, ".css") {
			w.Header().Add("Content-Type", "text/css")
		} else if strings.Contains(r.URL.Path, ".js") {
			w.Header().Add("Content-Type", "application/javascript")
		} else {
			w.Header().Add("Content-Type", http.DetectContentType(p.Body))
		}

		w.Write(p.Body)
	}

	p.Session = nil
	p.Body = nil
	p.R = nil
	p = nil
	context.Clear(r)
	return
}

func loadPage(title string) (*Page, error) {

	if roottitle := (title == "/"); roottitle {
		webbase := "web/"
		fname := fmt.Sprintf("%s%s", webbase, "index.html")
		body, err := Asset(fname)
		if err != nil {
			fname = fmt.Sprintf("%s%s", webbase, "index.tmpl")
			body, err = Asset(fname)
			if err != nil {
				return nil, err
			}
			return &Page{Body: body, isResource: false}, nil
		}

		return &Page{Body: body, isResource: true}, nil

	}

	filename := fmt.Sprintf("web%s.tmpl", title)

	if body, err := Asset(filename); err != nil {
		filename = fmt.Sprintf("web%s.html", title)

		if body, err = Asset(filename); err != nil {
			filename = fmt.Sprintf("web%s", title)

			if body, err = Asset(filename); err != nil {
				return nil, err
			} else {
				if strings.Contains(title, ".tmpl") {
					return nil, nil
				}
				return &Page{Body: body, isResource: true}, nil
			}
		} else {
			return &Page{Body: body, isResource: true}, nil
		}
	} else {
		return &Page{Body: body, isResource: false}, nil
	}

}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
func equalz(args ...interface{}) bool {
	if args[0] == args[1] {
		return true
	}
	return false
}
func nequalz(args ...interface{}) bool {
	if args[0] != args[1] {
		return true
	}
	return false
}

func netlt(x, v float64) bool {
	if x < v {
		return true
	}
	return false
}
func netgt(x, v float64) bool {
	if x > v {
		return true
	}
	return false
}
func netlte(x, v float64) bool {
	if x <= v {
		return true
	}
	return false
}

func GetLine(fname string, match string) int {
	intx := 0
	file, err := os.Open(fname)
	if err != nil {
		color.Red("Could not find a source file")
		return -1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		intx = intx + 1
		if strings.Contains(scanner.Text(), match) {

			return intx
		}

	}

	return -1
}
func netgte(x, v float64) bool {
	if x >= v {
		return true
	}
	return false
}

type Page struct {
	Title      string
	Body       []byte
	isResource bool
	R          *http.Request
	Session    *sessions.Session
}

func ReadyTemplate(body []byte) string {
	return strings.Replace(strings.Replace(strings.Replace(string(body), "/{", "\"{", -1), "}/", "}\"", -1), "`", "\"", -1)
}

type SampleForm struct {
	TestField string           `title:"Hi world!",valid:"unique",placeholder:"Testfield prompt"`
	Count     int              `placeholder:"Count"`
	Name      netform.Password `valid:"required",title:"Input title"`
	FieldTwo  netform.Radio    `title:"Enter Email",valid:"email,unique,required",select:"blue,orange,red,green"`
	FieldF    netform.Select   `placeholder:"Prompt?",valid:"email,unique,required",select:"blue,orange,red,green"`
	Created   netform.Date
	Text      netform.Paragraph `title:"Enter a description."`
	Photo     netform.File      `file:"image/*"`
	Emal      netform.Email
	Terms     bool `title:"Accept terms of use."`
}

func NetcastSampleForm(args ...interface{}) *SampleForm {

	s := SampleForm{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructSampleForm() *SampleForm { return &SampleForm{} }

type fInput struct {
	Name, Type, Placeholder, Title, Class, Value string
	Choices                                      []string
	Required                                     bool
}

func NetcastfInput(args ...interface{}) *fInput {

	s := fInput{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructfInput() *fInput { return &fInput{} }

type afInput struct {
	Name, Type, Placeholder, Title, Class, Value string
	Choices                                      []string
	Required                                     bool
	ModelName                                    string
}

func NetcastafInput(args ...interface{}) *afInput {

	s := afInput{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructafInput() *afInput { return &afInput{} }

type fForm struct {
	Target, Method, Token, ButtonClass, CTA string
	Input                                   []fInput
}

func NetcastfForm(args ...interface{}) *fForm {

	s := fForm{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructfForm() *fForm { return &fForm{} }

type aForm struct {
	Target, Method, Token, ButtonClass, CTA string
	Input                                   []afInput
	ModelName                               string
}

func NetcastaForm(args ...interface{}) *aForm {

	s := aForm{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructaForm() *aForm { return &aForm{} }

//
func NetaC(args ...interface{}) string {

	return "}}"

}

//
func NetaO(args ...interface{}) string {

	return "{{"

}

//
func NetIsIn(args ...interface{}) bool {

	return strings.Contains(args[0].(string), args[1].(string))

}

//
func NetHasBody(args ...interface{}) bool {

	return args[0].(*http.Request).ContentLength > 0

}

//
func NetForm(args ...interface{}) string {

	err := netform.Form(args[0].(*http.Request), args[1])
	return err.Error()

}

//
func NetTokenizeForm(args ...interface{}) (form fForm) {

	v := reflect.ValueOf(args[0]).Elem()
	//t := reflect.TypeOf(item)
	bso := netform.ToBson(mResponse(args[0]))
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		fieldtype := strings.ToLower(field.Type.String())
		requird := strings.Contains(string(field.Tag), "required")
		opts := netform.GetSel(string(field.Tag))
		title := field.Tag.Get("title")
		placehlder := netform.GetPl(string(field.Tag))

		if strings.Contains(fieldtype, "bool") {
			form.Input = append(form.Input, fInput{field.Name, "checkbox", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
		} else if strings.Contains(fieldtype, "int") || strings.Contains(fieldtype, "float") {
			form.Input = append(form.Input, fInput{field.Name, "number", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
		} else if strings.Contains(fieldtype, "string") {
			form.Input = append(form.Input, fInput{field.Name, "text", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
		} else if strings.Contains(fieldtype, "email") {
			form.Input = append(form.Input, fInput{field.Name, "email", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
		} else if strings.Contains(fieldtype, "password") {
			form.Input = append(form.Input, fInput{field.Name, "password", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
		} else if strings.Contains(fieldtype, "select") {

			form.Input = append(form.Input, fInput{field.Name, "select", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird})

		} else if strings.Contains(fieldtype, "radio") {

			form.Input = append(form.Input, fInput{field.Name, "radio", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird})

		} else if strings.Contains(fieldtype, "selectmult") {

			form.Input = append(form.Input, fInput{field.Name, "selectmult", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird})

		} else if strings.Contains(fieldtype, "date") {
			form.Input = append(form.Input, fInput{field.Name, "date", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
		} else if strings.Contains(fieldtype, "file") {
			placehlder = field.Tag.Get("file")
			form.Input = append(form.Input, fInput{field.Name, "file", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
		} else if strings.Contains(fieldtype, "paragraph") {
			form.Input = append(form.Input, fInput{field.Name, "textarea", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
		} else {
			form.Input = append(form.Input, fInput{field.Name, "invalid", placehlder, title, netform.InputClass, "", nil, requird})
		}
	}

	return

}

//
func NetTokenizeFormAng(args ...interface{}) (form aForm) {

	v := reflect.ValueOf(args[0]).Elem()
	//t := reflect.TypeOf(item)
	modelClass := args[1].(string)
	bso := netform.ToBson(mResponse(args[0]))
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		fieldtype := strings.ToLower(field.Type.String())
		requird := strings.Contains(string(field.Tag), "required")
		opts := netform.GetSel(string(field.Tag))
		title := field.Tag.Get("title")
		placehlder := netform.GetPl(string(field.Tag))

		if strings.Contains(fieldtype, "bool") {
			form.Input = append(form.Input, afInput{field.Name, "checkbox", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
		} else if strings.Contains(fieldtype, "int") || strings.Contains(fieldtype, "float") {
			form.Input = append(form.Input, afInput{field.Name, "number", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
		} else if strings.Contains(fieldtype, "string") {
			form.Input = append(form.Input, afInput{field.Name, "text", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
		} else if strings.Contains(fieldtype, "email") {
			form.Input = append(form.Input, afInput{field.Name, "email", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
		} else if strings.Contains(fieldtype, "password") {
			form.Input = append(form.Input, afInput{field.Name, "password", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
		} else if strings.Contains(fieldtype, "select") {

			form.Input = append(form.Input, afInput{field.Name, "select", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird, modelClass})

		} else if strings.Contains(fieldtype, "radio") {

			form.Input = append(form.Input, afInput{field.Name, "radio", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird, modelClass})

		} else if strings.Contains(fieldtype, "selectmult") {

			form.Input = append(form.Input, afInput{field.Name, "selectmult", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird, modelClass})

		} else if strings.Contains(fieldtype, "date") {
			form.Input = append(form.Input, afInput{field.Name, "date", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
		} else if strings.Contains(fieldtype, "file") {
			placehlder = field.Tag.Get("file")
			form.Input = append(form.Input, afInput{field.Name, "file", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
		} else if strings.Contains(fieldtype, "paragraph") {
			form.Input = append(form.Input, afInput{field.Name, "textarea", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
		} else {
			form.Input = append(form.Input, afInput{field.Name, "invalid", placehlder, title, netform.InputClass, "", nil, requird, modelClass})
		}
	}

	return

}

// Return HTML form
// Argument 1 : Interface{} - Interface to build form with. Submit a variable with data to prepopulate form.
// Argument 2 : String - Target URL to submit form to.
// Argument 3 : String - Method of form submission (GET,POST,PUT etc...).
// Argument 4 : String - Call to action of form button.
// Argument 5 : *sessions.Session Current user session. Must be passed to ensure secure communication.
func NetBuild(args ...interface{}) string {

	if len(args) < 5 {
		return "<h1>No enough arguments to build form.</h1>"
	}
	form := NetTokenizeForm(args[0])
	target := args[1].(string)
	session := args[4].(*sessions.Session)

	form.Token = netform.GenerateToken(target, session.Values["formtoken"].(string))
	form.Method = args[2].(string)
	form.Target = target
	form.ButtonClass = netform.ButtonClass
	form.CTA = args[3].(string)
	return btForm(form)

}

// Return string of Angular form
// Argument 1 : Interface{} - Interface to build form with. Submit a variable with data to prepopulate form.
// Argument 2 : String - Target URL to submit form to. This is used to generate a token only valid for the specified target URL path.
// Argument 3 : String - JS Function to use with form's submit ng-click .
// Argument 4 : Call to action of form button.
// Argument 5 : String - variable name to be used as a local scope object to hold form data.
// Argument 5 : *sessions.Session Current user session. Must be passed to ensure secure communication.
func NetAngularForm(args ...interface{}) string {

	if len(args) < 6 {
		return "<h1>No enough arguments to build form.</h1>"
	}
	modelclass := args[4].(string)
	form := NetTokenizeFormAng(args[0], modelclass)
	target := args[1].(string)
	session := args[5].(*sessions.Session)
	form.ModelName = modelclass
	form.Token = netform.GenerateToken(target, session.Values["formtoken"].(string))
	form.Method = args[2].(string)
	form.Target = target
	form.ButtonClass = netform.ButtonClass
	form.CTA = args[3].(string)
	return batForm(form)

}

// Return  string of form token.
// Argument 0 : String of URI your form is submitting to
// Argument 1 : *sessions.Session (github.com/gorilla/sessions)
func NetGenerateToken(args ...interface{}) string {

	session := args[1].(*sessions.Session)
	return netform.GenerateToken(args[0].(string), session.Values["formtoken"].(string))

}

func NettInput(args ...interface{}) string {

	var d fInput
	filename := "tmpl/input.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (tInput) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = fInput{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("tInput")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func btInput(d fInput) string {
	return NetbtInput(d)
}

//
func NetbtInput(d fInput) string {

	filename := "tmpl/input.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("tInput")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (tInput) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func NetctInput(args ...interface{}) (d fInput) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = fInput{}
	}
	return
}

func ctInput(args ...interface{}) (d fInput) {
	if len(args) > 0 {
		d = NetctInput(args[0])
	} else {
		d = NetctInput()
	}
	return
}

func NetatInput(args ...interface{}) string {

	var d afInput
	filename := "tmpl/ainput.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (atInput) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = afInput{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("atInput")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func batInput(d afInput) string {
	return NetbatInput(d)
}

//
func NetbatInput(d afInput) string {

	filename := "tmpl/ainput.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("atInput")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (atInput) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func NetcatInput(args ...interface{}) (d afInput) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = afInput{}
	}
	return
}

func catInput(args ...interface{}) (d afInput) {
	if len(args) > 0 {
		d = NetcatInput(args[0])
	} else {
		d = NetcatInput()
	}
	return
}

func NettForm(args ...interface{}) string {

	var d fForm
	filename := "tmpl/form.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (tForm) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = fForm{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("tForm")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func btForm(d fForm) string {
	return NetbtForm(d)
}

//
func NetbtForm(d fForm) string {

	filename := "tmpl/form.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("tForm")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (tForm) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func NetctForm(args ...interface{}) (d fForm) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = fForm{}
	}
	return
}

func ctForm(args ...interface{}) (d fForm) {
	if len(args) > 0 {
		d = NetctForm(args[0])
	} else {
		d = NetctForm()
	}
	return
}

func NetatForm(args ...interface{}) string {

	var d aForm
	filename := "tmpl/aform.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (atForm) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = aForm{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("atForm")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func batForm(d aForm) string {
	return NetbatForm(d)
}

//
func NetbatForm(d aForm) string {

	filename := "tmpl/aform.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("atForm")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (atForm) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func NetcatForm(args ...interface{}) (d aForm) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = aForm{}
	}
	return
}

func catForm(args ...interface{}) (d aForm) {
	if len(args) > 0 {
		d = NetcatForm(args[0])
	} else {
		d = NetcatForm()
	}
	return
}

func dummy_timer() {
	dg := time.Second * 5
	log.Println(dg)
}
func FileServer() http.Handler {
	return http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "web"})
}
