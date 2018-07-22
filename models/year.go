package models

import (
	sq "github.com/Masterminds/squirrel"
)

func init() {
}

type Year struct {
	ID        uint       `hash:"ignore" json:"-"`
	Name      string
	Gpa       float64
	Semesters []Semester `db:"-"`
	StudentID uint       `hash:"ignore" json:"-" db:"student_id"`
}

var yearSql = sq.Select("*").From("years")

func GetYears() (years []Year, err error) {
	sql, _, _ := yearSql.ToSql()
	err = DB.Select(&years, sql)
	return
}

func (student *Student) GetYears() (years []Year, err error) {
	sql, args, _ := yearSql.Where(sq.Eq{"student_id": student.ID}).ToSql()
	err = DB.Select(&years, sql, args...)
	for i := range years {
		years[i].Semesters, err = GetSemestersByYear(&years[i])
	}
	return
}

func GetYearsByStudent(student *Student) (years []Year, err error) {
	sql, args, _ := yearSql.Where(sq.Eq{"student_id": student.ID}).ToSql()
	err = DB.Select(&years, sql, args...)
	for i := range years {
		years[i].Semesters, err = GetSemestersByYear(&years[i])
	}
	return
}

func InsertYear(year *Year) (err error) {
	sql, args, _ := sq.Insert("years").Columns("student_id", "name", "gpa").Values(year.StudentID, year.Name, year.Gpa).ToSql()
	result, err := DB.Exec(sql, args...)
	if err != nil {
		return err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	year.ID = uint(lastId)
	return
}

func (student *Student) InsertYears() {
	stmt, _ := DB.Prepare("INSERT INTO years(name, gpa, student_id) VALUES(?,?,?)")
	defer stmt.Close()
	for i, y := range student.Years {
		if res, err := stmt.Exec(y.Name, y.Gpa, student.ID); err == nil {
			lastId, _ := res.LastInsertId()
			student.Years[i].ID = uint(lastId)
		}
	}
}

func InsertYearByStudent(year *Year, student *Student) (err error) {
	sql, args, _ := sq.Insert("years").Columns("student_id", "name", "gpa").Values(student.ID, year.Name, year.Gpa).ToSql()
	result, err := DB.Exec(sql, args...)
	if err != nil {
		return err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	year.ID = uint(lastId)
	return
}

func (student *Student) UpdateYears() {
	stmt, _ := DB.Prepare("update years set name=?, gpa=? where id=?")
	defer stmt.Close()
	for _, y := range student.Years {
		stmt.Exec(y.Name, y.Gpa, y.ID)
	}
}

func UpdateYear(year *Year) (err error) {
	sql, args, _ := sq.Update("years").Set("name", year.Name).Set("gpa", year.Gpa).Where(sq.Eq{"id": year.ID}).ToSql()
	_, err = DB.Exec(sql, args...)
	return
}

func DeleteYear(year *Year) (err error) {
	sql, args, _ := sq.Delete("years").Where(sq.Eq{"id": year.ID}).ToSql()
	_, err = DB.Exec(sql, args...)
	return
}
