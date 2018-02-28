package command

import "gopkg.in/urfave/cli.v2"

// Config represents CLI config
type Config struct {
	endpoint string
	key      string
	cert     string
	cacert   string
}

// Cfg is a instance of Config in current context
var Cfg = &Config{}
var defaultEndpoint = "http://127.0.0.1:2379"

// CliFlags represents CLI flags
var CliFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "endpoint",
		Usage:       "Etcd endpoint",
		EnvVars:     []string{"ETCD_ENDPOINT"},
		Value:       defaultEndpoint,
		Destination: &Cfg.endpoint,
	},
	&cli.StringFlag{
		Name:        "key",
		Usage:       "Etcd TLS client key",
		EnvVars:     []string{"ETCD_KEY"},
		Destination: &Cfg.key,
	},
	&cli.StringFlag{
		Name:        "cert",
		Usage:       "Etcd TLS client certificate",
		EnvVars:     []string{"ETCD_CERT"},
		Destination: &Cfg.cert,
	},
	&cli.StringFlag{
		Name:        "cacert",
		Usage:       "Etcd server TLS CA certificate",
		EnvVars:     []string{"ETCD_CA_CERT"},
		Destination: &Cfg.cacert,
	},
}
