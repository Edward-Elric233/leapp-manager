package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"leapp-manager/common"
	"leapp-manager/model"
	"log"
	"net/http"
	"strconv"
)

func GetAllTasks(c *gin.Context) {
	p, _ := strconv.Atoi(c.Query("p"))
	if p < 0 {
		p = 0
	}
	tasks, err := model.GetAllTasks(p*common.ItemsPerPage, common.ItemsPerPage)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    tasks,
	})
	return
}

func SearchTasks(c *gin.Context) {
	keyword := c.Query("keyword")
	tasks, err := model.SearchTasks(keyword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    tasks,
	})
	return
}

func GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	task, err := model.GetTaskById(id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    task,
	})
	return
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	err = model.DeleteTaskById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "删除任务成功",
		})
	}
}

func CreateTask(c *gin.Context) {
	var task model.Task
	err := json.NewDecoder(c.Request.Body).Decode(&task)
	if err != nil || task.TaskName == "" || task.Ip == "" || task.Port == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if err := task.Insert(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

//TODO: func UpdateTask

func StartTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	task, err := model.GetTaskById(id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	//TODO: 将用户名和密码放入数据库，需要用户输入
	username := "root"
	password := "Tlinux12#$"
	command := "cat /etc/os-release"
	// 执行命令
	output, err := common.ExecuteRemoteCommand(command, task.Ip, strconv.Itoa(task.Port), username, password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	//连接成功
	task.Status = common.TaskRun
	task.Info = output

	//TODO: run shell script
	command = "dnf install git -y && git clone https://gitee.com/EdwardElric233/leapp-repository.git && cd leapp-repository && bash run.sh"
	// 执行命令
	output, err = common.ExecuteRemoteCommand(command, task.Ip, strconv.Itoa(task.Port), username, password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	task.Info += "\n开始升级"
	err = task.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    task,
	})

	return
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var task2websocket = make(map[int]*websocket.Conn)

func RegisterWebSocket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %v", err)
		return
	}
	//defer conn.Close()  -> in removeWebSocket
	task2websocket[id] = conn

	err = conn.WriteMessage(websocket.TextMessage, []byte("Socket has been created"))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}

	// 处理 WebSocket 连接
	//for {
	//	t, msg, err := conn.ReadMessage()
	//	if err != nil {
	//		break
	//	}
	//	log.Printf("Received message: %s", msg)
	//	err = conn.WriteMessage(t, msg)
	//	if err != nil {
	//		break
	//	}
	//}
}

func RemoveWebSocket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	conn, ok := task2websocket[id]
	if !ok {
		log.Printf("Failed to get websocket for task %d: %v", id, err)
		return
	}
	defer conn.Close()
}
