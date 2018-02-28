package command

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"gopkg.in/urfave/cli.v2"
)

func init() {
	Commands = append(Commands, listCmd)
}

var listCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List all keys",
	Action: func(c *cli.Context) error {

		if c.NArg() > 1 {
			return fmt.Errorf("Usage: list [<prefix>]")
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

		resp, err := clientv3.NewKV(client).Get(context.Background(),
			key, clientv3.WithFromKey(), clientv3.WithKeysOnly())
		if err != nil {
			return err
		}

		for _, kv := range resp.Kvs {
			fmt.Println(string(kv.Key))
		}
		return nil
	},
}
