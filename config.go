package main

import (
	"encoding/json"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/mchmarny/ws-collector/common"
	"log"
	"os"
)

var args = Config{}

func init() {

	loadConfig("./config/config.json", &args)
	trace := common.GetEnvVarAsBool(os.Getenv("TRACE"), false)
	if trace {
		args.Trace = trace
	}
	if args.Backend.Args == nil {
		args.Backend.Args = make(map[string]string)
	}

	// are we on CF?
	cf, err := cfenv.Current()
	if err == nil && cf != nil {
		if args.Trace {
			common.Printout(cf)
		}
		args.ID = cf.ID + "-" + string(cf.Index)
		args.Server.Port = cf.Port
		args.Server.Host = cf.Host

		// get baackend URI from CF
		cfSrvName := args.Backend.Args["service_name"]
		if len(cfSrvName) > 0 {
			srv, err := cf.Services.WithName(cfSrvName)
			if err == nil {
				args.Backend.URI = srv.Credentials["uri"]
			}
		}
	}
	if args.Trace {
		common.Printout(args)
	}

}

type ServerConfig struct {
	Root  string `json:"root,omitempty"`
	Host  string `json:"host,omitempty"`
	Port  int    `json:"port,omitempty"`
	Token string `json:"token,omitempty"`
}

type BackendConfig struct {
	Type string            `json:"type,omitempty"`
	URI  string            `json:"uri,omitempty"`
	Args map[string]string `json:"args,omitempty"`
}

type Config struct {
	ID      string        `json:"id,omitempty"`
	Trace   bool          `json:"trace,omitempty"`
	Server  ServerConfig  `json:"server,omitempty"`
	Backend BackendConfig `json:"backend,omitempty"`
}

func loadConfig(path string, c *Config) {
	log.Printf("loading config from file: %s ...", path)
	f, err := os.Open(path)
	if err != nil {
		log.Panicf("error while reading config file: %s - %v", path, err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		log.Panicf("error while parsing config file: %s - %v", path, err)
	}
}
