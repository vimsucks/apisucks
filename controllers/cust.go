package controllers

import (
	"github.com/labstack/echo"
	"strconv"
	"net/http"
	"fmt"
	"net/http/cookiejar"
	"time"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"errors"
	"github.com/labstack/gommon/log"
	. "github.com/vimsucks/apisucks/models"
	"github.com/vimsucks/apisucks/util"
	"github.com/mitchellh/hashstructure"
)

const jwglLoginUrl = "http://jwgl.cust.edu.cn/teachwebsl/login.aspx"
const jwglScoreUrl = "http://jwgl.cust.edu.cn/teachweb/cjcx/StudentGrade.aspx"

var JwglConnectionError = errors.New("error when connecting to jwgl")
var JwglLoginError = errors.New("user id or password wrong")

const (
	ModeAll      = iota
	ModeYear
	ModeSemester
)

// METHOD: POST
// URL: /cust/:sid/score/all
// FORM_FIELDS: [password]
func GetAllScore(c echo.Context) (err error) {
	return getScore(c, ModeAll)
}

// METHOD: POST
// URL: /cust/:sid/score/year/:year
// FORM_FIELDS: [password]
func GetYearScore(c echo.Context) (err error) {
	return getScore(c, ModeYear)
}

// METHOD: POST
// URL: /cust/:sid/score/year/:year/semester/:semester
// FORM_FIELDS: [password]
func GetYearSemesterScore(c echo.Context) (err error) {
	return getScore(c, ModeSemester)
}

func getScore(c echo.Context, mode int) (err error) {
	sid := c.Param("sid")
	password := c.FormValue("password")
	var year, semester int
	if mode == ModeYear || mode == ModeSemester {
		if year, err = strconv.Atoi(c.Param("year")); err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "year should be int (e.g., /cust/2015001234/year/1")
		}
	}
	if mode == ModeSemester {
		if semester, err = strconv.Atoi(c.Param("semester")); err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "semester should be int (e.g., /cust/2015001234/year/1")
		}
	}

	var student Student
	if student, err = GetStudentBySID(sid); err == nil {
		if student.Password != password {
			return c.String(http.StatusBadRequest, JwglLoginError.Error())
		}
	}

	var client *http.Client
	if client, err = loginHttpClient(sid, password); err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	log.Info(fmt.Sprintf("student %s login successfully", sid))

	student = Student{SID: sid, Password: password}
	InsertStudent(&student)

	var years []Year
	if years, err = getYears(client); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	} else {
		if year == -1 {
			year = len(years)
		}
		if semester == -1 {
			semester = len(years[year-1].Semesters)
		}
		switch mode {
		case ModeAll:
			return c.JSON(http.StatusOK, years)
		case ModeYear:
			return c.JSON(http.StatusOK, years[year-1])
		case ModeSemester:
			return c.JSON(http.StatusOK, years[year-1].Semesters[semester-1])
		default:
			return errors.New("unrecognized mode")
		}
	}
}

func newRegister(client *http.Client, student *Student) (years []Year, err error) {
	if err = InsertStudent(student); err != nil {
		return nil, err
	}
	if years, err = getYears(client); err != nil {
		return nil, err
	}
	student.Years = years
	go func() {
		student.InsertYears()
		for i := range student.Years {
			student.Years[i].InsertSemesters()
			for j := range student.Years[i].Semesters {
				fmt.Println(student.Years[i].Semesters[j].ID)
				student.Years[i].Semesters[j].InsertScores()
			}
		}
	}()
	return years, nil
}

func checkUpdate(client *http.Client, student *Student) (years []Year, err error) {
	if student.Years, err = student.GetYears(); err != nil {
		return nil, err
	}
	if years, err = getYears(client); err != nil {
		return nil, err
	}
	fmt.Println(hashstructure.Hash(student.Years, nil))
	fmt.Println(hashstructure.Hash(years, nil))
	i := 0
	if util.CompareStruct(student.Years, years) {
		return years, nil
	}
	for ; i < len(student.Years); i++ {
		if ! util.CompareStruct(student.Years[i], years[i]) {
			student.Years[i].Gpa = years[i].Gpa
			UpdateYear(&student.Years[i])

			j := 0
			for ; j < len(student.Years[i].Semesters); j++ {
				if !util.CompareStruct(student.Years[i].Semesters[j], years[i].Semesters[j]) {
					student.Years[i].Semesters[j].Gpa = years[i].Semesters[j].Gpa
					UpdateSemester(&student.Years[i].Semesters[j])

					k := 0
					for ; k < len(student.Years[i].Semesters[j].Scores); k++ {
						if !util.CompareStruct(student.Years[i].Semesters[j].Scores[k], years[i].Semesters[j].Scores[k]) {
							years[i].Semesters[j].Scores[k].ID = student.Years[i].Semesters[j].Scores[k].ID
							years[i].Semesters[j].Scores[k].SemesterID = student.Years[i].Semesters[j].Scores[k].SemesterID
							student.Years[i].Semesters[j].Scores[k] = years[i].Semesters[j].Scores[k]
							UpdateScore(&student.Years[i].Semesters[j].Scores[k])
						}
					}

					for ; k < len(years[i].Semesters[j].Scores); k++ {
						InsertScoreBySemester(&years[i].Semesters[j].Scores[k], &years[i].Semesters[j])
					}
				}
			}
			for ; j < len(years[i].Semesters); j++ {
				InsertSemesterByYear(&years[i].Semesters[j], &student.Years[i])
				(&years[i].Semesters[j]).InsertScores()
			}
		}
	}
	for ; i < len(years); i++ {
		InsertYearByStudent(&years[i], student)
		years[i].InsertSemesters()
		for _, semester := range years[i].Semesters {
			semester.InsertScores()
		}
	}
	return years, nil
}

func getYears(client *http.Client) (years []Year, err error) {
	resp, err := client.Get(jwglScoreUrl)
	if err != nil {
		return nil, JwglConnectionError
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	var year Year
	var semester Semester
	var score Score
	firstYear := true
	doc.Find("tr").Next().Each(func(k int, s *goquery.Selection) {
		selTd := s.Find("td")
		length := selTd.Length()
		selTd = selTd.First()
		if length == 12 {
			// 新学年
			if firstYear {
				firstYear = false
			} else {
				year.Semesters = append(year.Semesters, semester)
				years = append(years, year)
			}
			yearGpa, _ := strconv.ParseFloat(selTd.Next().Text(), 64)
			year = Year{Name: selTd.Text(), Gpa: yearGpa}
			semesterGpa, _ := strconv.ParseFloat(selTd.Next().Next().Next().Text(), 64)
			semester = Semester{Name: selTd.Next().Next().Text(), Gpa: semesterGpa}
			score = parseScore(selTd.Next().Next().Next().Next())
			semester.Scores = append(semester.Scores, score)
		} else if length == 10 {
			// 新学期
			year.Semesters = append(year.Semesters, semester)
			semesterGpa, _ := strconv.ParseFloat(selTd.Next().Text(), 64)
			semester = Semester{Name: selTd.Text(), Gpa: semesterGpa}
			score = parseScore(selTd.Next().Next())
			semester.Scores = append(semester.Scores, score)
		} else {
			score = parseScore(selTd)
			semester.Scores = append(semester.Scores, score)
		}
	})
	year.Semesters = append(year.Semesters, semester)
	years = append(years, year)

	return years, nil
}

func getYearsInDB(student *Student) (years []Year, err error) {
	if years, err = GetYearsByStudent(student); err != nil {
		return nil, err
	}
	for i, _ := range years {
		if years[i].Semesters, err = GetSemestersByYear(&years[i]); err != nil {
			return nil, err
		}
		for j, _ := range years[i].Semesters {
			if years[i].Semesters[j].Scores, err = GetScoresBySemester(&years[i].Semesters[j]); err != nil {
				return nil, err
			}
		}
	}
	return years, nil
}

// 使用成绩行的第一个 td 元素，解析出成绩 struct
func parseScore(s *goquery.Selection) Score {
	score := Score{}
	score.Name = s.Text()
	score.Type = s.Next().Text()
	credit, _ := strconv.ParseFloat(s.Next().Next().Text(), 64)
	score.Credit = credit
	period, _ := strconv.Atoi(s.Next().Next().Next().Text())
	score.Period = period
	score.Score = s.Next().Next().Next().Next().Text()
	score.Review = s.Next().Next().Next().Next().Next().Text()
	score.ExamType = s.Next().Next().Next().Next().Next().Next().Text()
	return score
}

func loginHttpClient(sid, password string) (*http.Client, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
	// 先 GET 一遍登录页面，获取登录 Form 需要的 __EVENTVALIDATION 和 __VIEWSTATE 值
	resp, err := client.Get(jwglLoginUrl)
	if err != nil || resp.StatusCode != 200 {
		return nil, JwglConnectionError
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	eventValidation, _ := doc.Find("#__EVENTVALIDATION").Attr("value")
	viewState, _ := doc.Find("#__VIEWSTATE").Attr("value")

	// 构造表单
	form := url.Values{}
	form.Add("__EVENTVALIDATION", eventValidation)
	form.Add("__VIEWSTATE", viewState)
	form.Add("txtUserName", sid)
	form.Add("txtPassWord", password)
	form.Add("Button1", "登录")
	resp, err = client.PostForm(jwglLoginUrl, form)
	if err != nil || resp.StatusCode != 200 {
		return nil, JwglConnectionError
	}

	// 检查是否登录成功，若成功会重定向至 index，检查首页是否出现学生姓名
	doc, _ = goquery.NewDocumentFromReader(resp.Body)
	if doc.Find("#StudentNameValueLabel").Length() == 0 {
		return nil, JwglLoginError
	}
	return client, nil
}
