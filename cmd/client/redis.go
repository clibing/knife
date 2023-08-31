package client


import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

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
	re                          = regexp.MustCompile(`role:(.*)`)
	host, values, keys, expires []string // redis server 多服务地址
	password                    string   // redis的密码
	database                    int      // redis database
	command                     string   // redis command
	sentinel                    bool     // redis sentinel model
	redisCmd                    = &cobra.Command{
		Use:   "redis",
		Short: "redis client",
		Long: `使用redis发送请求数据

1.快速获取key: name的值
  knife redis -c GET -k name

2.获取key: name, userr的值
  knife redis -H 127.0.0.1:6379 -c GET -k name -k user

3.设置key: name, value: 123456, expire: 120s
  knife redis -H 127.0.0.1:6379 -c SET -k name -v 123456 -e 120s

4.删除key: name, user
  knife redis -c DEL -k name -k user

5.链接多redis服务，获取key: name
  knife redis -H 127.0.0.1:6379 -H 127.0.0.1:6380 -c GET -k name

6.删除哨兵的key, 会自动判断当前redis server的role是否为master
  knife redis -H 127.0.0.1:6379 -H 127.0.0.1:6380 -c DEL -k name

7.当链接哨兵模式时，读取配置时，只返回一个key对应的值
  knife redis -H 127.0.0.1:6379 -H 127.0.0.1:6380 -s -c GET -k name.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			size := len(host)
			for i, s := range host {
				rdb := redis.NewClient(&redis.Options{
					Addr:     s,
					Password: password,
					DB:       database,
				})
				if size > 1 {
					fmt.Printf("%d: %s role(%s)\n", i, s, role(ctx, rdb))
				}
				abort, err := execute(ctx, rdb, command, keys, values, expires, args)
				if err != nil {
					fmt.Println("err message: ", err)
				}
				if abort {
					break
				}
			}
		},
	}
)

func execute(ctx context.Context, rdb *redis.Client, cmd string, keys []string, values []string, expires []string, args []string) (abort bool, err error) {
	keySize := len(keys)
	if keySize == 0 {
		fmt.Println("请输入key")
		return
	}
	switch cmd {
	case GET:
		for _, key := range keys {
			value, err := rdb.Get(ctx, key).Result()
			if err != nil {
				continue
			}
			ttl, _ := rdb.TTL(ctx, key).Result()
			fmt.Printf("  [%s]-->[%s] expire: %ds\n", key, value, int(ttl.Seconds()))
			if sentinel {
				abort = true
				break
			}
		}
	case DEL:
		if master(ctx, rdb) {
			for _, key := range keys {
				rdb.Del(ctx, key).Result()
			}
		}
	case SET:
		if master(ctx, rdb) {
			valueSize := len(values)
			expireSize := len(expires)
			if keySize > 0 && keySize == valueSize && valueSize == expireSize {
				for i := 0; i < keySize; i++ {
					key := keys[i]
					value := values[i]
					expire := expires[i]
					var expireDuration time.Duration
					if len(expire) == 0 {
						expireDuration, _ = time.ParseDuration(expire)
					} else {
						expireDuration, err = time.ParseDuration(expire)
						if err != nil {
							abort = true
							fmt.Println("expire format error: " + expire)
							return
						}
					}
					var result string
					result, err = rdb.Set(ctx, key, value, expireDuration).Result()
					if err != nil {
						fmt.Println(err.Error())
						abort = true
						return
					}
					fmt.Printf("%s. set %s %v\r\n", result, key, value)
				}
			} else {
				fmt.Println("key, value, expire丢失")
				abort = true
				return
			}
		}
	default:
		fmt.Println("Not Found")
	}
	return
}

type Parameter struct {
	Key      string
	Value    interface{}
	Duration time.Duration
}

func ParseArgs(args []string) ([]*Parameter, error) {
	parameters := make([]*Parameter, 0, len(args))
	for i, _ := range args {
		parameters[i] = nil
	}
	return parameters, nil
}

func role(ctx context.Context, rdb *redis.Client) string {
	if master(ctx, rdb) {
		return "master"
	}
	return "slave"
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
	redisCmd.Flags().StringSliceVarP(&host, "host", "H", []string{"127.0.0.1:6379"}, "host list")
	redisCmd.Flags().StringVarP(&password, "password", "p", "", "连接redis的密码")
	redisCmd.Flags().IntVarP(&database, "database", "d", 0, "redis的数据，默认为0")
	redisCmd.Flags().BoolVarP(&sentinel, "sentinel", "s", false, "是否sentinel模式")
	redisCmd.Flags().StringVarP(&command, "command", "c", "", "redis执行的命令: GET SET DEL")
	redisCmd.Flags().StringSliceVarP(&keys, "key", "k", nil, "key list")
	redisCmd.Flags().StringSliceVarP(&values, "value", "v", nil, "value list")
	redisCmd.Flags().StringSliceVarP(&expires, "expire", "e", nil, "expire list")
}

func NewRedisClient() *cobra.Command {
	return redisCmd
}
