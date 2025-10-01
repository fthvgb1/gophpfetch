package main

import (
	"github.com/fthvgb1/wp-go/helper"
	"io"
	"net/http"
	"os"
	"path"
)

func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Limit file size to 10MB. This line saves you from those accidental 100MB uploads!
		r.ParseMultipartForm(10 << 20)
		if r.MultipartForm.File == nil {
			return
		}
		for key := range r.MultipartForm.File {
			// Retrieve the file from form data
			file, _, err := r.FormFile(key)
			if err != nil {
				http.Error(w, "Error retrieving the file", http.StatusBadRequest)
				return
			}
			defer file.Close()

			err = helper.IsDirExistAndMkdir(path.Dir(key), 0775)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			f, err := os.Create(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

	})
	http.ListenAndServe(":17778", nil)

}
