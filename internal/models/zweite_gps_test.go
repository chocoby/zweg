package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFirstTitle(t *testing.T) {
	tests := []struct {
		name   string
		points []Point
		want   string
	}{
		{"empty slice", nil, ""},
		{"no titles", []Point{{Tl: ""}, {Tl: ""}}, ""},
		{"first point has title", []Point{{Tl: "Morning Run"}, {Tl: ""}}, "Morning Run"},
		{"later point has title", []Point{{Tl: ""}, {Tl: "Evening Walk"}}, "Evening Walk"},
		{"multiple titles - first wins", []Point{{Tl: "First"}, {Tl: "Second"}}, "First"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstTitle(tt.points); got != tt.want {
				t.Errorf("FirstTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPoint_DecodeFullSpec(t *testing.T) {
	raw := `{
		"tm": 1609459200, "lo": 139.745438, "la": 35.658581,
		"al": "150.0", "sp": "1.0", "co": 90, "th": 85, "he": 88, "ds": "0.0",
		"ap": 1013.25, "dp": "memo", "ha": 5.0, "ms": 1, "ow": "iPhone",
		"ra": 12.5, "tl": "Morning Run", "va": 3.0, "ws": 1234, "xa": 1.5,
		"gx": 0.01, "gy": -0.02, "gz": 0.99,
		"ax": 0.1, "ay": 0.2, "az": 0.3,
		"ep": 0.5, "er": -0.5, "ey": 1.0, "pf": 2.5
	}`

	var p Point
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checks := []struct {
		name string
		got  any
		want any
	}{
		{"Tm", p.Tm, int64(1609459200)},
		{"Lo", p.Lo, 139.745438},
		{"La", p.La, 35.658581},
		{"Al", p.Al, "150.0"},
		{"Sp", p.Sp, "1.0"},
		{"Co", p.Co, 90},
		{"Th", p.Th, 85},
		{"He", p.He, 88},
		{"Ds", p.Ds, "0.0"},
		{"Ap", p.Ap, 1013.25},
		{"Dp", p.Dp, "memo"},
		{"Ha", p.Ha, 5.0},
		{"Ms", *p.Ms, MeansJogging},
		{"Ow", p.Ow, "iPhone"},
		{"Ra", p.Ra, 12.5},
		{"Tl", p.Tl, "Morning Run"},
		{"Va", p.Va, 3.0},
		{"Ws", p.Ws, 1234},
		{"Xa", p.Xa, 1.5},
		{"Gx", p.Gx, 0.01},
		{"Gy", p.Gy, -0.02},
		{"Gz", p.Gz, 0.99},
		{"Ax", p.Ax, 0.1},
		{"Ay", p.Ay, 0.2},
		{"Az", p.Az, 0.3},
		{"Ep", p.Ep, 0.5},
		{"Er", p.Er, -0.5},
		{"Ey", p.Ey, 1.0},
		{"Pf", p.Pf, 2.5},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
}

func TestPoint_Means(t *testing.T) {
	walking := MeansWalking
	jogging := MeansJogging
	bicycle := MeansBicycle
	motorcycle := MeansMotorCycle
	automobile := MeansAutoMobile
	train := MeansTrain
	misc := MeansMisc

	tests := []struct {
		name string
		raw  string
		want *Means
	}{
		{"walking", `{"ms":0}`, &walking},
		{"jogging", `{"ms":1}`, &jogging},
		{"bicycle", `{"ms":2}`, &bicycle},
		{"motorcycle", `{"ms":3}`, &motorcycle},
		{"automobile", `{"ms":4}`, &automobile},
		{"train", `{"ms":5}`, &train},
		{"misc", `{"ms":6}`, &misc},
		{"omitted", `{}`, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Point
			if err := json.Unmarshal([]byte(tt.raw), &p); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			switch {
			case tt.want == nil && p.Ms != nil:
				t.Errorf("Ms = %v, want nil", *p.Ms)
			case tt.want != nil && p.Ms == nil:
				t.Errorf("Ms = nil, want %v", *tt.want)
			case tt.want != nil && p.Ms != nil && *p.Ms != *tt.want:
				t.Errorf("Ms = %d, want %d", *p.Ms, *tt.want)
			}
		})
	}
}

func TestMeans_String(t *testing.T) {
	tests := []struct {
		m    Means
		want string
	}{
		{MeansWalking, "Walking"},
		{MeansJogging, "Jogging"},
		{MeansBicycle, "Bicycle"},
		{MeansMotorCycle, "MotorCycle"},
		{MeansAutoMobile, "AutoMobile"},
		{MeansTrain, "Train"},
		{MeansMisc, "Misc"},
		{Means(99), ""},
	}
	for _, tt := range tests {
		if got := tt.m.String(); got != tt.want {
			t.Errorf("Means(%d).String() = %q, want %q", tt.m, got, tt.want)
		}
	}
}

func TestFirstMeans(t *testing.T) {
	walking := MeansWalking
	bicycle := MeansBicycle

	tests := []struct {
		name   string
		points []Point
		want   Means
		wantOK bool
	}{
		{"empty slice", nil, 0, false},
		{"no means", []Point{{}, {}}, 0, false},
		{"first has means", []Point{{Ms: &walking}, {}}, MeansWalking, true},
		{"later has means", []Point{{}, {Ms: &bicycle}}, MeansBicycle, true},
		{"multiple - first wins", []Point{{Ms: &walking}, {Ms: &bicycle}}, MeansWalking, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FirstMeans(tt.points)
			if ok != tt.wantOK {
				t.Errorf("FirstMeans ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && got != tt.want {
				t.Errorf("FirstMeans = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint_TimestampIn(t *testing.T) {
	tests := []struct {
		name string
		tm   int64
		loc  *time.Location
		want time.Time
	}{
		{
			name: "UTC",
			tm:   1609459200, // 2021-01-01 00:00:00 UTC
			loc:  time.UTC,
			want: time.Unix(1609459200, 0).UTC(),
		},
		{
			name: "JST (+09:00)",
			tm:   1609459200,
			loc:  time.FixedZone("", 9*3600),
			want: time.Unix(1609459200, 0).In(time.FixedZone("", 9*3600)),
		},
		{
			name: "EST (-05:00)",
			tm:   1609459200,
			loc:  time.FixedZone("", -5*3600),
			want: time.Unix(1609459200, 0).In(time.FixedZone("", -5*3600)),
		},
		{
			name: "India Standard Time (+05:30)",
			tm:   1609459200,
			loc:  time.FixedZone("", 5*3600+30*60),
			want: time.Unix(1609459200, 0).In(time.FixedZone("", 5*3600+30*60)),
		},
		{
			name: "negative timestamp with positive offset",
			tm:   -1,
			loc:  time.FixedZone("", 9*3600),
			want: time.Unix(-1, 0).In(time.FixedZone("", 9*3600)),
		},
		{
			name: "unix epoch in UTC",
			tm:   0,
			loc:  time.UTC,
			want: time.Unix(0, 0).UTC(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Point{Tm: tt.tm}
			got := p.TimestampIn(tt.loc)
			if !got.Equal(tt.want) {
				t.Errorf("TimestampIn() = %v, want %v", got, tt.want)
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
