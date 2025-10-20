package models

import (
	"testing"
	"time"
)

func TestPoint_Timestamp(t *testing.T) {
	tests := []struct {
		name     string
		tm       int64
		wantTime time.Time
	}{
		{
			name:     "unix epoch",
			tm:       0,
			wantTime: time.Unix(0, 0),
		},
		{
			name:     "specific timestamp",
			tm:       1609459200,
			wantTime: time.Unix(1609459200, 0),
		},
		{
			name:     "negative timestamp",
			tm:       -1,
			wantTime: time.Unix(-1, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Point{Tm: tt.tm}
			got := p.Timestamp()
			if !got.Equal(tt.wantTime) {
				t.Errorf("Timestamp() = %v, want %v", got, tt.wantTime)
			}
		})
	}
}

func TestPoint_Altitude(t *testing.T) {
	tests := []struct {
		name    string
		al      string
		want    float64
		wantErr bool
	}{
		{
			name:    "valid altitude",
			al:      "123.45",
			want:    123.45,
			wantErr: false,
		},
		{
			name:    "zero altitude",
			al:      "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "negative altitude",
			al:      "-10.5",
			want:    -10.5,
			wantErr: false,
		},
		{
			name:    "empty string",
			al:      "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid altitude",
			al:      "invalid",
			want:    0,
			wantErr: true,
		},
		{
			name:    "partial invalid",
			al:      "12.34abc",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Point{Al: tt.al}
			got, err := p.Altitude()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Altitude() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Altitude() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("Altitude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint_Speed(t *testing.T) {
	tests := []struct {
		name    string
		sp      string
		want    float64
		wantErr bool
	}{
		{
			name:    "valid speed",
			sp:      "25.5",
			want:    25.5,
			wantErr: false,
		},
		{
			name:    "zero speed",
			sp:      "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "empty string",
			sp:      "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid speed",
			sp:      "fast",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Point{Sp: tt.sp}
			got, err := p.Speed()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Speed() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Speed() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("Speed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint_Distance(t *testing.T) {
	tests := []struct {
		name    string
		ds      string
		want    float64
		wantErr bool
	}{
		{
			name:    "valid distance",
			ds:      "1500.75",
			want:    1500.75,
			wantErr: false,
		},
		{
			name:    "zero distance",
			ds:      "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "empty string",
			ds:      "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid distance",
			ds:      "far",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Point{Ds: tt.ds}
			got, err := p.Distance()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Distance() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Distance() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}
