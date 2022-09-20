package job

import (
	"crontab/internal/comm"
	"encoding/json"
	"sync"
	"time"
)

/*
*  任务列表管理（添加，删除，更新）
 */

type Job struct {
	Time    string   `json:"time"`    //任务执行时间
	Cmd     string   `json:"cmd"`     //可执行程序
	Args    []string `json:"args"`    //执行参数
	Out     string   `json:"out"`     //输出文件
	Comment string   `json:"comment"` //任务备注
	Start   string   `json:"start"`   //任务单次执行仅用作状态使用
	comm.Schedule
}

func NewJobs() *JobContainer {
	return &JobContainer{mj: make(map[string]*Job), lk: new(sync.RWMutex)}
}

type JobContainer struct {
	mj map[string]*Job
	lk *sync.RWMutex
}

func (jobs *JobContainer) add(k string, v *Job) {
	jobs.lk.Lock()
	defer jobs.lk.Unlock()
	jobs.mj[k] = v
}

func (jobs *JobContainer) del(k string) {
	jobs.lk.Lock()
	defer jobs.lk.Unlock()
	delete(jobs.mj, k)
}

func (jobs *JobContainer) json() ([]byte, error) {
	jobs.lk.RLock()
	defer jobs.lk.RUnlock()
	return json.Marshal(jobs.mj)
}

func (jobs *JobContainer) getJobs() map[string]*Job {
	jobs.lk.RLock()
	defer jobs.lk.RUnlock()
	return jobs.mj
}

func (jobs *JobContainer) replaceJobs(mj map[string]*Job) {
	jobs.lk.Lock()
	defer jobs.lk.Unlock()
	jobs.mj = mj
}

func (jobs *JobContainer) runJobs() {
	t := time.Now()
	if t.Second() == 0 {
		jobs.lk.Lock()
		defer jobs.lk.Unlock()
		minute := t.Minute()
		hour := t.Hour()
		dom := t.Day()
		month := int(t.Month())
		dow := int(t.Weekday())
		for _, j := range jobs.mj {
			if inArray(j.Minute, minute) &&
				inArray(j.Hour, hour) &&
				inArray(j.Dom, dom) &&
				inArray(j.Month, month) &&
				inArray(j.Dow, dow) {
				go runJob(*j)
			}
		}
	}
}
