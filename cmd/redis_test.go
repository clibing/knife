package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func TestRedis(t *testing.T) {
	v := make([]string, 0)
	v = append(v, "admin", "clibing", "value")
	sliceTest(v...)

	result := re.FindStringSubmatch(`# Replication
role:master
connected_slaves:0
master_replid:cb1a21229c978f44fd43bfb321e63e282890ac0c
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0`)
	role := result[1]
	if strings.Compare(role, "master") == 0 {
		fmt.Println(role)
	}

}

func sliceTest(hosts ...string) {
	for i, host := range hosts {
		fmt.Println(i, " ", host)
	}

}
