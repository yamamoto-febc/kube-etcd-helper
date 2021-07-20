package command

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/coreos/etcd/pkg/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func etcdConn() (*clientv3.Client, error) {

	var tlsConfig *tls.Config
	if len(Cfg.cert) != 0 || len(Cfg.key) != 0 || len(Cfg.cacert) != 0 {
		tlsInfo := transport.TLSInfo{
			CertFile: Cfg.cert,
			KeyFile:  Cfg.key,
			CAFile:   Cfg.cacert,
		}
		var err error
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("ERROR: unable to create client config: %v", err)
		}
	}

	config := clientv3.Config{
		Endpoints:   []string{Cfg.endpoint},
		TLS:         tlsConfig,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		return nil, fmt.Errorf("ERROR: unable to connect to etcd: %v", err)
	}

	return client, nil
}
