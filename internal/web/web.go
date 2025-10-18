//go:build !noui
// +build !noui

package web

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

//go:embed build/*
var webUI embed.FS

type webFS struct {
	Fs http.FileSystem
}

func (fs *webFS) Open(name string) (http.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return fs.Fs.Open("index.html")
	}
	return f, err
}

// FrontEndHandler Provides Handler to Serve Frontend
var FrontEndHandler http.Handler

func init() {
	contentStatic, _ := fs.Sub(fs.FS(webUI), "build")
	FrontEndHandler = withEnvVars(http.FileServer(&webFS{Fs: http.FS(contentStatic)}))
}

// envs you want to expose (keep it a build-time list or read it in)
var envKeys = []string{"NO_AUTH"}

// withEnvVars returns a handler that injects window.ENV = {…}
func withEnvVars(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			next.ServeHTTP(w, r)
			return
		}

		// build JS object literal  {KEY:"value",…}
		buf := bytes.NewBufferString("window.ENV={")
		for i, k := range envKeys {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "%s:%q", k, os.Getenv(k))
		}
		buf.WriteString("};")

		// read embedded index.html
		data, _ := webUI.ReadFile("build/index.html")
		html := strings.Replace(string(data), "</head>",
			fmt.Sprintf(`<script>%s</script></head>`, buf.String()), 1)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
}
