// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

//日志功能模块。
package log

//初始化: 使用configBase.ConfigBase对象创建一个写日志的对象
//启用的各项参数分别为:
//   nbs.logfile =>日志文件
//   nbs.audit 启用审计模式
//   nbs.level 设定日志级别
//   nbs.logsize 日志滚动 size
//   nbs.logcount 保留日志个数
//日志对象使用结束，需要调用Release方法销毁日志对象，确保日志完整写入

import "os"
import "github.com/TingYunAPM/go/utils/pool"
import "time"
import "github.com/TingYunAPM/go/utils/cache_config"
import "github.com/TingYunAPM/go/utils/service"
import "fmt"
import "strings"
import "archive/tar"
import "compress/gzip"
import "io"

type Logger struct {
	messagePool pool.SerialReadPool
	svc         service.Service
	fp          *os.File
	pid         int
	configs     *cache_config.Configuration
	topSize     int64
	maxLogSize  int64
}

const (
	LevelOff      = 0x00
	LevelCritical = 0x01
	LevelError    = 0x02
	LevelWarning  = 0x03
	LevelInfo     = 0x04
	LevelVerbos   = 0x05
	LevelDebug    = 0x06
	LevelMask     = 0x07
)
const ConfigStringNBSLevel = 4
const ConfigStringNBSLogFileName = 5
const ConfigBoolNBSAudit = 3
const ConfigIntegerNBSMaxLogSize = 3
const ConfigIntegerNBSMaxLogCount = 4
const Audit = 0x08

//创建一个写日志对象,使用 configBase.ConfigBase对象初始化
func New(conf *cache_config.Configuration) *Logger {
	return new(Logger).init(conf)
}

//格式化字符串输出,功能参考fmt.Printf
func (l *Logger) Printf(level int, format string, a ...interface{}) {
	if l == nil || !l.levelEnable(level) {
		return
	}
	l.Append(fmt.Sprintf("%s (pid:%d) %s :", time.Now().Format("2006-01-02 15:04:05.000"), l.pid, levelMap[level&LevelMask]) + fmt.Sprintf(format, a...))
}

//功能参考fmt.Print
func (l *Logger) Print(level int, a ...interface{}) {
	if l == nil || !l.levelEnable(level) {
		return
	}
	l.Append(fmt.Sprintf("%s (pid:%d) %s :", time.Now().Format("2006-01-02 15:04:05.000"), l.pid, levelMap[level&LevelMask]) + fmt.Sprint(a...))
}

//功能参考fmt.Println
func (l *Logger) Println(level int, a ...interface{}) {
	if l == nil || !l.levelEnable(level) {
		return
	}
	l.Append(fmt.Sprintf("%s (pid:%d) %s :", time.Now().Format("2006-01-02 15:04:05.000"), l.pid, levelMap[level&LevelMask]) + fmt.Sprintln(a...))
}

//追加文本，不打印时间pid等信息
func (l *Logger) Append(data string) {
	if l == nil {
		return
	}
	//	fmt.Println("append message:", data)
	l.messagePool.Put(&message{data})
}

//停止日志对象，不再写入
func (l *Logger) Release() {
	if l == nil {
		return
	}
	l.svc.Stop()
}
func (l *Logger) levelEnable(level int) bool {
	if l.configLevel() < level&LevelMask || level&LevelMask < LevelCritical { //log level filter
		return false
	}
	audit := l.configs.CBools.Read(ConfigBoolNBSAudit, false)
	if level&Audit != 0 && !audit { //audit filter
		return false
	}
	return true
}

func compress(sfile string, tfile string) bool {
	fw, err := os.OpenFile(tfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	defer fw.Close()
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()
	fr, err := os.Open(sfile)
	defer fr.Close()
	state, err := fr.Stat()
	head, err := tar.FileInfoHeader(state, "")
	err = tw.WriteHeader(head)
	if err != nil {
		return false
	}
	_, err = io.Copy(tw, fr)
	if err != nil {
		return false
	}
	return true
}
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func rollFile(filename string, num int) {
	getName := func(num int) string {
		if num == 0 {
			return filename
		}
		return fmt.Sprintf("%s.%d", filename, num)
	}
	top := getName(num)
	if fileExist(top) {
		os.Remove(top)
	}
	for i := num; i > 0; i-- {
		next := getName(i - 1)
		if fileExist(next) {
			os.Rename(next, top)
		}
		top = next
	}
}
func rollAndCompress(filename string, num int) {
	gzFile := filename + ".tar.gz"
	rollFile(gzFile, num)
	if compress(filename, gzFile) {
		os.Remove(filename)
	} else {
		rollFile(filename, num)
	}
}

type message struct {
	data string
}

var levelMap = [7]string{"OFF", "CRITICAL", "ERROR", "WARNING", "INFO", "VERBOSE", "DEBUG"}

var dict = map[string]int{"off": 0, "critical": 1, "error": 2, "warning": 3, "info": 4, "verbose": 5, "debug": 6}

func (l *Logger) configLevel() int {
	levelString := l.configs.CStrings.Read(ConfigStringNBSLevel, "info")
	ret, found := dict[strings.ToLower(levelString)]
	if found {
		return ret
	}
	return LevelInfo
}
func (l *Logger) init(conf *cache_config.Configuration) *Logger {
	l.configs = conf
	l.messagePool.Init()
	l.fp = nil
	l.pid = os.Getpid()
	l.topSize = -1
	l.updateLogSize()
	l.svc.Start(l.loop)
	return l
}

func (l *Logger) write(msg string) {
	if l.fp == nil {
		return
	}
	n, err := l.fp.WriteString(msg)
	if err == nil {
		l.topSize += int64(n)
	}
}
func (l *Logger) openFile() {
	if l.fp != nil {
		return
	}
	logfile := l.configs.CStrings.Read(ConfigStringNBSLogFileName, "")
	if logfile == "" {
		return
	}
	//	fmt.Println("open log file:", logfile)
	fp, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		l.fp = fp
		stat, err := fp.Stat()
		if err == nil {
			l.topSize = stat.Size()
		}
	} else {
		l.fp = nil
		fmt.Print(err)
	}
}
func (l *Logger) updateLogSize() {
	l.maxLogSize = int64(l.configs.CIntegers.Read(ConfigIntegerNBSMaxLogSize, 10)) * 1024 * 1024
	//	fmt.Println("maxLogSize=", l.maxLogSize)
}
func (l *Logger) closeFile() {
	if l.fp != nil {
		l.fp.Close()
		l.fp = nil
	}
}
func (l *Logger) processMessage() {
	updated := false
	for l.messagePool.Size() > 0 {
		if !updated {
			l.updateLogSize()
		}
		for msg := l.messagePool.Get(); msg != nil; msg = l.messagePool.Get() {
			l.openFile()
			l.write(msg.(*message).data)
			if l.maxLogSize <= l.topSize {
				l.closeFile()
				l.logsShift()
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	l.closeFile()
}
func (l *Logger) logsShift() {
	logfile := l.configs.CStrings.Read(ConfigStringNBSLogFileName, "")
	state, err := os.Stat(logfile)
	if err != nil {
		return
	}
	l.updateLogSize()
	if state.Size() < l.maxLogSize {
		return
	}
	logCount := l.configs.CIntegers.Read(ConfigIntegerNBSMaxLogCount, 3)
	rollAndCompress(logfile, int(logCount))

}
func (l *Logger) loop(running func() bool) {
	for running() {
		l.processMessage()
		time.Sleep(10 * time.Millisecond)
	}
	l.processMessage()
}
