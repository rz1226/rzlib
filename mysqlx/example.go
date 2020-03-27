package mysqlx

/**



func testmany() {
	sqlstr := "select * from tai"
	var u  []*repo.Tai
	res := mysqlx.SqlStr(sqlstr).Query(repo.Kit).ToStruct(&u)
	fmt.Println(res)
	fmt.Println("data=")
	for k, v := range u {
		fmt.Println("id=", k, "  v=", v)
	}
	//sql, err := mysqlx.NewBM(&u).ToSqlInsert("tai")
	//fmt.Println(sql.Info() , err)

	//m := mysqlx.NewBM( &u).ToArray()
	//fmt.Println("m1===",m)

}
func testone() {
	sqlstr := "select * from tai"
	u := new(repo.Tai )
	res := mysqlx.SqlStr(sqlstr).AddParams().Query(repo.Kit).ToStruct( u)
	fmt.Println(res)
	fmt.Println("data=", u)

	m := mysqlx.NewBM(u ).ToMap()
	fmt.Println("m==",m)

	var d []*repo.Tai
	mysqlx.Map2Struct(m,&d )
	fmt.Println(mysqlx.NewBM(&d).ToSqlInsert("x"))


	//sql, err := mysqlx.NewBM(u).ToSqlInsert("tai")
	//fmt.Println(sql.Info(), err)



	//sql2, err := mysqlx.NewBM(u).ToSqlUpdate("tai", nil,"")
	//fmt.Println(sql2.Info(), err)


	//condition := mysqlx.SqlStr(" where id =? and name =?").AddParams(2,"mike")
	//sql3 := sql2.ConcatSql(condition )
	//fmt.Println(sql3.Info())

}


*/
