package mysqlx

import(
	"testing"
	"fmt"
)

func Test_batchwhere( t *testing.T ){

	data := []string{"a","b","c","d"}

	sql := SQLStr("select * from table where a= 1").AddParams().AndIn("key",data )

	fmt.Println(sql.Info())

}


