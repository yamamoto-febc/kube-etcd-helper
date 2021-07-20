module kube-etcd-helper

go 1.16

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/urfave/cli/v2 v2.3.0
	go.etcd.io/etcd/client/v3 v3.5.0
	k8s.io/apimachinery v0.21.3
	k8s.io/kube-aggregator v0.21.3
	k8s.io/kubectl v0.21.3

)
