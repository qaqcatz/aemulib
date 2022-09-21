package aemulib

import (
	"errors"
	"github.com/qaqcatz/adclib"
	"github.com/qaqcatz/nanoshlib"
	"strconv"
	"strings"
	"time"
)

// AEmu: emulator -avd AvdName -port Port ExtraParam.
type AEmu struct {
	EmulatorPath string 	  // e.g /home/hzy/Android/Sdk/emulator/emulator
	AvdName      string       // avd name. see it from emulator -list-avds or ls ~/.android/avd
	Port         string       // adb's port,  the emulator name during runtime will be emulator-Port
	ExtraParam   string       // emulator parameters other than -avd xxx -port xxx
	adbs         *adclib.AdbS // communicate with the emulator
}

// NewAEmu create an AEmu object.
//
// - emulatorPath, e.g. /home/hzy/Android/Sdk/emulator/emulator
//
// - avdName, avd name. see it from emulator -list-avds or ls ~/.android/avd
//
// - Port, adb's port,  the emulator name during runtime will be emulator-Port
//
// - extraParam, emulator parameters other than -avd xxx
//
// - adbPath, httpForwardPort: adclib.NewAdbS(adbPath, "emulator-"+port, "127.0.0.1", httpForwardPort).
func NewAEmu(emulatorPath string, avdName string, port string, extraParam string,
	adbPath string, httpForwardPort string) *AEmu {
	adbs := adclib.NewAdbS(adbPath, "emulator-"+port, "127.0.0.1", httpForwardPort)
	return &AEmu{
		EmulatorPath: emulatorPath,
		AvdName: avdName,
		Port: port,
		ExtraParam: extraParam,
		adbs: adbs,
	}
}

// Exec encapsulates *(adclib.AdbS).Exec().
func (aemu *AEmu) Exec(cmdStr string, timoutMs int) ([]byte, []byte, error, bool) {
	return aemu.adbs.Exec(cmdStr, timoutMs)
}

// GetPid: ps -ef | grep -e /emulator/qemu/.*AvdName | grep -v 'grep -e /emulator/qemu/.*AvdName' | awk '{print $2}'
//  example:
//	hzy        86288   53101 38 19:50 pts/1    00:02:01 /home/hzy/Android/Sdk/emulator/qemu/linux-x86_64/qemu-system-x86_64 -avd test
// 	hzy        87319   85954  0 19:55 pts/4    00:00:00 grep -e /emulator/qemu/.*test
//	output:
//  86288
// Return (pid, nil) or ("", error: ps -ef error or pid is not number).
// May wait 1s.
func (aemu *AEmu) GetPid() (string, error) {
	outStream, _, err := nanoshlib.Exec("ps -ef | grep -e /emulator/qemu/.*"+aemu.AvdName+
		" | grep -v 'grep -e /emulator/qemu/.*"+aemu.AvdName+"' | awk '{print $2}'", 1000)
	if err != nil {
		return "", err
	}
	pid := strings.TrimSpace(string(outStream))
	_, err = strconv.Atoi(pid)
	if err != nil {
		return "", errors.New("pid is not number: " + err.Error())
	}
	return pid, nil
}

// Kill: kill -9 GetPid().
// May wait 1s.
func (aemu *AEmu) Kill() error {
	pid, err := aemu.GetPid()
	if err != nil {
		return err
	}
	_, _, err = nanoshlib.Exec("kill -9 "+pid, 1000)
	if err != nil {
		return err
	}
	return nil
}

// hasResumedActivity: return 'adb -s xxx shell dumpsys activity activities | grep mResumedActivity' not (error && empty).
// May wait 1s
func (aemu *AEmu) hasResumedActivity() bool {
	outStream, _, err, _ := aemu.Exec("shell dumpsys activity activities | grep mResumedActivity", 1000)
	if err != nil || strings.TrimSpace(string(outStream)) == "" {
		return false
	}
	return true
}

// Restart:
//
// 1. Kill(), ignore error.
//
// 2. emulator -avd AvdName -port Port ExtraParam;
//
// 3. wait for any resumed activity or start error or timeout
//
// 4. sleep 3s.
//
// May wait timeoutMs + 3s.
func (aemu *AEmu) Restart(timeoutMs int) error {
	// 1. kill
	_ = aemu.Kill()
	// 2. emulator -avd AvdName -port Port ExtraParam
	doneChan, _, err := nanoshlib.Exec0s(aemu.EmulatorPath+" -avd "+aemu.AvdName+" -port "+aemu.Port+" "+aemu.ExtraParam)
	if err != nil {
		return errors.New("emulator restart error: " + err.Error())
	}
	// 3. wait for any resumed activity or start error or timeout
	timeout := time.After(time.Duration(timeoutMs) * time.Millisecond)
	hra := make(chan bool)
	ok := false
	err = nil
	for {
		go func() { time.Sleep(3*time.Second); hra <- aemu.hasResumedActivity() }()
		select {
		case err_ := <-doneChan:
			err = err_
		case <-timeout:
			err = errors.New("timeout")
		case res := <-hra:
			ok = res
		}
		if err != nil {
			return errors.New("emulator restart error: " + err.Error())
		}
		_, err = aemu.GetPid()
		if err != nil {
			return errors.New("emulator restart error: can not find the process: " + err.Error())
		}
		if ok {
			break
		}
	}
	// 4. sleep 3s
	// After startup, some settings-related activities may be run first, so wait for a while
	time.Sleep(3*time.Second)
	return nil
}

// RestartE is an extension of Restart, it will restart the emulator with extraParam
func (aemu *AEmu) RestartE(extraParam string, timeoutMs int) error {
	// 1. kill
	_ = aemu.Kill()
	// 2. emulator -avd AvdName -port Port ExtraParam
	doneChan, _, err := nanoshlib.Exec0s(aemu.EmulatorPath+" -avd "+aemu.AvdName+" -port "+aemu.Port+" "+extraParam)
	if err != nil {
		return errors.New("emulator restart error: " + err.Error())
	}
	// 3. wait for any resumed activity or start error or timeout
	timeout := time.After(time.Duration(timeoutMs) * time.Millisecond)
	hra := make(chan bool)
	ok := false
	err = nil
	for {
		go func() { time.Sleep(3*time.Second); hra <- aemu.hasResumedActivity() }()
		select {
		case err_ := <-doneChan:
			err = err_
		case <-timeout:
			err = errors.New("timeout")
		case res := <-hra:
			ok = res
		}
		if err != nil {
			return err
		}
		_, err = aemu.GetPid()
		if err != nil {
			return errors.New("emulator restart error: can not find the process: " + err.Error())
		}
		if ok {
			break
		}
	}
	// 4. sleep 3s
	// After startup, some settings-related activities may be run first, so wait for a while
	time.Sleep(3*time.Second)
	return nil
}
