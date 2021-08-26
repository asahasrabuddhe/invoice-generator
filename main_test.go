package main

import (
	"html/template"
	"testing"
	"time"
)

func TestOrdinal(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args {
				1,
			},
			want: "1st",
		},
		{
			name: "2",
			args: args {
				2,
			},
			want: "2nd",
		},
		{
			name: "3",
			args: args {
				3,
			},
			want: "3rd",
		},
		{
			name: "4",
			args: args {
				4,
			},
			want: "4th",
		},
		{
			name: "5",
			args: args {
				5,
			},
			want: "5th",
		},
		{
			name: "6",
			args: args {
				6,
			},
			want: "6th",
		},
		{
			name: "7",
			args: args {
				7,
			},
			want: "7th",
		},
		{
			name: "8",
			args: args {
				8,
			},
			want: "8th",
		},
		{
			name: "9",
			args: args {
				9,
			},
			want: "9th",
		},
		{
			name: "10",
			args: args {
				10,
			},
			want: "10th",
		},
		{
			name: "11",
			args: args {
				11,
			},
			want: "11th",
		},
		{
			name: "12",
			args: args {
				12,
			},
			want: "12th",
		},
		{
			name: "13",
			args: args {
				13,
			},
			want: "13th",
		},
		{
			name: "14",
			args: args {
				14,
			},
			want: "14th",
		},
		{
			name: "15",
			args: args {
				15,
			},
			want: "15th",
		},
		{
			name: "16",
			args: args {
				16,
			},
			want: "16th",
		},
		{
			name: "17",
			args: args {
				17,
			},
			want: "17th",
		},
		{
			name: "18",
			args: args {
				18,
			},
			want: "18th",
		},
		{
			name: "19",
			args: args {
				19,
			},
			want: "19th",
		},
		{
			name: "20",
			args: args {
				20,
			},
			want: "20th",
		},
		{
			name: "21",
			args: args {
				21,
			},
			want: "21st",
		},
		{
			name: "22",
			args: args {
				22,
			},
			want: "22nd",
		},
		{
			name: "23",
			args: args {
				23,
			},
			want: "23rd",
		},
		{
			name: "24",
			args: args {
				24,
			},
			want: "24th",
		},
		{
			name: "25",
			args: args {
				25,
			},
			want: "25th",
		},
		{
			name: "26",
			args: args {
				26,
			},
			want: "26th",
		},
		{
			name: "27",
			args: args {
				27,
			},
			want: "27th",
		},
		{
			name: "28",
			args: args {
				28,
			},
			want: "28th",
		},
		{
			name: "29",
			args: args {
				29,
			},
			want: "29th",
		},
		{
			name: "30",
			args: args {
				30,
			},
			want: "30th",
		},
		{
			name: "31",
			args: args {
				31,
			},
			want: "31st",
		},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ordinal(tt.args.x); got != tt.want {
				t.Errorf("Ordinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdinalDate(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "16/09/2019",
			args: args{
				date: time.Date(2019, time.September, 16, 0, 0, 0, 0, loc),
			},
			want: "16th September 2019",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrdinalDate(tt.args.date); got != tt.want {
				t.Errorf("OrdinalDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDescription(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want template.HTML
	}{
		{
			name: "1",
			args: args{
				//line: "5 days of work done in between 16th September 2019 and "
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDescription(tt.args.line); got != tt.want {
				t.Errorf("FormatDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}