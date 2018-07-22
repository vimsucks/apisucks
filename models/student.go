package models

import (
	sq "github.com/Masterminds/squirrel"
)

func init() {
}

type Student struct {
	ID       uint
	SID      string `db:"sid"`
	Password string
	Years    []Year `db:"-"`
}

var selectSql = sq.Select("*").From("students")

func GetStudents() (students []Student, err error) {
	sql, _, _ := selectSql.ToSql()
	err = DB.Select(&students, sql)
	return
}

func GetStudentByID(id int) (student Student, err error) {
	sql, args, _ := selectSql.Where(sq.Eq{"id": id}).Limit(1).ToSql()
	err = DB.Get(&student, sql, args...)
	return
}

func GetStudentBySID(sid string) (student Student, err error) {
	sql, args, _ := selectSql.Where(sq.Eq{"sid": sid}).Limit(1).ToSql()
	err = DB.Get(&student, sql, args...)
	return
}

func InsertStudent(student *Student) (err error) {
	sql, args, _ := sq.Insert("students").Columns("sid", "password").Values(student.SID, student.Password).ToSql()
	result, err := DB.Exec(sql, args...)
	if err != nil {
		return err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	student.ID = uint(lastId)
	return
}

func InsertStudents(students []Student) {
	stmt, _ := DB.Prepare("INSERT INTO students(sid, password) VALUES(?,?)")
	defer stmt.Close()
	for i, s := range students {
		if res, err := stmt.Exec(s.SID, s.Password); err == nil {
			lastId, _ := res.LastInsertId()
			students[i].ID = uint(lastId)
		}
	}
}

func UpdateStudent(student *Student) (err error) {
	sql, args, _ := sq.Update("students").Set("sid", student.SID).Set("password", student.Password).Where(sq.Eq{"id": student.ID}).ToSql()
	_, err = DB.Exec(sql, args...)
	return
}

func DeleteStudent(student *Student) (err error) {
	sql, args, _ := sq.Delete("students").Where(sq.Eq{"id": student.ID}).ToSql()
	_, err = DB.Exec(sql, args...)
	return
}
