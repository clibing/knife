package client

import (
	"context"
	"fmt"
	"math"
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
			if err == redis.Nil {
				fmt.Printf("key (%s) does not exist\n", key)
				continue
			}
			if err != nil {
				fmt.Println(err.Error())
				return true, err
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
				n, err := rdb.Del(ctx, key).Result()
				if err != nil {
					fmt.Println(err.Error())
					return true, err
				}
				fmt.Printf("del %s: %d\n", key, n)
			}
		} else {
			fmt.Println("当前节点非master, 跳过DEL")
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
					if len(expire) > 0 {
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
		} else {
			fmt.Println("当前节点非master, 跳过SET")
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
		if len(result) < 2 {
			return false
		}
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

	redisCmd.AddCommand(TTLCmd())
	redisCmd.AddCommand(GetCmd())
	redisCmd.AddCommand(SetCmd())
	redisCmd.AddCommand(DelCmd())
	redisCmd.AddCommand(ZSetCmd())
}

func TTLCmd() *cobra.Command {
	redisTTLCmd := &cobra.Command{
		Use:   "ttl",
		Short: "key ttl",
		Long: `获取key对应的ttl

1.获取key的ttl 
  knife client redis ttl key`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			client, key, err := NewRedisClusterClient(ctx, cmd)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer client.Close()
			if len(key) == 0 {
				fmt.Println("key不能为空， 使用--key $key")
				return
			}
			expire, err := cmd.Flags().GetDuration("expire")
			if err != nil {
				fmt.Println("过期时间错误", err)
				return
			}
			if expire > 0 {
				fmt.Println("设置key的过期时间", expire)
				v, err := client.Expire(ctx, key, expire).Result()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("set expire: ", v)
				}
			}

			val, err := client.TTL(ctx, key).Result()
			if err == redis.Nil {
				fmt.Printf("key (%s) does not exist\n", key)
			} else if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("TTL:", val)
			}
		},
	}
	redisTTLCmd.Flags().String("nodes", "192.168.1.8:6391,192.168.1.8:6392,192.168.1.8:6393,192.168.1.8:6394,192.168.1.8:6395,192.168.1.8:6396", "redis cluster node list")
	redisTTLCmd.Flags().String("key", "", "指定key的值")
	redisTTLCmd.Flags().Duration("expire", 0, "指定key的过期时间")
	return redisTTLCmd
}

func GetCmd() *cobra.Command {
	redisGetCmd := &cobra.Command{
		Use:   "get",
		Short: "get key",
		Long: `从redis cluster集群中获取key对应的value

1.获取key: name的值
  knife client redis get name`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			client, key, err := NewRedisClusterClient(ctx, cmd)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer client.Close()
			if len(key) == 0 {
				fmt.Println("key不能为空， 使用--key $key")
				return
			}

			val, err := client.Get(ctx, key).Result()
			if err == redis.Nil {
				fmt.Printf("key (%s) does not exist\n", key)
			} else if err != nil {
				fmt.Println(err)
			} else {
				ttl, _ := client.TTL(ctx, key).Result()

				fmt.Println("value:", val)
				fmt.Println("  ttl:", ttl)
			}
		},
	}
	redisGetCmd.Flags().String("nodes", "192.168.1.8:6391,192.168.1.8:6392,192.168.1.8:6393,192.168.1.8:6394,192.168.1.8:6395,192.168.1.8:6396", "redis cluster node list")
	redisGetCmd.Flags().String("key", "", "根据key获取对应的值")
	return redisGetCmd
}
func SetCmd() *cobra.Command {
	redisSetCmd := &cobra.Command{
		Use:   "set",
		Short: "set key value",
		Long: `从redis cluster集群中设置key对应的value

1.设置key: 对应的值
  knife client redis set name value expire unit`,
		Run: func(cmd *cobra.Command, args []string) {
			value, err := cmd.Flags().GetString("value")
			if err != nil {
				fmt.Println("value error", err)
				return
			}
			expire, err := cmd.Flags().GetDuration("expire")
			if err != nil {
				fmt.Println("expire error", err)
				return
			}
			ctx := context.Background()
			client, key, err := NewRedisClusterClient(ctx, cmd)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer client.Close()
			if len(key) == 0 {
				fmt.Println("key不能为空， 使用--key $key")
				return
			}

			val, err := client.Set(ctx, key, value, expire).Result()
			if err == redis.Nil {
				fmt.Printf("key (%s) does not exist\n", key)
			} else if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(val)
			}
		},
	}
	redisSetCmd.Flags().String("nodes", "192.168.1.8:6391,192.168.1.8:6392,192.168.1.8:6393,192.168.1.8:6394,192.168.1.8:6395,192.168.1.8:6396", "redis cluster node list")
	redisSetCmd.Flags().String("key", "", "key")
	redisSetCmd.Flags().String("value", "", "value")
	redisSetCmd.Flags().Duration("expire", time.Minute, "expire default 60s")
	return redisSetCmd
}

func DelCmd() *cobra.Command {
	redisDelCmd := &cobra.Command{
		Use:   "del",
		Short: "根据key删除redis cluster集群中的数据",
		Long: `从redis cluster集群中删除指定key

1.删除key
  knife client redis del --nodes=127.0.0.1:7000,127.0.0.1:7001,127.0.0.1:7002 --key=key1`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			client, key, err := NewRedisClusterClient(ctx, cmd)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer client.Close()
			if len(key) == 0 {
				fmt.Println("key不能为空， 使用--key $key")
				return
			}

			val, err := client.Del(ctx, key).Result()
			if err == redis.Nil {
				fmt.Printf("key (%s) does not exist\n", key)
			} else if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(val)
			}
		},
	}
	redisDelCmd.Flags().String("nodes", "192.168.1.8:6391,192.168.1.8:6392,192.168.1.8:6393,192.168.1.8:6394,192.168.1.8:6395,192.168.1.8:6396", "redis cluster node list")
	redisDelCmd.Flags().String("key", "", "删除指定的key")
	return redisDelCmd
}

func ZSetCmd() *cobra.Command {
	zSetCmd := &cobra.Command{
		Use:   "zset",
		Short: "zset相关操作",
		Long: `从redis cluster集群中操作zset数据

1.删除key
  knife client redis del --nodes=127.0.0.1:7000,127.0.0.1:7001,127.0.0.1:7002 --key=key1`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	zSetCmd.Flags().String("nodes", "192.168.1.8:6391,192.168.1.8:6392,192.168.1.8:6393,192.168.1.8:6394,192.168.1.8:6395,192.168.1.8:6396", "redis cluster node list")
	zSetCmd.Flags().String("key", "", "zset key")

	zSetAdd := &cobra.Command{
		Use:   "add",
		Short: "zset添加",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			client, key, err := NewRedisClusterClient(ctx, cmd)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer client.Close()

			member, err := cmd.Flags().GetString("member")
			if err != nil {
				fmt.Println("member error", err)
				return
			}

			score, err := cmd.Flags().GetFloat64("score")
			if err != nil {
				fmt.Println("score error", err)
				return
			}

			expire, err := cmd.Flags().GetDuration("expire")
			if err != nil {
				fmt.Println("expire error", err)
				return
			}
			if expire > 0 {
				if _, err := client.Expire(ctx, key, expire).Result(); err != nil {
					fmt.Println("expire error", err)
					return
				}
			}

			val, err := client.ZAdd(ctx, key, redis.Z{
				Score:  score,
				Member: member,
			}).Result()
			if err == redis.Nil {
				fmt.Printf("key (%s) does not exist\n", key)
			} else if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(val)
			}

		},
	}
	zSetAdd.Flags().String("nodes", "192.168.1.8:6391,192.168.1.8:6392,192.168.1.8:6393,192.168.1.8:6394,192.168.1.8:6395,192.168.1.8:6396", "redis cluster node list")
	zSetAdd.Flags().String("key", "", "zset的key")
	zSetAdd.Flags().String("member", "", "zset的member")
	zSetAdd.Flags().Float64("score", float64(time.Now().UnixMilli()), "zset的member对应的score")
	zSetAdd.Flags().Duration("expire", 0, "expire default 0s(不过期)")

	zSetGet := &cobra.Command{
		Use:   "get",
		Short: "zset获取",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			client, key, err := NewRedisClusterClient(ctx, cmd)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer client.Close()

			start, err := cmd.Flags().GetInt64("start")
			if err != nil {
				fmt.Println("start error", err)
				return
			}
			stop, err := cmd.Flags().GetInt64("stop")
			if err != nil {
				fmt.Println("stop error", err)
				return
			}
			offset, err := cmd.Flags().GetInt64("offset")
			if err != nil {
				fmt.Println("offset error", err)
				return
			}
			count, err := cmd.Flags().GetInt64("count")
			if err != nil {
				fmt.Println("count error", err)
				return
			}
			rev, err := cmd.Flags().GetBool("rev")
			if err != nil {
				fmt.Println("rev error", err)
				return
			}

			members, err := client.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
				Key:     key,
				Start:   start,
				Stop:    stop,
				ByScore: true,
				Offset:  offset,
				Count:   count,
				Rev:     rev,
			}).Result()
			if err == redis.Nil {
				fmt.Printf("key (%s) does not exist\n", key)
			} else if err != nil {
				fmt.Println(err)
			} else {
				for _, member := range members {
					fmt.Println(member.Member, " ", member.Score)
				}
			}
		},
	}
	zSetGet.Flags().String("nodes", "192.168.1.8:6391,192.168.1.8:6392,192.168.1.8:6393,192.168.1.8:6394,192.168.1.8:6395,192.168.1.8:6396", "redis cluster node list")
	zSetGet.Flags().String("key", "", "zset的key")
	zSetGet.Flags().Int64("start", -math.MaxInt64, "start score")
	zSetGet.Flags().Int64("stop", math.MaxInt64, "stop score")
	zSetGet.Flags().Int64("offset", 0, "offset 偏移量")
	zSetGet.Flags().Int64("count", 32, "count 分页数")
	zSetGet.Flags().Bool("rev", true, "rev 反转")

	zSetCmd.AddCommand(
		zSetAdd,
		zSetGet,
	)
	return zSetCmd
}

func NewRedisClusterClient(ctx context.Context, cmd *cobra.Command) (rdb *redis.ClusterClient, key string, err error) {
	nodes, err := cmd.Flags().GetString("nodes")
	if err != nil {
		fmt.Println(err)
		return
	}
	key, err = cmd.Flags().GetString("key")
	if err != nil {
		fmt.Println(err)
		return
	}

	addrs := strings.Split(nodes, ",")
	rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
	})

	// 连接测试
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("无法连接 Redis 集群")
		return
	}

	// fmt.Println("成功连接到 Redis 集群")
	return
}

func NewRedisClient() *cobra.Command {
	return redisCmd
}
