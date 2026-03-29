package services

import (
	"database/sql"
	"fmt"
)

type Attendance struct {
	DB *sql.DB
}

func (a *Attendance) getWeeklyAttendance(userId int) {

	fmt.Println("get the weekly attendance for user", userId)
}

func (a *Attendance) getAllAttendance(userId int) {
	fmt.Println("get all attendance for user", userId)
}

func (a *Attendance) clockIn(userId int) {
	fmt.Println("clockin the user", userId)
}

func (a *Attendance) clockOut(userId int) {
	fmt.Println("clockout the user", userId)
}
