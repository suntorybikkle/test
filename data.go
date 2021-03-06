package main

import (
	"database/sql"
	// "fmt"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=ldb dbname=ldb password=ldb sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func GetLastStudyInfo(id int) (studyInfo StudyInfo, err error) {
	studyInfo = StudyInfo{}
	err = Db.QueryRow("SELECT id, user_id, subject_id, study_time, date_time FROM study_infos WHERE user_id = $1 ORDER BY Id DESC LIMIT 1", id).Scan(
		&studyInfo.Id, &studyInfo.UserId, &studyInfo.SubjectId, &studyInfo.StudyTime, &studyInfo.DateTime)
	return
}

func GetAllStudyInfo(userId int) (studyInfos []StudyInfo, err error) {
	rows, err := Db.Query("SELECT id, user_id, subject_id, study_time, date_time FROM study_infos WHERE user_id = $1", userId)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		studyInfo := StudyInfo{}
		err = rows.Scan(&studyInfo.Id, &studyInfo.UserId, &studyInfo.SubjectId, &studyInfo.StudyTime, &studyInfo.DateTime)
		if err != nil {
			log.Println(err)
			return
		}
		studyInfos = append(studyInfos, studyInfo)
	}
	return
}

func (studyInfo *StudyInfo) Create() (err error) {
	statement := "INSERT INTO study_infos(user_id, subject_id, study_time, date_time) VALUES($1, $2, $3, $4) RETURNING id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(studyInfo.UserId, studyInfo.SubjectId, studyInfo.StudyTime, studyInfo.DateTime).Scan(&studyInfo.Id)
	return
}
