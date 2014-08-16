package vmx

import (
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestMarshal(t *testing.Test) {
	type VM struct {
		Hwversion    string `vmx:virtualHW.version`
		HwProdCompat string `vmx:virtualHW.productCompatibility`
		Memsize      string `vmx:memsize`
		Numvcpus	 string `vmx:numvcpus`
		MemHotAdd    bool   `vmx:mem.hotadd`
		DisplayName  string `vmx:displayName`
		GuestOS      string `vmx:guestOS`
		Autoanswer   bool   `vmx:msg.autoAnswer`
	}

	var vm VM
	data, err := Marshal(&vm)
	ok(t, err)
}

func TestMarshalEmbedded(t *testing.Test) {
	type Vhardware struct {
		version string `vmx:version`
		compat  string `vmx:productCompatibility`
	}

	type VM struct {
		Hwversion   Vhardware `vmx:virtualHW`
		Memsize     string    `vmx:memsize`
		MemHotAdd   bool      `vmx:mem.hotadd`
		DisplayName string    `vmx:displayName`
		GuestOS     string    `vmx:guestOS`
	}

	var vm VM
	data, err := Marshal(&vm)
	ok(t, err)
}

func TestMarshalArray(t *testing.Test) {
	type Vhardware struct {
		version string `vmx:version`
		compat  string `vmx:productCompatibility`
	}

	type Ethernet struct {
		Present              bool   `vmx:present`
		ConnectionType       string `vmx:connectionType`
		VirtualDev           string `vmx:virtualDev`
		WakeOnPcktRcv        bool   `vmx:wakeOnPcktRcv`
		AddressType          string `vmx:addressType`
		LinkStatePropagation bool   `vmx:linkStatePropagation.enable`
	}

	type VM struct {
		Hwversion   Vhardware  `vmx:virtualHW`
		Memsize     string     `vmx:memsize`
		MemHotAdd   bool       `vmx:mem.hotadd`
		DisplayName string     `vmx:displayName`
		GuestOS     string     `vmx:guestOS`
		Autoanswer  bool       `vmx:msg.autoAnswer`
		Ethernet    []Ethernet `vmx:ethernet`
	}

	var vm VM
	data, err := Marshal(&vm)
	ok(t, err)
}

func TestMarshalUSB(t *testing.T) {

}

func TestMarshalIDEDevices(t *testing.T) {
	// Maximum two IDE ports: primary and secondary
	// Maximum two devices per port
	// Only same device type per port
	// ide<port>:<id>.present
}

func TestMarshalParallelPorts(t *testing.T) {
	// max 3
}
