package gitserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"

	"github.com/pnocera/novaco/internal/config"
)

func GitServer(config *config.Config) error {
	cfg = config

	client, _ := api.NewClient(api.DefaultConfig())

	schema := "http"
	if cfg.SSLEnabled {
		schema = "https"
	}
	hcheck := fmt.Sprintf("%s://%s:%d/health", schema, cfg.Hostname, cfg.Port)

	client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: "novaco-gitserver",
		Tags: []string{"gitserver"},
		Port: cfg.Port,
		Check: &api.AgentServiceCheck{
			HTTP:     hcheck,
			Interval: "10s"},
	})

	host := fmt.Sprintf("%s:%d", config.Hostname, config.Port)
	hall := fmt.Sprintf(":%d", config.Port)
	log.Println("[INFO] Starting git http server at ", host)

	r := mux.NewRouter()
	attachHandler(r)

	if cfg.SSLEnabled {
		if err := http.ListenAndServeTLS(hall, cfg.CertFilePath, cfg.KeyFilePath, r); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(hall, r); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func attachHandler(r *mux.Router) {
	//git methods Handler
	r.HandleFunc(`/{user-name}/{repo-name}/info/refs`, basicAuthentication(infoserviceHandler)).Methods("GET")
	r.HandleFunc(`/{user-name}/{repo-name}/git-upload-pack`, basicAuthentication(uploadPackHandler2)).Methods("POST")
	r.HandleFunc(`/{user-name}/{repo-name}/git-receive-pack`, basicAuthentication(receivePackHandler2)).Methods("POST")

	//APIs handlers
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc(GetRepoCreateURL(), basicAuthentication(repoCreateHandler)).Methods("POST")
	r.HandleFunc(GetReposURL(), repoIndexHandler).Methods("GET")
	r.HandleFunc(GetRepoURL(), repoShowHandler).Methods("GET")
	r.HandleFunc(GetBranchesURL(), branchIndexHandler).Methods("GET")
	r.HandleFunc(GetBranchURL(), branchShowHandler).Methods("GET")
}
