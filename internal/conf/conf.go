package conf

import (
	"crontab/internal/job"
	"crontab/internal/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/*
* 任务配置文件，读取&更新
 */

type ConfStruct struct {
	jobs []job.Job
}

var confFilePath *string

func Load(conf *string) (bool, error) {
	if conf != nil {
		confFilePath = conf
	}
	logger.SysWriteLn("Load config start ...")

	content, err := ioutil.ReadFile(*confFilePath)
	if err != nil {
		logger.SysPrintf("Err %s .\n", err)
		return false, err
	}
	allJobs := new([]job.Job)
	err = json.Unmarshal(content, &allJobs)
	if err != nil {
		logger.SysPrintf("Err %s .\n", err)
		return false, err
	}

	successful, err := job.AddAll(allJobs)
	if err != nil {
		logger.SysPrintf("Err %s .\n", err)
		return false, err
	}

	if !successful {
		logger.SysPrintf("Undefined error")
		return false, fmt.Errorf("Undefined error")
	}

	logger.SysWriteLn("Load config end.")
	return true, nil
}

func Save() (bool, error) {
	logger.SysWriteLn("Saving config ...")
	fp, err := os.Create(*confFilePath)
	if err != nil {
		logger.SysWriteLn(err)
		return false, err
	}
	defer fp.Close()
	ret, err := job.Serialize(false)
	if err != nil {
		logger.SysWriteLn(err)
		return false, err
	}
	_, err = fmt.Fprint(fp, ret)
	if err != nil {
		logger.SysWriteLn(err)
		return false, err
	}
	logger.SysWriteLn("Config saved.")
	return true, nil
}
