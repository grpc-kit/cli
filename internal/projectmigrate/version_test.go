package projectmigrate

import "testing"

func TestValidateTargetCLIVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "without v prefix", input: "0.4.0", want: "0.4.0"},
		{name: "with v prefix", input: "v0.4.0", want: "0.4.0"},
		{name: "empty", input: "", wantErr: true},
		{name: "development default", input: "0.0.0", wantErr: true},
		{name: "prerelease", input: "0.3.9-beta.1", wantErr: true},
		{name: "non canonical", input: "0.4", wantErr: true},
		{name: "build metadata", input: "0.4.0+local", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ValidateTargetCLIVersion(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("ValidateTargetCLIVersion(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("ValidateTargetCLIVersion(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}
