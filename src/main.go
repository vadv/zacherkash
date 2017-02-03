package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"transport"
)

var (
	BuildVersion = "unknown"
	config_path  = flag.String("config", "/etc/zacherkash.yaml", "path to config file")
	pid_file     = flag.String("pid", "", "path to pid file, if needed")
	version      = flag.Bool("version", false, "print version and exit")
)

type Config struct {
	Bind        string            `yaml:"bind"`
	LogFile     string            `yaml:"log_file"`
	BodyRewrite map[string]string `yaml:"body_rewrite"`
	Upstream    string            `yaml:"upstream"`
}

func main() {

	if !flag.Parsed() {
		flag.Parse()
	}

	if *version {
		fmt.Printf("version: %s\n", BuildVersion)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*config_path)
	if err != nil {
		log.Fatalf("[ERROR] Open config file: %s\n", err.Error())
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("[ERROR] Config error: %s\n", err.Error())
	}

	if err := transport.BuildBodyRewrites(config.BodyRewrite); err != nil {
		log.Fatalf("[ERROR] Compile rewite rules error: %s\n", err.Error())
	}

	os.MkdirAll(filepath.Dir(config.LogFile), 0755)
	fd, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("[ERROR] Open log file: %s\n", err.Error())
	}

	if *pid_file != "" {
		if err := ioutil.WriteFile(*pid_file, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
			log.Fatalf("[ERROR] Write pid file: %s\n", err.Error())
		}
	}

	// переписываем host
	if config.Upstream != "" {
		transport.Upstream = config.Upstream
	}

	log.SetOutput(fd)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = r.Host
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.Transport = &transport.BodyRewriter{http.DefaultTransport}
		proxy.ServeHTTP(w, r)
	})

	log.Printf("[INFO] Start listen at: %s\n", config.Bind)
	if err := http.ListenAndServe(config.Bind, nil); err != nil {
		log.Fatalf("[ERROR] Listen: %s\n", err.Error())
	}
}
