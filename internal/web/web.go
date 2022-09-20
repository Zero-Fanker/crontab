package web

import (
	"bufio"
	"crontab/internal/conf"
	consts "crontab/internal/const"
	"crontab/internal/job"
	"crontab/internal/logger"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

/*
*  get	获取任务列表
*  set	设置任务/添加任务
*  del	删除任务
*  log	任务执行日志
*  load 重新加载任务列表
*  stop 停止任务触发，正在执行的任务正常执行
*  start开始任务触发，
*  status获取正在执行的任务
 */

func Get(w http.ResponseWriter, r *http.Request) {
	allJobs, err := job.Serialize(false)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", allJobs)
	}
}

func Set(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	// h := r.FormValue("h")
	j := r.FormValue("j")
	j = strings.TrimSpace(j)
	if j == "" {
		fmt.Fprintf(w, "%s", "job empty")
	}
	successful, err := job.Add(j)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	} else if !successful {
		fmt.Fprintf(w, "undefined error")
		return
	}
	_, err = conf.Save()
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	fmt.Fprintf(w, "%s", "success")
}

func Del(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	h := r.FormValue("h")
	h = strings.TrimSpace(h)
	job.Delete(h)
	_, err = conf.Save()
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", "success")
	}
}

func Loger(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	dayNum := r.FormValue("d")

	reg := regexp.MustCompile(`^[0-9]{8}$`)
	b := reg.MatchString(dayNum)

	if !b {
		fmt.Fprintf(w, "%s", "invalid day")
		return
	}
	file := logger.FilePath() + dayNum + "_" + consts.RUN_LOG_POSTFIX

	fp, err := os.Open(file)

	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	defer fp.Close()
	rd := bufio.NewReader(fp)
	rd.WriteTo(w)
}

func Load(w http.ResponseWriter, r *http.Request) {
	loaded, loadErr := conf.Load(nil)
	if loaded {
		fmt.Fprintf(w, "%s", "success")
	} else {
		fmt.Fprintf(w, "%s", loadErr)
	}

}

func Status(w http.ResponseWriter, r *http.Request) {
	brunning, err := job.Serialize(true)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", brunning)
	}
}

func Stop(w http.ResponseWriter, r *http.Request) {
	job.Stop()
	fmt.Fprintf(w, "%s", "success")
}

func Start(w http.ResponseWriter, r *http.Request) {
	job.Stop()
	fmt.Fprintf(w, "%s", "success")
}
