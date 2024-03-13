package model

import (
	"errors"
)

type Task struct {
	Id       int    `json:"id"`
	TaskName string `json:"task_name" gorm:"unique;index" validate:"max=20"`
	Ip       string `json:"ip" gorm:"unique;index" validate:"min=7,max=15"`
	Port     int    `json:"port" gorm:"type:int;default:22"` //ssh端口号
	Status   int    `json:"status" gorm:"type:int"`          //0-未开始 1-正在升级 2-升级成功 3-升级失败
	//TODO: creator 根据创建者进行筛选
}

func GetMaxTaskId() int {
	var task Task
	DB.Last(&task)
	return task.Id
}

func GetAllTasks(startIdx int, num int) (tasks []*Task, err error) {
	err = DB.Select([]string{"id", "task_name", "ip", "port", "status"}).Offset(startIdx).Limit(num).Order("id desc").Find(&tasks).Error
	return tasks, err
}

func SearchTasks(keyword string) (tasks []*Task, err error) {
	err = DB.Select([]string{"id", "task_name", "ip", "port", "status"}).Where("id = ? or task_name LIKE ?", keyword, keyword+"%").Order("id desc").Find(&tasks).Error
	return tasks, err
}

func GetTaskById(id int, selectAll bool) (*Task, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	task := Task{Id: id}
	var err error = nil
	if selectAll {
		err = DB.First(&task, "id = ?", id).Error
	} else {
		err = DB.Select([]string{"id", "task_name", "ip", "port", "status"}).First(&task, "id = ?", id).Error
	}
	return &task, err
}

func DeleteTaskById(id int) (err error) {
	if id == 0 {
		return errors.New("id 为空！")
	}
	task := Task{Id: id}
	return task.Delete()
}

func (task *Task) Insert() error {
	var err error
	err = DB.Create(task).Error
	return err
}

func (task *Task) Update() error {
	var err error
	err = DB.Model(task).Updates(task).Error
	return err
}

func (task *Task) Delete() error {
	if task.Id == 0 {
		return errors.New("id 为空！")
	}
	err := DB.Delete(task).Error
	return err
}
