package httpkit

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func Test_post2(t *testing.T ){

	httpC :=  NewHTTPClient(12,2)
	robot := "https://oapi.dingtalk.com/robot/send?access_token=3d6371d7e82bb6b5355b8e637443cd624af50af772a0cb20a60a85cd5846994e"
	message :=`{
				  "msgtype": "markdown",
				  "markdown": {
					"title": "xx",
					"text": "@xx，`+"haha"+`"
				  },
				  "at": {
					"atMobiles": [
					  "18122183630"
					],
					"isAtAll": false
				  }
				}`
	var buf io.Reader
	buf = strings.NewReader(message)
	str, err := httpC.Post2(robot, "application/json;charset=utf-8", buf )

	fmt.Println( str, err )


}