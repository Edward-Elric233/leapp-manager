package common

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"html/template"
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		log.Println(err)
	}
}

func GetIp() (ip string) {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return ip
	}

	for _, a := range ips {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				if strings.HasPrefix(ip, "10") {
					return
				}
				if strings.HasPrefix(ip, "172") {
					return
				}
				if strings.HasPrefix(ip, "192.168") {
					return
				}
				ip = ""
			}
		}
	}
	return
}

var sizeKB = 1024
var sizeMB = sizeKB * 1024
var sizeGB = sizeMB * 1024

func Bytes2Size(num int64) string {
	numStr := ""
	unit := "B"
	if num/int64(sizeGB) > 1 {
		numStr = fmt.Sprintf("%.2f", float64(num)/float64(sizeGB))
		unit = "GB"
	} else if num/int64(sizeMB) > 1 {
		numStr = fmt.Sprintf("%d", int(float64(num)/float64(sizeMB)))
		unit = "MB"
	} else if num/int64(sizeKB) > 1 {
		numStr = fmt.Sprintf("%d", int(float64(num)/float64(sizeKB)))
		unit = "KB"
	} else {
		numStr = fmt.Sprintf("%d", num)
	}
	return numStr + " " + unit
}

func Seconds2Time(num int) (time string) {
	if num/31104000 > 0 {
		time += strconv.Itoa(num/31104000) + " 年 "
		num %= 31104000
	}
	if num/2592000 > 0 {
		time += strconv.Itoa(num/2592000) + " 个月 "
		num %= 2592000
	}
	if num/86400 > 0 {
		time += strconv.Itoa(num/86400) + " 天 "
		num %= 86400
	}
	if num/3600 > 0 {
		time += strconv.Itoa(num/3600) + " 小时 "
		num %= 3600
	}
	if num/60 > 0 {
		time += strconv.Itoa(num/60) + " 分钟 "
		num %= 60
	}
	time += strconv.Itoa(num) + " 秒"
	return
}

func Interface2String(inter interface{}) string {
	switch inter.(type) {
	case string:
		return inter.(string)
	case int:
		return fmt.Sprintf("%d", inter.(int))
	case float64:
		return fmt.Sprintf("%f", inter.(float64))
	}
	return "Not Implemented"
}

func UnescapeHTML(x string) interface{} {
	return template.HTML(x)
}

func IntMax(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func GetUUID() string {
	code := uuid.New().String()
	code = strings.Replace(code, "-", "", -1)
	return code
}

func Max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

// 连接到SSH服务器并执行命令的函数
func ExecuteRemoteCommand(command, hostname, port, username, password string) (string, error) {
	// 设置SSH客户端配置
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// 非生产环境可以使用不安全的HostKeyCallback
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 拼接地址
	address := hostname + ":" + port

	// 连接到SSH服务器
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// 执行命令
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(command)
	if err != nil {
		return stderr.String(), err
	}

	// 返回命令输出
	return stdout.String(), nil
}

// 连接到SSH服务器并执行命令，实时获取命令的输出
// 需要传入一个回调，用来处理实时输出
func ExecuteRemoteCommandRT(command, hostname, port, username, password string, deal func(string)) error {
	// 设置SSH客户端配置
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// 非生产环境可以使用不安全的HostKeyCallback
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 拼接地址
	address := hostname + ":" + port

	// 连接到SSH服务器
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return err
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// 执行命令
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	err = session.Start(command)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Failed to read line: %v", err)
			break
		}
		// 在这里处理读到的每一行日志
		deal(line)
	}

	err = session.Wait()
	if err != nil {
		return err
	}
	return nil
}
