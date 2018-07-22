package models

import (
	sq "github.com/Masterminds/squirrel"
)

type Semester struct {
	ID     uint    `hash:"ignore" json:"-"`
	Name   string
	Gpa    float64
	Scores []Score `db:"-"`
	YearID uint    `hash:"ignore" json:"-" db:"year_id"`
}

var semesterSql = sq.Select("*").From("semesters")

func GetSemesters() (semesters []Semester, err error) {
	sql, args, _ := semesterSql.ToSql()
	err = DB.Select(&semesters, sql, args...)
	return
}

func (year *Year) GetSemesters() (semesters []Semester, err error) {
	sql, args, _ := semesterSql.Where(sq.Eq{"year_id": year.ID}).ToSql()
	err = DB.Select(&semesters, sql, args...)
	for i := range semesters {
		semesters[i].Scores, _ = GetScoresBySemester(&semesters[i])
	}
	return
}
func GetSemestersByYear(year *Year) (semesters []Semester, err error) {
	sql, args, _ := semesterSql.Where(sq.Eq{"year_id": year.ID}).ToSql()
	err = DB.Select(&semesters, sql, args...)
	for i := range semesters {
		semesters[i].Scores, _ = GetScoresBySemester(&semesters[i])
	}
	return
}

func InsertSemester(semester *Semester) (err error) {
	sql, args, _ := sq.Insert("semesters").Columns("year_id", "name", "gpa").Values(semester.YearID, semester.Name, semester.Gpa).ToSql()
	result, err := DB.Exec(sql, args...)
	if err != nil {
		return err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	semester.ID = uint(lastId)
	return
}

func (year *Year) InsertSemesters() {
	stmt, _ := DB.Prepare("INSERT INTO semesters(name, gpa, year_id) VALUES(?,?,?)")
	defer stmt.Close()
	for i, s := range year.Semesters {
		if res, err := stmt.Exec(s.Name, s.Gpa, year.ID); err == nil {
			lastId, _ := res.LastInsertId()
			year.Semesters[i].ID = uint(lastId)
		}
	}
}

func InsertSemesterByYear(semester *Semester, year *Year) (err error) {
	defer func() {
		for _, score := range semester.Scores {
			InsertScoreBySemester(&score, semester)
		}
	}()
	sql, args, _ := sq.Insert("semesters").Columns("year_id", "name", "gpa").Values(year.ID, semester.Name, semester.Gpa).ToSql()
	result, err := DB.Exec(sql, args...)
	if err != nil {
		return err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	semester.ID = uint(lastId)
	return
}

func (year *Year) UpdateSemesters() {
	stmt, _ := DB.Prepare("update semesters set name=?, gpa=? where id=?")
	defer stmt.Close()
	for _, s := range year.Semesters {
		stmt.Exec(s.Name, s.Gpa, s.ID)
	}
}

func UpdateSemester(semester *Semester) (err error) {
	sql, args, _ := sq.Update("semesters").Set("name", semester.Name).Set("gpa", semester.Gpa).Where(sq.Eq{"id": semester.ID}).ToSql()
	_, err = DB.Exec(sql, args...)
	return
}

func DeleteSemester(semester *Semester) (err error) {
	sql, args, _ := sq.Delete("semesters").Where(sq.Eq{"id": semester.ID}).ToSql()
	_, err = DB.Exec(sql, args...)
	return
}
