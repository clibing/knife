/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

const (
	GET string = "GET"
	SET string = "SET"
	DEL string = "DEL"
)

// redisCmd represents the redis command
var (
	re                = regexp.MustCompile(`role:(.*)`)
	host, value, keys []string // redis server 多服务地址
	password          string   // redis的密码
	database          int      // redis database
	command           string   // redis command
	redisCmd          = &cobra.Command{
		Use:   "redis",
		Short: "redis client",
		Long:  `使用redis发送请求数据.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			for i, s := range host {
				rdb := redis.NewClient(&redis.Options{
					Addr:     s,
					Password: password,
					DB:       database,
				})
				fmt.Printf("%d: %s\n", i, s)
				execute(ctx, rdb, command, value)

			}
		},
	}
)

func execute(ctx context.Context, rdb *redis.Client, cmd string, args []string) (abort bool, err error) {
	fmt.Println(cmd)

	switch cmd {
	case GET:
		for _, key := range args {
			value, err := rdb.Get(ctx, key).Result()
			if err != nil {
				continue
			}
			fmt.Printf("  [%s]-->[%s]\n", key, value)
		}
	case DEL:
		if master(ctx, rdb) {
			for _, key := range args {
				rdb.Del(ctx, key).Result()
			}
		}
	case SET:
		if master(ctx, rdb) {
			size := len(args)
			if size > 0 && size%2 == 0 {
				for i := 0; i < size; {
					fmt.Println("key: ", args[i], " value: ", args[i+1])
					i = i + 2
				}
			}
		}
	default:
		fmt.Println("Not Found")
	}
	return
}

func master(ctx context.Context, rdb *redis.Client) bool {
	replication, _ := rdb.Info(ctx, "replication").Result()
	if len(replication) > 0 {
		result := re.FindStringSubmatch(replication)
		role := result[1]
		role = strings.Replace(role, "\n", "", -1)
		role = strings.Replace(role, "\r", "", -1)
		if strings.Compare(role, "master") == 0 {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(redisCmd)
	redisCmd.Flags().StringSliceVarP(&host, "host", "H", nil, "host list")
	redisCmd.Flags().StringVarP(&password, "password", "p", "", "连接redis的密码")
	redisCmd.Flags().IntVarP(&database, "database", "d", 0, "redis的数据，默认为0")
	redisCmd.Flags().StringVarP(&command, "command", "c", "", `redis执行的命令:
	GET: redis get命令
	SET: redis set命令
	DEL: redis 删除 value slice的信息.`)
	redisCmd.Flags().StringSliceVarP(&value, "value", "v", nil, "执行的参数")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// redisCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// redisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}