package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"gopkg.in/urfave/cli.v2"

	"encoding/json"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
	"os"
)

func init() {
	Commands = append(Commands, getCmd)
}

var getCmd = &cli.Command{
	Name:  "get",
	Usage: "Get value",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "pretty",
			Usage:   "Enable JSON pretty format",
			EnvVars: []string{"ETCD_JSON_PRETTY"},
		},
	},
	Action: func(c *cli.Context) error {

		if c.NArg() != 1 {
			return fmt.Errorf("Usage: get <key>")
		}

		client, err := etcdConn()
		if err != nil {
			return err
		}
		defer client.Close() // nolint : error return value not checked

		key := c.Args().First()
		resp, err := clientv3.NewKV(client).Get(context.Background(), key)
		if err != nil {
			return err
		}

		decoder := scheme.Codecs.UniversalDeserializer()
		encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, false)
		objJSON := &bytes.Buffer{}
		for _, kv := range resp.Kvs {
			obj, gvk, err := decoder.Decode(kv.Value, nil, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARN: unable to decode %s: %v\n", kv.Key, err)
				continue
			}
			fmt.Println(gvk)
			objJSON.Reset()
			err = encoder.Encode(obj, objJSON)
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARN: unable to decode %s: %v\n", kv.Key, err)
				continue
			}
			var objMap map[string]interface{}
			if err := json.Unmarshal(objJSON.Bytes(), &objMap); err != nil {
				return nil
			}

			var jsonData []byte
			if c.Bool("pretty") {
				jsonData, err = json.MarshalIndent(objMap, "", "  ")
			} else {
				jsonData, err = json.Marshal(objMap)
			}
			if err != nil {
				return err
			}

			fmt.Println(string(jsonData))
		}

		return nil

	},
}
