package vmx

import "testing"

func TestParsingTag(t *testing.T) {
	tests := []struct {
		tag       string
		name      string
		omitempty bool
		err       string
	}{
		{"vmx:displayname", "", false, "Tag name has to be enclosed in double quotes: vmx:displayname"},
		{"vmx:", "", false, "Invalid tag: vmx:"},
		{`vmx:""`, "", false, `Tag name is missing: vmx:""`},
		{"vm", "", false, "Invalid tag: vm"},
		{`vmx:"displayname,omitempty`, "displayname", true, ""},
		{`vmx:"displayname,blah"`, "displayname", false, ""},
		{`vmx:"-"`, "-", false, ""},
	}

	for _, tt := range tests {
		name, omitempty, err := parseTag(tt.tag)
		equals(t, tt.name, name)
		equals(t, tt.omitempty, omitempty)
		if err != nil {
			equals(t, tt.err, err.Error())
		} else {
			equals(t, tt.err, "")
		}
	}
}
