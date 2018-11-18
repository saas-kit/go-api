package signeddata

import (
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	sKey := "secret"
	userID := "12345qwerty"
	payload := map[string]interface{}{"user_id": userID}
	exp := time.Now().Add(120 * time.Second)

	type args struct {
		key     string
		payload map[string]interface{}
		exp     time.Time
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Encode(): successful",
			args:    args{sKey, payload, exp},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signedString, err := Encode(tt.args.key, tt.args.payload, tt.args.exp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err := Decode(sKey, signedString); err != nil {
				t.Errorf("Encode(): %v", err)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	sKey := "secret"
	userID := "12345qwerty"
	payload := map[string]interface{}{"user_id": userID}

	signedString1, _ := Encode(sKey, payload, time.Now().Add(120*time.Second))
	signedString2, _ := Encode(sKey, payload, time.Now().Add(-time.Hour))

	type args struct {
		key          string
		signedString string
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			"Decode(): successful",
			args{sKey, signedString1},
			payload,
			false,
		},
		{
			"Decode(): expired",
			args{sKey, signedString2},
			payload,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.key, tt.args.signedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got["user_id"] != userID {
				t.Errorf("Decode() = %v, want %v", got["user_id"], tt.want["user_id"])
			}
		})
	}
}
