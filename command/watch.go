package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"gopkg.in/urfave/cli.v2"

	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	Commands = append(Commands, watchCmd)
}

var watchCmd = &cli.Command{
	Name:  "watch",
	Usage: "watch values",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "excludes",
			Usage:   "Exclude prefixes",
			Aliases: []string{"e"},
		},
		&cli.BoolFlag{
			Name:    "pretty",
			Usage:   "Enable JSON pretty format",
			EnvVars: []string{"ETCD_JSON_PRETTY"},
		},
		&cli.StringFlag{
			Name:  "output-dir",
			Usage: "Output in the hierarchical structure of keys",
		},
	},
	Action: func(c *cli.Context) error {

		if c.NArg() > 1 {
			return fmt.Errorf("Usage: watch [<prefix>]")
		}

		client, err := etcdConn()
		if err != nil {
			return err
		}
		defer client.Close() // nolint : error return value not checked

		key := "/"
		if c.NArg() == 1 {
			key = c.Args().First()
		}

		isNeedOutfile := false
		outputDir := c.String("output-dir")
		if len(outputDir) > 0 {
			isNeedOutfile = true
			dir, err := homedir.Expand(filepath.Clean(outputDir))
			if err != nil {
				return err
			}
			if _, err := os.Stat(dir); err == nil {
				return fmt.Errorf("output-dir[%q] is already exists", dir)
			}
			outputDir = dir
		}

		decoder := scheme.Codecs.UniversalDeserializer()
		encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, false)
		objJSON := &bytes.Buffer{}

		rch := client.Watch(context.Background(), key, clientv3.WithPrefix())
		for wresp := range rch {
		loop:
			for _, ev := range wresp.Events {
				kv := ev.Kv
				for _, exclude := range c.StringSlice("excludes") {
					if strings.HasPrefix(string(kv.Key), exclude) {
						continue loop
					}
				}

				if ev.Type == clientv3.EventTypeDelete {
					fmt.Fprintf(os.Stderr, "INFO: key %q is deleted\n", kv.Key)
					continue
				}

				objJSON.Reset()
				obj, _, err := decoder.Decode(kv.Value, nil, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "WARN: error decoding value %q: %v\n	", string(kv.Value), err)
					continue
				}
				if err = encoder.Encode(obj, objJSON); err != nil {
					fmt.Fprintf(os.Stderr, "WARN: error encoding object %#v as 	JSON: %v", obj, err)
					continue
				}
				var objMap map[string]interface{}
				if err := json.Unmarshal(objJSON.Bytes(), &objMap); err != nil {
					return nil
				}
				kvData := etcd3kv{
					Key:            string(kv.Key),
					Value:          objMap,
					CreateRevision: kv.CreateRevision,
					ModRevision:    kv.ModRevision,
					Version:        kv.Version,
					Lease:          kv.Lease,
				}

				var jsonData []byte
				if c.Bool("pretty") {
					jsonData, err = json.MarshalIndent(kvData, "", "  ")
				} else {
					jsonData, err = json.Marshal(kvData)
				}
				if err != nil {
					return err
				}

				if isNeedOutfile {
					key := string(kv.Key)
					if strings.HasPrefix(key, "/") {
						key = strings.TrimLeft(key, "/")
					}

					dir := filepath.Join(outputDir, key)
					// create dir if not exists
					if _, err := os.Stat(dir); err != nil {
						if err = os.MkdirAll(dir, 0755); err != nil {
							return err
						}
					}

					files, err := ioutil.ReadDir(dir)
					if err != nil {
						return err
					}
					fileCount := len(files)
					fileName := fmt.Sprintf("%012d.json", fileCount+1)
					if err := ioutil.WriteFile(filepath.Join(dir, fileName), jsonData, 0755); err != nil {
						return err
					}

				} else {
					fmt.Println(string(jsonData))
				}
			}
		}

		return nil
	},
}
