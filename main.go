package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/ghodss/yaml"
	validator "gopkg.in/go-playground/validator.v9"
)

type Trigger struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Token string `json:"token"`
}

type Deployment struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Token    string `json:"token"`
	UseValue bool   `json:"useValue"`
}

type Config struct {
	Deployments []Deployment `json:"deployments"`
}

func LoadConfig() (config Config, err error) {
	configPath := os.Getenv("CONFIG")
	if len(configPath) < 1 {
		configPath = "./config.yaml"
	}
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return
	}
	validate := validator.New()
	err = validate.Struct(config)

	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		config, err := LoadConfig()
		if err != nil {
			log.Println(err)
			if strings.Contains(err.Error(), "no such file or directory") {
				log.Println("Failed to load config file")
				os.Exit(1)
			}
		}

		var t Trigger
		err = json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		found := false
		var deployment Deployment
		for _, deploy := range config.Deployments {
			if deploy.Name == t.Name && deploy.Token == t.Token {
				found = true
				deployment = deploy
			}
		}
		if !found {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !deployment.UseValue {
			t.Value = ""
		}

		output, err := exec.Command(deployment.Command, t.Value).Output()
		if err != nil {
			log.Println(err)
			http.Error(w, "internal-server-error", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(string(output)))

	default:
		fmt.Fprintf(w, "Tendang!")
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8000"
	}
	fmt.Printf("Starting tendang at port " + port + "\n")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
