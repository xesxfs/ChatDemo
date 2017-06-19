package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id      int
	Name    string
	Profile *Profile `orm:"rel(one)"`
	Post    []*Post  `orm:"reverse(many)"`
}

type Profile struct {
	Id   int
	Age  int16
	User *User `orm:"reverse(one)"`
}

type Post struct {
	Id    int
	Title string
	User  *User  `orm:"rel(fk)"`
	Tags  []*Tag `orm:"rel(m2m)"`
}

type Tag struct {
	Id    int
	Name  string
	Posts []*Post `orm:"reverse(many)"`
}

func init() {
	// orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:think123@/orm_test?charset=utf8")
	orm.RegisterModel(new(User), new(Post), new(Profile), new(Tag))

}

func main() {
	orm.Debug = true
	orm.RunSyncdb("default", false, true)
	o := orm.NewOrm()
	o.Using("default")

	profile := new(Profile)
	profile.Age = 22

	user := new(User)
	user.Profile = profile
	user.Name = "Lucy"
	fmt.Println("---------------------------")

	fmt.Println(o.Insert(profile))
	fmt.Println(o.Insert(user))

}
