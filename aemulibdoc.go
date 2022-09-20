//aemulib is responsible for managing the startup and shutdown of the android emulator,
//as well as communicating with the emulator.
//
//adclib(https://github.com/qaqcatz/adclib) provide stable communication interfaces to android device.
//However, we still need to manually open the android device, and restart the device when something wrong.
//
//Management of remote devices is difficult, but at least we can automatically manage local emulators.
//
//That's what aemulib does.
//
//For more convenience, aemulib also encapsulates adclib to communicate with the emulator.
//
//Make sure there is only one emulator instance for an avd name!
package aemulib
