package formula

import (
	"bytes"
	"testing"
	"time"
)

func Test_printDays(t *testing.T) {
	type args struct {
		days  []WorkDays
		now   time.Time
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			name: "Print 3 days",
			args: args{
				days: []WorkDays{
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-06",
								Time: "09:00",
							},
							{
								Date: "2020-10-06",
								Time: "12:00",
							},
							{
								Date: "2020-10-06",
								Time: "13:00",
							},
							{
								Date: "2020-10-06",
								Time: "18:00",
							},
						},
					},
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-07",
								Time: "09:00",
							},
							{
								Date: "2020-10-07",
								Time: "12:00",
							},
							{
								Date: "2020-10-07",
								Time: "13:00",
							},
							{
								Date: "2020-10-07",
								Time: "19:00",
							},
						},
					},
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-08",
								Time: "01:00",
							},
							{
								Date: "2020-10-08",
								Time: "23:00",
							},
						},
					},
				},
			},
			wantW: `Work Hours:
---
Data: 2020-10-06
- 09:00
- 12:00
- 13:00
- 18:00
WorkTime: 08:00
---
Data: 2020-10-07
- 09:00
- 12:00
- 13:00
- 19:00
WorkTime: 09:00
---
Data: 2020-10-08
- 01:00
- 23:00
WorkTime: 22:00
WeekTime: 39:00
`,
		},
		{
			name: "Print invalid days",
			args: args{
				days: []WorkDays{
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-06",
								Time: "09:00",
							},
							{
								Date: "2020-10-06",
								Time: "12:00",
							},
							{
								Date: "2020-10-06",
								Time: "13:00",
							},
						},
					},
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-07",
								Time: "09:00",
							},
							{
								Date: "2020-10-07",
								Time: "12:00",
							},
							{
								Date: "2020-10-07",
								Time: "13:00",
							},
							{
								Date: "2020-10-07",
								Time: "19:00",
							},
						},
					},
				},
			},
			wantW: `Work Hours:
---
Data: 2020-10-06
- 09:00
- 12:00
- 13:00
WorkTime: invalid
---
Data: 2020-10-07
- 09:00
- 12:00
- 13:00
- 19:00
WorkTime: 09:00
WeekTime: 09:00
`,
		},
		{
			name: "Invalid but today",
			args: args{
				days: []WorkDays{
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-06",
								Time: "09:00",
							},
							{
								Date: "2020-10-06",
								Time: "12:00",
							},
							{
								Date: "2020-10-06",
								Time: "13:00",
							},
						},
					},
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-07",
								Time: "09:00",
							},
							{
								Date: "2020-10-07",
								Time: "12:00",
							},
							{
								Date: "2020-10-07",
								Time: "13:00",
							},
						},
					},
				},
				now: func() time.Time {
					t, _ := time.Parse("2006-01-02 15:04", "2020-10-07 14:00")
					return t
				}(),
			},
			wantW: `Work Hours:
---
Data: 2020-10-06
- 09:00
- 12:00
- 13:00
WorkTime: invalid
---
Data: 2020-10-07
- 09:00
- 12:00
- 13:00
WorkTime: 04:00
WeekTime: 04:00
`,
		},
		{
			name: "Invalid today byt only 1 date",
			args: args{
				days: []WorkDays{
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-06",
								Time: "09:00",
							},
							{
								Date: "2020-10-06",
								Time: "12:00",
							},
							{
								Date: "2020-10-06",
								Time: "13:00",
							},
						},
					},
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{

							{
								Date: "2020-10-07",
								Time: "09:00",
							},
						},
					},
				},
				now: func() time.Time {
					t, _ := time.Parse("2006-01-02 15:04", "2020-10-07 14:00")
					return t
				}(),
			},
			wantW: `Work Hours:
---
Data: 2020-10-06
- 09:00
- 12:00
- 13:00
WorkTime: invalid
---
Data: 2020-10-07
- 09:00
WorkTime: 05:00
WeekTime: 05:00
`,
		},
		{
			name: "empty",
			args: args{
				days: []WorkDays{
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{},
					},
					{
						TimeCards: []struct {
							Date string `json:"date"`
							Time string `json:"time"`
						}{},
					},
				},
			},
			wantW: `Work Hours:
---
WorkTime: 00:00
---
WorkTime: 00:00
WeekTime: 00:00
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			printDays(tt.args.days, w, tt.args.now)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("printDays() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
