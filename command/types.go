package command

type etcd3kv struct {
	Key            string                 `json:"key,omitempty"`
	Value          map[string]interface{} `json:"value,omitempty"`
	CreateRevision int64                  `json:"create_revision,omitempty"`
	ModRevision    int64                  `json:"mod_revision,omitempty"`
	Version        int64                  `json:"version,omitempty"`
	Lease          int64                  `json:"lease,omitempty"`
}
