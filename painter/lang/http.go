package lang

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

func HttpHandler(loop *painter.Loop, p *Parser) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var in io.Reader = r.Body
		if r.Method == http.MethodGet {
			// Заменяем запятые на переносы строк
			cmd := strings.ReplaceAll(r.URL.Query().Get("cmd"), ",", "\n")
			in = strings.NewReader(cmd)
		}

		cmds, err := p.Parse(in)
		if err != nil {
			log.Printf("Bad script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, cmd := range cmds {
			loop.Post(cmd)
		}
		rw.WriteHeader(http.StatusOK)
	})
}
