package job

import (
	"bufio"
	consts "crontab/internal/const"
	"crontab/internal/logger"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*
* 任务执行
* 开始 结束 日志
 */

var stopCh chan bool = make(chan bool)
var startCh chan bool = make(chan bool)
var configJobs *JobContainer
var runningJobs *JobContainer

func Reset() {
	configJobs = NewJobs()
}

func Run() {
	go jobHandle()
}

func Stop() {
	stopCh <- true
}

func AddAll(jobs *[]Job) (bool, error) {
	for _, job := range *jobs {
		_, err := job.Schedule.Parse(job.Time)
		if err != nil {
			logger.SysPrintf("Err %s %s.\n", err, job.Time)
			return false, err
		}

		jret, err := json.Marshal(job)
		if err != nil {
			logger.SysPrintf("Err %s %s.\n", err, job.Time)
			return false, err
		}

		md5Val := md5.New()
		io.WriteString(md5Val, string(jret))
		hsum := fmt.Sprintf("%x", md5Val.Sum(nil))
		logger.SysPrintf("Adding task [%s], (%s,%v)\n", hsum, job.Comment, job.Time)
		configJobs.add(hsum, &job)
	}
	return true, nil
}

func Add(jstr string) (bool, error) {
	decode := json.NewDecoder(strings.NewReader(jstr))
	var j Job
	if decErr := decode.Decode(&j); decErr != nil {
		logger.SysPrintf("Err %s %s.\n", decErr, jstr)
		return false, decErr
	}
	_, err := j.Schedule.Parse(j.Time)
	if err != nil {
		logger.SysPrintf("Err %s %s.\n", err, jstr)
		return false, err
	}

	h := md5.New()
	io.WriteString(h, jstr)
	hsum := fmt.Sprintf("%x", h.Sum(nil))
	configJobs.add(hsum, &j)

	return true, nil
}

func Delete(key string) {
	configJobs.del(key)
}

func Serialize(current bool) ([]byte, error) {
	if current {
		return runningJobs.json()
	}
	return configJobs.json()
}

func init() {
	Reset()
	runningJobs = NewJobs()
}

func jobHandle() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-stopCh:
			tick.Stop()
			logger.SysWriteLn("Stop crontab")
		case <-startCh:
			tick = time.NewTicker(time.Second)
			logger.SysWriteLn("Start crontab")
		case <-tick.C:
			go configJobs.runJobs()
		}
	}
}

func runJob(j Job) {
	cmd := exec.Command(j.Cmd, j.Args...)
	outpipe, outErr := cmd.StdoutPipe()
	if outErr != nil {
		logger.Printf("[Err] %s %s %s %s\n", j.Cmd, j.Args, j.Out, outErr)
	}
	startErr := cmd.Start()
	if startErr != nil {
		logger.Printf("[Err] %s %s %s %s\n", j.Cmd, j.Args, j.Out, startErr)
		return
	}
	pid := cmd.Process.Pid
	spid := strconv.Itoa(pid)
	j.Start = time.Now().Format(consts.TIMEFORMAT)
	runningJobs.add(spid, &j)
	defer func() {
		runningJobs.del(spid)
		logger.Printf("[End] pid.%d %s %s %s\n", pid, j.Cmd, j.Args, j.Out)
	}()
	logger.Printf("[Start] pid.%d %s %s %s\n", pid, j.Cmd, j.Args, j.Out)
	if j.Out != "" {
		of, ofErr := os.OpenFile(j.Out, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if ofErr != nil {
			logger.Printf("[Err] pid.%d %s %s %s %s", pid, j.Cmd, j.Args, j.Out, ofErr)
		} else {
			defer of.Close()
			outrd := bufio.NewReader(outpipe)
			outrd.WriteTo(of)
		}
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		logger.Printf("[Err] pid.%d %s %s %s %s\n", pid, j.Cmd, j.Args, j.Out, waitErr)
	}
}

func inArray(array []int, item int) bool {
	if len(array) < 1 {
		return false
	}
	if array[0] == -1 {
		return true
	}
	for _, v := range array {
		if item == v {
			return true
		}
	}
	return false
}
