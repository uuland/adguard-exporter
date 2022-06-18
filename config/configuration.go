package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
	log "github.com/sirupsen/logrus"
)

// Config is the exporter CLI configuration.
type Config struct {
	LogLevel         string        `config:"loglevel"`
	AdguardProtocol  string        `config:"adguard_protocol"`
	AdguardHostname  []string      `config:"adguard_hostname"`
	AdguardUsername  string        `config:"adguard_username"`
	AdguardPassword  string        `config:"adguard_password"`
	AdguardPort      string        `config:"adguard_port"`
	ServerPort       string        `config:"server_port"`
	Interval         time.Duration `config:"interval"`
	LogLimit         string        `config:"log_limit"`
	RDnsEnabled      bool          `config:"rdns_enabled"`
	PasswordFromFile bool          `config:"password_from_file"`
}

func getDefaultConfig() *Config {
	return &Config{
		LogLevel:         "debug",
		AdguardProtocol:  "http",
		AdguardHostname:  []string{"127.0.0.1"},
		AdguardUsername:  "",
		AdguardPassword:  "",
		AdguardPort:      "80",
		ServerPort:       "9617",
		Interval:         10 * time.Second,
		LogLimit:         "1000",
		RDnsEnabled:      true,
		PasswordFromFile: false,
	}
}

// Load method loads the configuration by using both flag or environment variables.
func Load() *Config {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	loader := confita.NewLoader(loaders...)

	cfg := getDefaultConfig()
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		log.Errorf("Could not load the configuration...")
		os.Exit(1)
	}

	//Set the adguard port based on the input configuration
	if cfg.AdguardPort == "" {
		if cfg.AdguardProtocol == "http" {
			cfg.AdguardPort = "80"
		} else if cfg.AdguardProtocol == "https" {
			cfg.AdguardPort = "443"
		} else {
			log.Errorf("protocol %s is invalid. Must be http or https.", cfg.AdguardProtocol)
			os.Exit(1)
		}
	}

	//Set the adguard password based on the input configuration
	if cfg.PasswordFromFile {
		secret, err := ioutil.ReadFile(cfg.AdguardPassword)
		if err != nil {
			log.Errorf("unable to read AdguardPassword from %s due to %s", cfg.AdguardPassword, err)
			os.Exit(1)
		}

		cfg.AdguardPassword = string(secret)
	}

	cfg.show()

	return cfg
}

func (c Config) show() {
	val := reflect.ValueOf(&c).Elem()
	log.Println("---------------------------------------")
	log.Println("- AdGuard Home exporter configuration -")
	log.Println("---------------------------------------")
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		// Do not print password or api token but do print the authentication method
		if typeField.Name != "AdguardPassword" {
			log.Println(fmt.Sprintf("%s : %v", typeField.Name, valueField.Interface()))
		} else {
			showAuthenticationMethod(typeField.Name, valueField.String())
		}
	}
	log.Println("---------------------------------------")
}

func showAuthenticationMethod(name, value string) {
	if len(value) > 0 {
		log.Println(fmt.Sprintf("AdGuard Authentication Method : %s", name))
	}
}
