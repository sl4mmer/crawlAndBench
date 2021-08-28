package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sl4mmer/crawlAndBench/pkg/common"
)

type Service struct {
	Queerer common.Queerer
}

func (s *Service) Handle(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if len(search) == 0 {
		writeErrStr(w, "empty request")
		return
	}
	result, err := s.Queerer.Query(r.Context(), search)
	if err != nil {
		writeErrStr(w, err.Error())
		return
	}
	response, err := json.Marshal(result)
	if err != nil {
		writeErrStr(w, err.Error())
		return
	}
	w.WriteHeader(200)
	w.Write(response)
}

func writeErrStr(writer http.ResponseWriter, err string) {
	resp := respErr{Error: err}
	r, _ := json.Marshal(resp)
	writer.WriteHeader(500)
	writer.Write(r)
}
