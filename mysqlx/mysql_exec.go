package mysqlx

import (
	dsql "database/sql"
	"errors"
	"fmt"
	"strings"
)

// 代表一个可执行的sql的字符串部分
type SQLStr string

// 生成完整的sql
func (ss SQLStr) AddParams(params ...interface{}) SQL {
	sql := SQL{}
	sql.str = string(ss)
	sql.params = params
	return sql
}

// 不需要参数直接query
/*
func (s Sql) Query(source interface{}) *QueryRes {
	res, error := queryCommon(source, string(s.str), s.params)
	return NewQueryRes(res, error)
}
*/
func (ss SQLStr) Query(source interface{}) (*QueryRes, error) {
	sql := ss.AddParams()
	return sql.Query(source)
}

// 不需要参数直接exec
func (ss SQLStr) Exec(source interface{}) (int64, error) {
	n, err := ss.AddParams().Exec(source)
	return n, err
}

// 代表一个可以执行的sql，一般由两部分组成，str，和变量
type SQL struct {
	str    string
	params []interface{}
}

// 补上一个条件
func (s SQL) ConcatSQL(s2 SQL) SQL {
	s.str += s2.str
	s.params = append(s.params, s2.params...)
	return s
}

// 执行exec   参数是*DB  or *DbTx
func (s SQL) Exec(source interface{}) (int64, error) {
	n, err := execCommon(source, s.str, s.params)
	return n, err
}
func (s SQL) Info() string {
	str := fmt.Sprint("str= ", s.str, "\n params=", s.params)
	return str
}

func execCommon(source interface{}, sqlStr string, args []interface{}) (int64, error) {
	if Conf.Log {
		fmt.Println("running.... exec sql = ", sqlStr, "\n args=", args)
	}
	p, ok := source.(*DB)
	if ok {
		result, err := p.realPool.Exec(sqlStr, args...)
		if err != nil {
			return int64(0), err
		}
		return affectedResult(sqlStr, result)
	}
	t, ok := source.(*DBTx)
	if ok {
		result, err := t.realtx.Exec(sqlStr, args...)
		if err != nil {
			return int64(0), err
		}
		return affectedResult(sqlStr, result)
	}
	return int64(0), errors.New("only support DbPool , DbTx")
}

// 从exec的result获取   当insert获取最后一个id， update，delete获取影响行数，replace获取最后一个id
func affectedResult(sqlStr string, result dsql.Result) (int64, error) {
	if isSQLUpdate(sqlStr) || isSQLDelete(sqlStr) {
		return result.RowsAffected() // 本身就是多个返回值
	}
	if isSQLInsert(sqlStr) {
		return result.LastInsertId() // 本身就是多个返回值
	}
	if isSQLReplace(sqlStr) {
		return result.LastInsertId() // 本身就是多个返回值
	}
	return int64(0), errors.New("only support update insert delete replace")
}

func isSQLReplace(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	return strings.HasPrefix(str, "replace")
}
func isSQLInsert(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	return strings.HasPrefix(str, "insert")
}

func isSQLUpdate(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	return strings.HasPrefix(str, "update")
}

func isSQLDelete(sqlStr string) bool {
	str := strings.TrimSpace(strings.ToLower(sqlStr))
	return strings.HasPrefix(str, "delete")
}
