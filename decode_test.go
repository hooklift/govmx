package vmx

import (
	"fmt"
	"testing"
)

var data = `.encoding = "UTF-8"
annotation = "Terraform VMWARE VIX test"
bios.bootorder = "hdd,CDROM"
checkpoint.vmstate = ""
cleanshutdown = "TRUE"
config.version = "8"
cpuid.corespersocket = "1"
displayname = "core01"
ehci.pcislotnumber = "-1"
ehci.present = "FALSE"
ethernet1.addressType = "generated"
ethernet1.connectionType = "bridged"
ethernet1.linkStatePropagation.enable = "true"
ethernet1.present = "TRUE"
ethernet1.startConnected = "true"
ethernet1.virtualDev = "e1000"
ethernet1.wakeOnPcktRcv = "false"
ethernet2.address = "00:50:56:aa:bb:cc"
ethernet2.addressType = "static"
ethernet2.connectionType = "nat"
ethernet2.present = "TRUE"
ethernet2.startConnected = "true"
ethernet2.virtualDev = "e1000"
ethernet2.wakeOnPcktRcv = "false"
ethernet3.addressType = "generated"
ethernet3.connectionType = "hostonly"
ethernet3.present = "TRUE"
ethernet3.startConnected = "true"
ethernet3.virtualDev = "e1000"
ethernet3.wakeOnPcktRcv = "false"
extendedconfigfile = "core01.vmxf"
floppy0.present = "FALSE"
guestos = "other3xlinux-64"
gui.fullscreenatpoweron = "FALSE"
gui.viewmodeatpoweron = "windowed"
hgfs.linkrootshare = "TRUE"
hgfs.maprootshare = "TRUE"
ide1:0.devicetype = "cdrom-image"
ide1:0.filename = "/Users/camilo/Dropbox/Development/cloudescape/dobby-boxes/coreos/packer_cache/e159f7e70f4ccc346ee76b2a32cbdf059549a7ca82e91edbeed5d747bcdd50f9.iso"
ide1:0.present = "TRUE"
isolation.tools.hgfs.disable = "FALSE"
memsize = "1024"
monitor.phys_bits_used = "40"
msg.autoanswer = "true"
numvcpus = "1"
nvram = "core01.nvram"
pcibridge0.pcislotnumber = "17"
pcibridge0.present = "TRUE"
pcibridge4.functions = "8"
pcibridge4.pcislotnumber = "21"
pcibridge4.present = "TRUE"
pcibridge4.virtualdev = "pcieRootPort"
pcibridge5.functions = "8"
pcibridge5.pcislotnumber = "22"
pcibridge5.present = "TRUE"
pcibridge5.virtualdev = "pcieRootPort"
pcibridge6.functions = "8"
pcibridge6.pcislotnumber = "23"
pcibridge6.present = "TRUE"
pcibridge6.virtualdev = "pcieRootPort"
pcibridge7.functions = "8"
pcibridge7.pcislotnumber = "24"
pcibridge7.present = "TRUE"
pcibridge7.virtualdev = "pcieRootPort"
policy.vm.mvmtid = ""
powertype.poweroff = "soft"
powertype.poweron = "soft"
powertype.reset = "soft"
powertype.suspend = "soft"
proxyapps.publishtohost = "FALSE"
remotedisplay.vnc.enabled = "TRUE"
remotedisplay.vnc.port = "5919"
replay.filename = ""
replay.supported = "FALSE"
scsi0.pcislotnumber = "16"
scsi0.present = "TRUE"
scsi0.virtualdev = "lsilogic"
scsi0:0.filename = "disk-cl1.vmdk"
scsi0:0.present = "TRUE"
scsi0:0.redo = ""
softPowerOff = "FALSE"
sound.startconnected = "FALSE"
tools.synctime = "TRUE"
tools.upgrade.policy = "upgradeAtPowerCycle"
usb.pcislotnumber = "-1"
usb.present = "FALSE"
uuid.action = "create"
uuid.bios = "56 4d 59 1a 1a 9b 5f d8-29 6c 70 d0 bf 20 41 99"
uuid.location = "56 4d 59 1a 1a 9b 5f d8-29 6c 70 d0 bf 20 41 99"
vc.uuid = ""
virtualhw.productcompatibility = "hosted"
virtualhw.version = "9"
vmci0.id = "1861462627"
vmci0.pcislotnumber = "35"
vmci0.present = "TRUE"
vmotion.checkpointfbsize = "67108864"
ethernet1.pciSlotNumber = "32"
ethernet2.pciSlotNumber = "33"
ethernet3.pciSlotNumber = "34"
ethernet1.generatedAddress = "00:0c:29:20:41:a3"
ethernet1.generatedAddressOffset = "10"
ethernet3.generatedAddress = "00:0c:29:20:41:b7"
ethernet3.generatedAddressOffset = "30"
`

func TestUnmarshal(t *testing.T) {
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
		Vhardware   Vhardware  `vmx:"virtualhw"`
		Memsize     uint       `vmx:"memsize"`
		Numvcpus    uint       `vmx:"numvcpus"`
		MemHotAdd   bool       `vmx:"mem.hotadd"`
		DisplayName string     `vmx:"displayName"`
		GuestOS     string     `vmx:"guestOS"`
		Autoanswer  bool       `vmx:"msg.autoAnswer"`
		Ethernet    []Ethernet `vmx:"ethernet"`
	}

	vm := new(VM)
	err := Unmarshal([]byte(data), vm)
	ok(t, err)
	fmt.Printf("%+v\n", vm)
}
