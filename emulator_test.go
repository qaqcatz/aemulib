package aemulib

import (
	"log"
	"testing"
)

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

// emulator -avd test first
func TestAEmu_Restart(t *testing.T) {
	aemu := NewAEmu("/home/hzy/Android/Sdk/emulator/emulator",
		"test", "5556", "-wipe-data",
		"/home/hzy/Android/Sdk/platform-tools/adb", "0")
	err := aemu.Restart(90000)
	if err != nil {
		t.Fatal("error: " + err.Error())
	} else {
		t.Log("ok!")
	}
}

// emulator -avd test first
func TestAEmu_RestartE(t *testing.T) {
	aemu := NewAEmu("/home/hzy/Android/Sdk/emulator/emulator",
		"test", "5556", "-wipe-data",
		"/home/hzy/Android/Sdk/platform-tools/adb", "0")
	err := aemu.RestartE("",90000)
	if err != nil {
		t.Fatal("error: " + err.Error())
	} else {
		t.Log("ok!")
	}
}