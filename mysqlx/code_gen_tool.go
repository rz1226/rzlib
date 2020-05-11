package mysqlx

import (
	"fmt"
	"strings"
)

type OriTableInfo struct {
	COLUMN_NAME    string `orm:"COLUMN_NAME"`
	DATA_TYPE      string `orm:"DATA_TYPE"`
	IS_NULLABLE    string `orm:"IS_NULLABLE"`
	COLUMN_COMMENT string `orm:"COLUMN_COMMENT"`
	EXTRA          string `orm:"EXTRA"`
}

type TableInfo struct {
	Field    string
	DataType string
	Tag      string
	Comment  string
}

//辅助根据表结构生成 bm

func GetBmStrFromTable(dbKit *DB, dbName string, tableName string) string {
	tableInfos := GetTableInfos(dbKit, dbName, tableName)

	str := ""
	for _, v := range tableInfos {
		str += pad(v.Field, 30) + " "
		str += pad(v.DataType, 12) + " "
		str += pad(v.Tag, 20) + " //" + v.Comment

		str += "\n"
	}
	//fmt.Println(str )
	return str
}

func GetTableInfos(dbKit *DB, dbName string, tableName string) []*TableInfo {
	sql := SqlStr("SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE,  COLUMN_COMMENT, EXTRA FROM INFORMATION_SCHEMA.COLUMNS where TABLE_SCHEMA = ? AND table_name = ? order by ordinal_position asc").AddParams(dbName, tableName)

	res, err := sql.Query(dbKit)
	if err != nil {
		fmt.Println("getFieldInfoStr err: ", err)
		return nil
	}
	var oriTableInfo []*OriTableInfo
	err = res.ToStruct(&oriTableInfo)
	if err != nil {
		fmt.Println("getFieldInfoStr err: ", err)
		return nil
	}

	tableInfos := make([]*TableInfo, 0, len(oriTableInfo))
	for _, v := range oriTableInfo {
		tableInfo := &TableInfo{}

		columnName := strings.ToUpper(v.COLUMN_NAME[0:1]) + v.COLUMN_NAME[1:]
		tableInfo.Field = columnName

		dataType := ""
		switch strings.ToUpper(v.DATA_TYPE) {
		case "INT", "BIGINT", "TINYINT", "MEDIUMINT":
			dataType = "int64"
		case "FLOAT", "DOUBLE":
			dataType = "float64"
		case "CHAR", "VARCHAR", "TIME", "TEXT", "DECIMAL", "BLOB", "GEOMETRY", "BIT", "DATETIME", "DATE", "TIMESTAMP":
			dataType = "string"
		default:
			dataType = v.DATA_TYPE
		}
		tableInfo.DataType = dataType

		strTag := "`orm:\"" + v.COLUMN_NAME + "\""
		if v.EXTRA == "auto_increment" {
			strTag += " auto:\"1\""
		}
		strTag += "`"
		tableInfo.Tag = strTag
		tableInfo.Comment = v.COLUMN_COMMENT
		//fmt.Println("数据类型： ",v.DATA_TYPE)
		//fmt.Println("是否为空： ",v.IS_NULLABLE)
		//fmt.Println("字段注释： ",v.COLUMN_COMMENT)
		//fmt.Println("字段额外： ",v.EXTRA)
		tableInfos = append(tableInfos, tableInfo)
	}
	return tableInfos
}

func pad(str string, lenAll int) string {
	length := len(str)
	if length >= lenAll {
		lenAll = length
	}
	if lenPad := lenAll - length; lenPad > 0 {
		return str + strings.Repeat(" ", lenPad)
	}

	return str

}
