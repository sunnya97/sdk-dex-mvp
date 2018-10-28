package orderbook

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidSortableDec(t *testing.T) {
	tests := []struct {
		name      string
		decString string
		want      bool
	}{
		{"correct", "0.00013003", true},
		{"correct", "110303030.00013003", true},
		{"too big", "1100000303030.00013003", false},
		{"too small", "0.0000000000000001", false},
		{"too much precsison", "12130.0002000000000001", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec, err := sdk.NewDecFromStr(tt.decString)
			if err != nil {
				if tt.want != false {
					t.Errorf("ValidSortableDec() = %v, want %v", true, tt.want)
				}
				return
			}
			if got := ValidSortableDec(dec); got != tt.want {
				t.Errorf("ValidSortableDec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppendWithSeperator(t *testing.T) {
	type args struct {
		byteslice1 []byte
		byteslice2 []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"nonempty w/ nonempty", args{[]byte("asdf"), []byte("jkl;")}, []byte("asdf/jkl;")},
		{"nonempty w/ empty", args{[]byte("asdf"), []byte("")}, []byte("asdf/")},
		{"empty w/ nonempty", args{[]byte(""), []byte("jkl;")}, []byte("/jkl;")},
		{"empty w/ empty", args{[]byte(""), []byte("")}, []byte("/")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendWithSeperator(tt.args.byteslice1, tt.args.byteslice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendWithSeperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitKeyAlongSeperator(t *testing.T) {
	tests := []struct {
		name    string
		fullKey []byte
		want    [][]byte
	}{
		{"2 layer", []byte("asdf/jkl;"), [][]byte{[]byte("asdf"), []byte("jkl;")}},
		{"3 layer", []byte("asdf/jkl;/qwerty"), [][]byte{[]byte("asdf"), []byte("jkl;"), []byte("qwerty")}},
		{"empty sections", []byte("asdf//qwerty"), [][]byte{[]byte("asdf"), []byte(""), []byte("qwerty")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitKeyAlongSeperator(tt.fullKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitKeyAlongSeperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSDKDecReciprocal(t *testing.T) {
	tests := []struct {
		name       string
		decString  string
		wantString string
	}{
		{"power ten", "0.00001", "100000"},
		{"common fractions", "0.25", "4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec, _ := sdk.NewDecFromStr(tt.decString)
			want, _ := sdk.NewDecFromStr(tt.wantString)

			if got := SDKDecReciprocal(dec); !reflect.DeepEqual(got, want) {
				t.Errorf("SDKDecReciprocal() = %v, want %v", got, want)
			}
		})
	}
}

func TestSortableSDKDecBytes(t *testing.T) {
	tests := []struct {
		name string
		dec  sdk.Dec
		want []byte
	}{
		{"one", sdk.OneDec(), []byte("00000000010000000000")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortableSDKDecBytes(tt.dec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortableSDKDecBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
