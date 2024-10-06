package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/swaggo/swag"
)

func (h *Handler) Swagger(w http.ResponseWriter, r *http.Request) {
	d, err := swag.ReadDoc()
	if err != nil {
		h.Log.Error(r.Context(), "swagger: read doc", err)
		return
	}

	data := make(map[string]interface{})
	data["host"] = r.Host
	if err := json.NewDecoder(strings.NewReader(d)).Decode(&data); err != nil {
		h.Log.Error(r.Context(), "swagger: decode doc", err)
		return
	}

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		h.Log.Error(r.Context(), "swagger: encode data", err)
		return
	}
}
