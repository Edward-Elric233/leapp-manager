package controller

import (
	"encoding/json"
	"fmt"
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

	//TODO: 在系统配置中设置SERVER_IP
	//环境变量位置需要和leapp-repository项目约定
	serverIp := "192.168.174.1"
	serverPort := 3000
	command := fmt.Sprintf("echo -e \"SERVER_IP=%s\\nSERVER_PORT=%d\\nTASK_ID=%d\" | sudo tee /root/.leapp.env > /dev/null && cat /etc/os-release && uname -r", serverIp, serverPort, id)
	// 执行命令
	output, err := common.ExecuteRemoteCommand(command, task.Ip, strconv.Itoa(task.Port), username, password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	//连接成功
	task.Status = common.TaskRun
	task.Info = output

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

	go func() {
		command = "dnf install git -y && rm leapp-repository -rf && git clone https://github.com/Edward-Elric233/leapp-repository.git && cd leapp-repository && bash run.sh&"
		// 执行命令
		output, err = common.ExecuteRemoteCommand(command, task.Ip, strconv.Itoa(task.Port), username, password)
		if err != nil {
			log.Printf("Failed to run leapp command for task %d: %v", id, err)
			return
		}
	}()

	//TODO: 放入RegisterWebSocket以实现即使用户退出页面，再进来也能看到日志
	go func() {
		conn, ok := task2websocket[id]
		if !ok {
			log.Printf("Failed to get websocket for task %d", id)
			return
		}
		deal := func(msg string) {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}
		}
		//TODO: 持久化log，这样即使退出页面再次进来还能获取以前的输出
		logFile := "/var/log/leapp/leapp-upgrade-stdout.log"
		command := fmt.Sprintf("while ! [ -f  %s ]; do sleep 1; done && tail -f %s", logFile, logFile)
		err = common.ExecuteRemoteCommandRT(command, task.Ip, strconv.Itoa(task.Port), username, password, deal)
		if err != nil {
			log.Printf("Failed to update log: %v", err)
			return
		}
	}()

	return
}

//TODO: 检测升级任务失败

func FinishTask(c *gin.Context) {
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

	//升级成功
	task.Status = common.TaskSuccess

	task.Info += "\n升级成功\n"
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

	write2WebSocket(id, "upgrade success!!!")

	return
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var task2websocket = make(map[int]*websocket.Conn)

func write2WebSocket(id int, msg string) {
	conn, ok := task2websocket[id]
	if !ok {
		log.Printf("Failed to get websocket for task %d", id)
		return
	}
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
}

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
	defer conn.Close()
	if !ok {
		log.Printf("Failed to get websocket for task %d", id)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}
