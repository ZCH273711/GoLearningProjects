package data

import (
	"fmt"
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// CreateSession 给一个已经存在的用户创建一个session
func (user *User) CreateSession() (session Session, err error) {
	statement_insert := "insert into sessions (uuid, email, user_id, created_at) values (?, ?, ?, ?)"
	statement_select := "select * from sessions where id = ?"
	result, err := Db.Exec(statement_insert, CreateUUID(), user.Email, user.Id, time.Now())
	if err != nil {
		fmt.Println("faile to insert session")
		fmt.Println(err.Error())
		return
	}
	session = Session{}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("faile to access last session id")
		fmt.Println(err.Error())
		return
	}
	row := Db.QueryRow(statement_select, id)
	err = row.Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		fmt.Println("faile to select session")
		fmt.Println(err.Error())
		return
	}
	// defer stmt.Close()
	// err = stmt.QueryRow(CreateUUID(), user.Email, user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// Session 获取存在的用户的session
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = ?", user.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// Check 检测session在数据库中是否合法
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = ?", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		fmt.Println("not valid")
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// DeleteByUUID
func (session *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(session.Uuid)
	return
}

// User 从Session中获取User信息
func (session *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = ?", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

func SessionDeleteAll() (err error) {
	statement := "delete from sessions"
	_, err = Db.Exec(statement)
	return
}

func (user *User) Create() (err error) {
	statement := "insert into users (uuid, name, email, password, created_at) values (?, ?, ?, ?, ?)"
	_, err = Db.Exec(statement, CreateUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now())
	if err != nil {
		fmt.Println("failed to prepare statement]")
		fmt.Println(err.Error())
		return
	}
	//defer stmt.Close()
	//_, err = stmt.Exec(CreateUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now())
	//if err != nil {
	//	fmt.Println("failed to execute insert")
	//}
	return
}

func (user *User) Delete() (err error) {
	statement := "delete from users where id = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id)
	return
}

func (user *User) Update() (err error) {
	statement := "update users set name = ?, email = ? where id = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Name, user.Email, user.Id)
	return
}

// Delete all users from database
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// Users 获取表中全部User信息
func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = ?", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func UserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = ?", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}
