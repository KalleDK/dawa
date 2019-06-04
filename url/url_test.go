package url

import (
	"fmt"
	"testing"
)

func TestUrl_UnmarshalJSON(t *testing.T) {

	var json_input = `"https://github.com/kalledk/dawa"`

	fmt.Print(json_input)
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		u       *URL
		args    args
		wantErr bool
		expect  string
	}{
		{
			"Empty",
			&URL{},
			args{
				[]byte(json_input),
			},
			false,
			"https://github.com/kalledk/dawa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Url.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.u.String()
			if (tt.expect != got) != tt.wantErr {
				t.Errorf("Url.UnmarshalJSON() expected = %v, got = %v, wantErr %v", tt.expect, got, tt.wantErr)
			}
		})
	}
}
