package models

import (
	sq "github.com/Masterminds/squirrel"
)

func init() {
}

type Score struct {
	ID         uint   `hash:"ignore" json:"-"`
	Name       string
	Type       string
	Credit     float64
	Period     int
	Score      string
	Review     string
	ExamType   string `db:"exam_type"`
	SemesterID uint   `hash:"ignore" json:"-" db:"semester_id"`
}

var scoreSql = sq.Select("*").From("scores")

func GetScores() (scores []Score, err error) {
	sql, _, _ := scoreSql.ToSql()
	err = DB.Select(&scores, sql)
	return
}

func (semester *Semester) GetScores() (scores []Score, err error) {
	sql, args, _ := scoreSql.Where(sq.Eq{"semester_id": semester.ID}).ToSql()
	err = DB.Select(&scores, sql, args...)
	return
}
func GetScoresBySemester(semester *Semester) (scores []Score, err error) {
	sql, args, _ := scoreSql.Where(sq.Eq{"semester_id": semester.ID}).ToSql()
	err = DB.Select(&scores, sql, args...)
	return
}

func InsertScore(score *Score) (err error) {
	sql, args, _ := sq.Insert("scores").Columns("name", "type", "credit", "period", "score", "review", "exam_type", "semester_id").Values(score.Name, score.Type, score.Credit, score.Period, score.Score, score.Review, score.ExamType, score.SemesterID).ToSql()
	result, err := DB.Exec(sql, args...)
	if err != nil {
		return err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	score.ID = uint(lastId)
	return
}

func (semester *Semester) InsertScores() {
	stmt, _ := DB.Prepare("INSERT INTO scores(name, type, credit, period, score, review, exam_type, semester_id) VALUES(?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	for i, s := range semester.Scores {
		if res, err := stmt.Exec(s.Name, s.Type, s.Credit, s.Period, s.Score, s.Review, s.ExamType,
			semester.ID); err == nil {
			lastId, _ := res.LastInsertId()
			semester.Scores[i].ID = uint(lastId)
		}
	}
}

func InsertScoreBySemester(score *Score, semester *Semester) (err error) {
	sql, args, _ := sq.Insert("scores").Columns("name", "type", "credit", "period", "score", "review", "exam_type", "semester_id").Values(score.Name, score.Type, score.Credit, score.Period, score.Score, score.Review, score.ExamType, semester.ID).ToSql()
	result, err := DB.Exec(sql, args...)
	if err != nil {
		return err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	score.ID = uint(lastId)
	return
}

func (semester *Semester) UpdateScores() {
	stmt, _ := DB.Prepare("update scores set name=?, type=?, credit=?, period=?, score=?, review=?, exam_type=? where id=?")
	defer stmt.Close()
	for _, s := range semester.Scores {
		stmt.Exec(s.Name, s.Type, s.Credit, s.Period, s.Score, s.Review, s.ExamType, s.ID)
	}
}

func UpdateScore(score *Score) (err error) {
	sql, args, _ := sq.Update("scores").Set("name", score.Name).Set("type", score.Type).Set("credit", score.Credit).Set("period", score.Period).Set("score", score.Score).Set("review", score.Review).Set("exam_type", score.ExamType).Where(sq.Eq{"id": score.ID}).ToSql()
	_, err = DB.Exec(sql, args...)
	return
}

func DeleteScore(score *Score) (err error) {
	sql, args, _ := sq.Delete("scores").Where(sq.Eq{"id": score.ID}).ToSql()
	_, err = DB.Exec(sql, args...)
	return
}
