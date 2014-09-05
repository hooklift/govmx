# GoVMX
[![GoDoc](https://godoc.org/github.com/cloudescape/govmx?status.svg)](https://godoc.org/github.com/cloudescape/govmx)
[![Build Status](https://travis-ci.org/cloudescape/govmx.svg?branch=master)](https://travis-ci.org/cloudesc        ape/govmx)

Data encoding and decoding library for VMware VMX files.


## Encoding

```go
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
```

`data` should be: 

```
.encoding = "utf-8"
annotation = "Test VM"
virtualHW.version = "10"
virtualHW.productCompatibility = "hosted"
memsize = "1024"
numvcpus = "2"
mem.hotadd = "false"
displayName = "test"
guestOS = "other3xlinux-64"
msg.autoAnswer = "true"
```

## Decoding

```go
var data = `.encoding = "utf-8"
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
err := Unmarshal([]byte(data), vm)
```

# Resources
* **The Laws of Reflection:** http://blog.golang.org/laws-of-reflection

# License
Copyright 2014 Cloudescape

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
