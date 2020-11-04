package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	b64 "encoding/base64"

	"github.com/ghodss/yaml"
	"github.com/sethvargo/go-password/password"
	"github.com/urfave/cli"
	"golang.org/x/crypto/pbkdf2"
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
	Hash     string `json:"hash"`
	UseValue bool   `json:"useValue"`
}

type Config struct {
	Deployments []Deployment `json:"deployments"`
	Salt        string       `json:"salt"`
}

var config *Config

func LoadConfig() (config *Config, err error) {
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

	return config, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var err error
		config, err = LoadConfig()
		if err != nil {
			log.Println(err)
			if strings.Contains(err.Error(), "no such file or directory") {
				log.Println("Failed to load config file")
				http.Error(w, "internal-server-error", http.StatusInternalServerError)
				return
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
		dk := pbkdf2.Key([]byte(t.Token), []byte(os.Getenv("SALT")), 4096, 16, sha1.New)
		sEnc := b64.StdEncoding.EncodeToString(dk)
		for _, deploy := range config.Deployments {
			if deploy.Name == t.Name && sEnc == deploy.Hash {
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

		// Reject any value that has characters other than: alphanumeric, -, and _
		re := regexp.MustCompile("^[a-zA-Z0-9_-]*$")
		if !re.MatchString(t.Value) {
			http.Error(w, "bad-request", http.StatusBadRequest)
			return
		}

		if len(t.Value) > 0 {
			deployment.Command += " " + t.Value
		}
		log.Println(t.Name, t.Value)
		log.Println(deployment.Command)
		output, err := exec.Command("bash", "-c", deployment.Command).Output()
		if err != nil {
			log.Println(err)
			http.Error(w, "internal-server-error", http.StatusInternalServerError)
			return
		}
		log.Println(string(output))
		w.Write([]byte("OK"))

	default:
		fmt.Fprintf(w, "Tendang!")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "tendang"
	app.Author = "BlankOn Developer"
	app.Email = "blankon-dev@googlegroups.com"

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Serving tendang proxy",
			Action: func(c *cli.Context) (err error) {
				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					handler(w, r)
				})
				port := os.Getenv("PORT")
				if len(port) < 1 {
					port = "8000"
				}
				fmt.Printf("Starting tendang at port " + port + "\n")
				if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:  "gen",
			Usage: "Generate key pair",
			Action: func(c *cli.Context) (err error) {
				token, err := password.Generate(64, 10, 0, false, true)
				if err != nil {
					return
				}
				dk := pbkdf2.Key([]byte(token), []byte(os.Getenv("SALT")), 4096, 16, sha1.New)
				sEnc := b64.StdEncoding.EncodeToString(dk)
				log.Println("Token: ", token)
				log.Println("Hash: ", sEnc)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
