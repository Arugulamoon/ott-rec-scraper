package timefmt_test

import (
	"testing"

	"eden-walker.com/home/ott-rec-scraper/pkg/timefmt"
)

func TestTranslateEvents(t *testing.T) {
	type test struct {
		input string
		want  []timefmt.TimeFmt
	}

	tests := []test{
		{
			input: "7 - 7:50 pm",
			want: []timefmt.TimeFmt{
				{Start: "19:00", End: "19:50"},
			},
		},
		{
			input: `11:20 am - 12:10 pm,
			12:20 - 1:10 pm`,
			want: []timefmt.TimeFmt{
				{Start: "11:20", End: "12:10"},
				{Start: "12:20", End: "13:10"},
			},
		},
		{
			input: `
			9 - 9:50 am
			`,
			want: []timefmt.TimeFmt{
				{Start: "9:00", End: "9:50"},
			},
		},
		{
			input: `
			4 - 4:50 pm 
			`,
			want: []timefmt.TimeFmt{
				{Start: "16:00", End: "16:50"},
			},
		},
	}

	for _, tc := range tests {
		got := timefmt.TranslateEvents(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("want %d; got %d", len(tc.want), len(got))
		}
		for i := 0; i < len(got); i++ {
			if got[i] != tc.want[i] {
				t.Errorf("want %v; got %v", tc.want[i], got[i])
			}
		}
	}
}

func TestTranslateTimeStrTo24H(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "11:20am", want: "11:20"},
		{input: "12:10pm", want: "12:10"},
		{input: "07:00pm", want: "19:00"},
		{input: "11:59pm", want: "23:59"},
		{input: "12:00am", want: "00:00"},
	}

	for _, tc := range tests {
		got := timefmt.TranslateTimeStrTo24H(tc.input)
		if got != tc.want {
			t.Errorf("want %s; got %s", tc.want, got)
		}
	}
}

func TestAppendAMPMToStartTime(t *testing.T) {
	type test struct {
		input timefmt.TimeFmt
		want  timefmt.TimeFmt
	}

	tests := []test{
		{
			input: timefmt.TimeFmt{Start: "7:00", End: "7:50pm"},
			want:  timefmt.TimeFmt{Start: "7:00pm", End: "7:50pm"},
		},
		{
			input: timefmt.TimeFmt{Start: "11:20am", End: "12:10pm"},
			want:  timefmt.TimeFmt{Start: "11:20am", End: "12:10pm"},
		},
	}

	for _, tc := range tests {
		got := timefmt.AppendAMPMToStartTime(tc.input)
		if got != tc.want {
			t.Errorf("want %v; got %v", tc.want, got)
		}
	}
}

func TestSplitEventTimes(t *testing.T) {
	type test struct {
		input string
		want  timefmt.TimeFmt
	}

	tests := []test{
		{
			input: "7-7:50pm",
			want:  timefmt.TimeFmt{Start: "7:00", End: "7:50pm"},
		},
		{
			input: "11:20am-12:10pm",
			want:  timefmt.TimeFmt{Start: "11:20am", End: "12:10pm"},
		},
	}

	for _, tc := range tests {
		got := timefmt.SplitEventTimes(tc.input)
		if got != tc.want {
			t.Errorf("want %s; got %s", tc.want, got)
		}
	}
}

func TestTranslateTimeStrToHHMM(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "11:20", want: "11:20"},
		{input: "7", want: "7:00"},
	}

	for _, tc := range tests {
		got := timefmt.TranslateTimeStrToHHMM(tc.input)
		if got != tc.want {
			t.Errorf("want %s; got %s", tc.want, got)
		}
	}
}

func TestSplitEvents(t *testing.T) {
	type test struct {
		input string
		want  []string
	}

	tests := []test{
		{input: "7-7:50pm", want: []string{"7-7:50pm"}},
		{
			input: "11:20am-12:10pm,12:20-1:10pm",
			want:  []string{"11:20am-12:10pm", "12:20-1:10pm"},
		},
	}

	for _, tc := range tests {
		got := timefmt.SplitEvents(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("want %d; got %d", len(tc.want), len(got))
		}
		for i := 0; i < len(got); i++ {
			if got[i] != tc.want[i] {
				t.Errorf("want %s; got %s", tc.want[i], got[i])
			}
		}
	}
}

func TestSanitizeTimes(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "7 - 7:50 pm", want: "7-7:50pm"},
		{
			input: `11:20 am - 12:10 pm, 
			12:20 - 1:10 pm`,
			want: "11:20am-12:10pm,12:20-1:10pm",
		},
	}

	for _, tc := range tests {
		got := timefmt.SanitizeTimes(tc.input)
		if got != tc.want {
			t.Errorf("want %s; got %s", tc.want, got)
		}
	}
}
