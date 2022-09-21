# aemulib
[![Go Reference](README.assets/aemulib.svg)](https://pkg.go.dev/github.com/qaqcatz/aemulib)

aemulib is responsible for managing the startup and shutdown of the android emulator, as well as communicating with the emulator.

adclib (https://github.com/qaqcatz/adclib) provide stable communication interfaces to android device.
However, we still need to manually open the android device, and restart the device when something wrong.

Management of remote devices is difficult, but at least we can automatically manage local emulators.

That's what aemulib does. 

For more convenience, aemulib also encapsulates adclib to communicate with the emulator.

**Make sure there is only one emulator instance for an avd name!**

# How to use

## Import

```golang
// go.mod
require github.com/qaqcatz/aemulib v1.1.0
// xxx.go
import "github.com/qaqcatz/aemulib"
```

## AEmu

```golang
// AEmu: emulator -avd AvdName -port Port ExtraParam
type AEmu struct {
	EmulatorPath string 	  // e.g /home/hzy/Android/Sdk/emulator/emulator
	AvdName      string       // avd name. see it from emulator -list-avds or ls ~/.android/avd
	Port         string       // adb's port,  the emulator name during runtime will be emulator-Port
	ExtraParam   string       // emulator parameters other than -avd xxx -port xxx
	adbs         *adclib.AdbS // communicate with the emulator
}
```

### NewAEmu

```golang
func NewAEmu(emulatorPath string, avdName string, port string, extraParam string, adbPath string, httpForwardPort string) *AEmu
```


NewAEmu create an AEmu object.

- emulatorPath, e.g. /home/hzy/Android/Sdk/emulator/emulator
- avdName, avd name. see it from emulator -list-avds or ls ~/.android/avd
- Port, adb's port, the emulator name during runtime will be emulator-Port
- extraParam, emulator parameters other than -avd xxx
- adbPath, httpForwardPort: adclib.NewAdbS(adbPath, "emulator-"+port, "127.0.0.1", httpForwardPort).

### Exec

```golang
func (aemu *AEmu) Exec(cmdStr string, timoutMs int) ([]byte, []byte, error, bool)
```

Exec encapsulates *(adclib.AdbS).Exec().

### Restart

func (aemu *AEmu) Restart(timeoutMs int) error
Restart:
1. Kill(), ignore error.
2. emulator -avd AvdName -port Port ExtraParam;
3. wait for any resumed activity or start error or timeout
4. sleep 3s.

May wait timeoutMs + 3s.

```golang
func ExampleAEmu_Restart() {
	// emulator -avd test first
	aemu := NewAEmu("/home/hzy/Android/Sdk/emulator/emulator",
		"test", "5556", "-wipe-data",
		"/home/hzy/Android/Sdk/platform-tools/adb", "0")
	err := aemu.Restart(90000)
	if err != nil {
		log.Fatal("error: " + err.Error())
	} else {
		log.Println("ok!")
	}
}
```

### RestartE

```golang
func (aemu *AEmu) RestartE(extraParam string, timeoutMs int) error
```

RestartE is an extension of Restart, it will restart the emulator with extraParam

```golang
func ExampleAEmu_RestartE() {
	// emulator -avd test first
	aemu := NewAEmu("/home/hzy/Android/Sdk/emulator/emulator",
		"test", "5556", "-wipe-data",
		"/home/hzy/Android/Sdk/platform-tools/adb", "0")
	err := aemu.RestartE("",90000)
	if err != nil {
		log.Fatal("error: " + err.Error())
	} else {
		log.Println("ok!")
	}
}
```

### Kill

```golang
func (aemu *AEmu) Kill() error
```


Kill: kill -9 GetPid(). May wait 1s.

### GetPid

```golang
func (aemu *AEmu) GetPid() (string, error)
```

GetPid: ps -ef | grep -e /emulator/qemu/.* AvdName | grep -v 'grep -e /emulator/qemu/.*AvdName' | awk '{print $2}'*

```shell
example:
 hzy        86288   53101 38 19:50 pts/1    00:02:01 /home/hzy/Android/Sdk/emulator/qemu/linux-x86_64/qemu-system-x86_64 -avd test
 hzy        87319   85954  0 19:55 pts/4    00:00:00 grep -e /emulator/qemu/.*test
 output:
 86288
```

Return (pid, nil) or ("", error: ps -ef error or pid is not number). May wait 1s.

