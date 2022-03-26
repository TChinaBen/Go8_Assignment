package main
//https://www.jianshu.com/p/ee87e989f149
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sqlerror/model/error"
	"sqlerror/model/user"

)

var db *sql.DB
const(
	IP = "127.0.0.1"
	PORT= "3306"
	dbName = "gorm_test"
)
func InitDatabase()*sql.DB{
	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",IP,PORT,dbName)
	db,_ = sql.Open("mysql",dsn)
	return db
}
func SelectUserById(id int) (user.User,*error.WrapError) {
	var user user.User
	QueryString := fmt.Sprintf("SELECT * FROM USERS WHERE id = ?", id)
	rows,err := db.Query(QueryString)
	if err == sql.ErrNoRows{
		errWarpMsg := &error.WrapError{
			Code: "ErrNoRows",
			Msg:  err.Error(),
			Err:  err,
		}
		return user,errWarpMsg
	}
	for rows.Next(){
		rows.Scan(&user.ID,&user.Name,&user.Email,&user.Age,&user.Birthday,&user.MemberNumber,&user.ActivatedAt,&user.CreatedAt,&user.UpdatedAt)
	}
	return user,nil
}
func main(){
    db = InitDatabase()
	fmt.Println(SelectUserById(1))
}