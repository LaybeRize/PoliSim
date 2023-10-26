package serving

import (
	"PoliSim/html/composition"
	"encoding/json"
	"net/http"
)

func InstallMarkdown() {
	composition.PatchHTMXFunctions[composition.MarkdownFormPage] = GetMarkdownItemFromForm
	composition.PatchHTMXFunctions[composition.MarkdownJsonPage] = GetMarkdownItemFromJson
}

type Content struct {
	Content string `input:"content" json:"content"`
}

func GetMarkdownItemFromForm(w http.ResponseWriter, r *http.Request) {
	content := &Content{}
	err := extractFormValuesForFields(content, r, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = composition.GetUpdatePreviewElement(content.Content).Render(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetMarkdownItemFromJson(w http.ResponseWriter, r *http.Request) {
	content := Content{}
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = composition.GetUpdatePreviewElement(content.Content).Render(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
