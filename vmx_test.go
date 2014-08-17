package vmx

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
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

func TestMarshal(t *testing.T) {
	type VM struct {
		Encoding     string `vmx:".encoding"`
		Annotation   string `vmx:"annotation"`
		Hwversion    uint8  `vmx:"virtualHW.version"`
		HwProdCompat string `vmx:"virtualHW.productCompatibility"`
		Memsize      uint   `vmx:"memsize"`
		Numvcpus     uint   `vmx:"numvcpus"`
		MemHotAdd    bool   `vmx:"mem.hotadd"`
		DisplayName  string `vmx:"displayName"`
		GuestOS      string `vmx:"guestOS"`
		Autoanswer   bool   `vmx:"msg.autoAnswer"`
	}

	vm := new(VM)
	vm.Encoding = "utf-8"
	vm.Annotation = "Test VM"
	vm.Hwversion = 10
	vm.HwProdCompat = "hosted"
	vm.Memsize = 1024
	vm.Numvcpus = 2
	vm.MemHotAdd = false
	vm.DisplayName = "test"
	vm.GuestOS = "other3xlinux-64"
	vm.Autoanswer = true

	data, err := Marshal(vm)
	ok(t, err)
	expected := `.encoding = "utf-8"
annotation = "Test VM"
virtualHW.version = "10"
virtualHW.productCompatibility = "hosted"
memsize = "1024"
numvcpus = "2"
mem.hotadd = "false"
displayName = "test"
guestOS = "other3xlinux-64"
msg.autoAnswer = "true"
`
	equals(t, expected, string(data))
}

func TestMarshalEmbedded(t *testing.T) {
	type Vhardware struct {
		Version string `vmx:"version"`
		Compat  string `vmx:"productCompatibility"`
	}

	type VM struct {
		Encoding    string    `vmx:".encoding"`
		Annotation  string    `vmx:"annotation"`
		Vhardware   Vhardware `vmx:"virtualHW"`
		Memsize     uint      `vmx:"memsize"`
		Numvcpus    uint      `vmx:"numvcpus"`
		MemHotAdd   bool      `vmx:"mem.hotadd"`
		DisplayName string    `vmx:"displayName"`
		GuestOS     string    `vmx:"guestOS"`
		Autoanswer  bool      `vmx:"msg.autoAnswer"`
	}

	vm := new(VM)
	vm.Encoding = "utf-8"
	vm.Annotation = "Test VM"
	vm.Vhardware = Vhardware{
		Version: "10",
		Compat:  "hosted",
	}
	vm.Memsize = 1024
	vm.Numvcpus = 2
	vm.MemHotAdd = false
	vm.DisplayName = "test"
	vm.GuestOS = "other3xlinux-64"
	vm.Autoanswer = true

	data, err := Marshal(vm)
	ok(t, err)
	expected := `.encoding = "utf-8"
annotation = "Test VM"
virtualHW.version = "10"
virtualHW.productCompatibility = "hosted"
memsize = "1024"
numvcpus = "2"
mem.hotadd = "false"
displayName = "test"
guestOS = "other3xlinux-64"
msg.autoAnswer = "true"
`
	equals(t, expected, string(data))
}

func TestMarshalArray(t *testing.T) {
	type Vhardware struct {
		Version string `vmx:"version"`
		Compat  string `vmx:"productCompatibility"`
	}

	type Ethernet struct {
		StartConnected       bool   `vmx:"startConnected"`
		Present              bool   `vmx:"present"`
		ConnectionType       string `vmx:"connectionType"`
		VirtualDev           string `vmx:"virtualDev"`
		WakeOnPcktRcv        bool   `vmx:"wakeOnPcktRcv"`
		AddressType          string `vmx:"addressType"`
		LinkStatePropagation bool   `vmx:"linkStatePropagation.enable,omitempty"`
	}

	type VM struct {
		Encoding    string     `vmx:".encoding"`
		Annotation  string     `vmx:"annotation"`
		Vhardware   Vhardware  `vmx:"virtualHW"`
		Memsize     uint       `vmx:"memsize"`
		Numvcpus    uint       `vmx:"numvcpus"`
		MemHotAdd   bool       `vmx:"mem.hotadd"`
		DisplayName string     `vmx:"displayName"`
		GuestOS     string     `vmx:"guestOS"`
		Autoanswer  bool       `vmx:"msg.autoAnswer"`
		Ethernet    []Ethernet `vmx:"ethernet"`
	}

	vm := new(VM)
	vm.Encoding = "utf-8"
	vm.Annotation = "Test VM"
	vm.Vhardware = Vhardware{
		Version: "9",
		Compat:  "hosted",
	}
	vm.Ethernet = []Ethernet{
		{
			StartConnected:       true,
			Present:              true,
			ConnectionType:       "bridged",
			VirtualDev:           "e1000",
			WakeOnPcktRcv:        false,
			AddressType:          "generated",
			LinkStatePropagation: true,
		},
		{
			StartConnected: true,
			Present:        true,
			ConnectionType: "nat",
			VirtualDev:     "e1000",
			WakeOnPcktRcv:  false,
			AddressType:    "generated",
		},
	}
	vm.Memsize = 1024
	vm.Numvcpus = 2
	vm.MemHotAdd = false
	vm.DisplayName = "test"
	vm.GuestOS = "other3xlinux-64"
	vm.Autoanswer = true

	data, err := Marshal(vm)
	ok(t, err)
	expected := `.encoding = "utf-8"
annotation = "Test VM"
virtualHW.version = "9"
virtualHW.productCompatibility = "hosted"
memsize = "1024"
numvcpus = "2"
mem.hotadd = "false"
displayName = "test"
guestOS = "other3xlinux-64"
msg.autoAnswer = "true"
ethernet0.startConnected = "true"
ethernet0.present = "true"
ethernet0.connectionType = "bridged"
ethernet0.virtualDev = "e1000"
ethernet0.wakeOnPcktRcv = "false"
ethernet0.addressType = "generated"
ethernet0.linkStatePropagation.enable = "true"
ethernet1.startConnected = "true"
ethernet1.present = "true"
ethernet1.connectionType = "nat"
ethernet1.virtualDev = "e1000"
ethernet1.wakeOnPcktRcv = "false"
ethernet1.addressType = "generated"
`
	equals(t, expected, string(data))
}
