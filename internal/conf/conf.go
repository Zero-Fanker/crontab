package conf

import (
	"bufio"
	"crontab/internal/job"
	"crontab/internal/logger"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
* 任务配置文件，读取&更新
 */

var confFilePath *string

func Load(conf *string) (bool, error) {
	if conf != nil {
		confFilePath = conf
	}
	logger.SysWriteLn("Load config start ...")
	fp, err := os.Open(*confFilePath)
	if err != nil {
		logger.SysPrintf("Err %s .\n", err)
		return false, err
	}
	defer fp.Close()
	rd := bufio.NewReader(fp)

	for {
		line, rdErr := rd.ReadString('\n')

		if rdErr != nil && rdErr != io.EOF {
			logger.SysPrintf("Err %s.\n", rdErr)
			return false, rdErr
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if rdErr == io.EOF {
				break
			}
			continue
		}
		successful, err := job.Add(line)
		if !successful || err != nil {
			return false, err
		}
	}

	logger.SysWriteLn("Load config end.")
	return true, nil
}

func Save() (bool, error) {
	logger.SysWriteLn("Flush config start ...")
	fp, err := os.Create(*confFilePath)
	if err != nil {
		logger.SysWriteLn(err)
		return false, err
	}
	defer fp.Close()
	tjobs := configJobs.getJobs()
	for _, j := range tjobs {
		b, _ := json.Marshal(j)
		fmt.Fprintf(fp, "%s\n", b)
	}
	logger.SysWriteLn("Flush config end.")
	return true, nil
}
