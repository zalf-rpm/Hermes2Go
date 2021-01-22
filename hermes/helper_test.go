package hermes

import (
	"reflect"
	"testing"
)

func TestDateConverter(t *testing.T) {
	type args struct {
		splitCenturyAt         int
		dateOldDE              string
		dateNewDE              string
		dateOldEN              string
		dateNewEN              string
		dateOldWithSeperatorDE string
		dateNewWithSeperatorDE string
		dateOldWithSeperatorEN string
		dateNewWithSeperatorEN string
	}
	tests := []struct {
		name string
		args args
	}{
		{"testset1", args{50, "100206", "10022006", "021006", "02102006", "10.02.06", "10.02.2006", "02.10.06", "02.10.2006"}},
		{"testset2", args{60, "100256", "10022056", "021056", "02.10.2056", "10.02.56", "10.02.2056", "02.10.56", "02.10.2056"}},
		{"testset3", args{50, "100256", "10021956", "021056", "02101956", "10.02.56", "10.02.1956", "02.10.56", "02.10.1956"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateOldDE := DateConverter(tt.args.splitCenturyAt, DateDEshort)
			dateNewDE := DateConverter(-1, DateDElong)
			dateOldEN := DateConverter(tt.args.splitCenturyAt, DateENshort)
			dateNewEN := DateConverter(-1, DateENlong)

			ztDatOldDE, masDatOldDE := dateOldDE(tt.args.dateOldDE)
			ztDatNewDE, masDatNewDE := dateNewDE(tt.args.dateNewDE)
			ztDatOldEN, masDatOldEN := dateOldEN(tt.args.dateOldEN)
			ztDatNewEN, masDatNewEN := dateNewEN(tt.args.dateNewEN)

			ztDatDatum, masDatDatum := datumOld(tt.args.dateOldDE, tt.args.splitCenturyAt)
			if !reflect.DeepEqual(ztDatOldDE, ztDatNewDE) || !reflect.DeepEqual(ztDatOldDE, ztDatDatum) {
				t.Errorf("DateConverter() '%s' ztDat: oldStyle= %v, newStyle= %v, want %v", tt.name, ztDatOldDE, ztDatNewDE, ztDatDatum)
			}
			if !reflect.DeepEqual(masDatOldDE, masDatNewDE) || !reflect.DeepEqual(masDatOldDE, masDatDatum) {
				t.Errorf("DateConverter() '%s' masDat: oldStyle= %v, newStyle= %v, want %v", tt.name, masDatOldDE, masDatNewDE, masDatDatum)
			}
			if !reflect.DeepEqual(ztDatOldEN, ztDatNewEN) || !reflect.DeepEqual(ztDatOldEN, ztDatDatum) {
				t.Errorf("DateConverter() '%s' ztDat: oldStyle= %v, newStyle= %v, want %v", tt.name, ztDatOldEN, ztDatNewEN, ztDatDatum)
			}
			if !reflect.DeepEqual(masDatOldEN, masDatNewEN) || !reflect.DeepEqual(masDatOldEN, masDatDatum) {
				t.Errorf("DateConverter() '%s' masDat: oldStyle= %v, newStyle= %v, want %v", tt.name, masDatOldEN, masDatNewEN, masDatDatum)
			}
			// test with seperators
			ztDatOldDEsep, masDatOldDEsep := dateOldDE(tt.args.dateOldWithSeperatorDE)
			ztDatNewDEsep, masDatNewDEsep := dateNewDE(tt.args.dateNewWithSeperatorDE)
			ztDatOldENsep, masDatOldENsep := dateOldEN(tt.args.dateOldWithSeperatorEN)
			ztDatNewENsep, masDatNewENsep := dateNewEN(tt.args.dateNewWithSeperatorEN)

			if !reflect.DeepEqual(ztDatOldDEsep, ztDatNewDEsep) || !reflect.DeepEqual(ztDatOldDEsep, ztDatDatum) {
				t.Errorf("DateConverter() '%s' ztDat: oldStyle= %v, newStyle= %v, want %v", tt.name, ztDatOldDEsep, ztDatNewDEsep, ztDatDatum)
			}
			if !reflect.DeepEqual(masDatOldDEsep, masDatNewDEsep) || !reflect.DeepEqual(masDatOldDEsep, masDatDatum) {
				t.Errorf("DateConverter() '%s' masDat: oldStyle= %v, newStyle= %v, want %v", tt.name, masDatOldDEsep, masDatNewDEsep, masDatDatum)
			}
			if !reflect.DeepEqual(ztDatOldENsep, ztDatNewENsep) || !reflect.DeepEqual(ztDatOldENsep, ztDatDatum) {
				t.Errorf("DateConverter() '%s' ztDat: oldStyle= %v, newStyle= %v, want %v", tt.name, ztDatOldENsep, ztDatNewENsep, ztDatDatum)
			}
			if !reflect.DeepEqual(masDatOldENsep, masDatNewEN) || !reflect.DeepEqual(masDatOldENsep, masDatDatum) {
				t.Errorf("DateConverter() '%s' masDat: oldStyle= %v, newStyle= %v, want %v", tt.name, masDatOldENsep, masDatNewENsep, masDatDatum)
			}
		})
	}
}

func TestKalenderDate(t *testing.T) {
	type args struct {
		MASDAT    int
		DateShort string
		DateLong  string
		cent      int
	}
	tests := []struct {
		name      string
		args      args
		wantYear  int
		wantMonth int
		wantDay   int
	}{
		{"testset3", args{1, "010101", "01011901", 0}, 1901, 1, 1},
		{"testset1", args{38392, "100206", "10022006", 50}, 2006, 2, 10},
		{"testset2", args{20129, "100256", "10021956", 50}, 1956, 2, 10},
		{"testset4", args{56654, "100256", "10022056", 60}, 2056, 2, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateOldDE := DateConverter(tt.args.cent, DateDEshort)
			dateNewDE := DateConverter(-1, DateDElong)

			_, masDatOldDE := dateOldDE(tt.args.DateShort)
			_, masDatNewDE := dateNewDE(tt.args.DateLong)

			gotYear, gotMonth, gotDay := KalenderDate(masDatOldDE)
			if gotYear != tt.wantYear {
				t.Errorf("KalenderDate() gotYear = %v, want %v", gotYear, tt.wantYear)
			}
			if gotMonth != tt.wantMonth {
				t.Errorf("KalenderDate() gotMonth = %v, want %v", gotMonth, tt.wantMonth)
			}
			if gotDay != tt.wantDay {
				t.Errorf("KalenderDate() gotDay = %v, want %v", gotDay, tt.wantDay)
			}

			gotYear, gotMonth, gotDay = KalenderDate(masDatNewDE)
			if gotYear != tt.wantYear {
				t.Errorf("KalenderDate() gotYear = %v, want %v", gotYear, tt.wantYear)
			}
			if gotMonth != tt.wantMonth {
				t.Errorf("KalenderDate() gotMonth = %v, want %v", gotMonth, tt.wantMonth)
			}
			if gotDay != tt.wantDay {
				t.Errorf("KalenderDate() gotDay = %v, want %v", gotDay, tt.wantDay)
			}

			gotYear, gotMonth, gotDay = KalenderDate(tt.args.MASDAT)
			if gotYear != tt.wantYear {
				t.Errorf("KalenderDate() gotYear = %v, want %v", gotYear, tt.wantYear)
			}
			if gotMonth != tt.wantMonth {
				t.Errorf("KalenderDate() gotMonth = %v, want %v", gotMonth, tt.wantMonth)
			}
			if gotDay != tt.wantDay {
				t.Errorf("KalenderDate() gotDay = %v, want %v", gotDay, tt.wantDay)
			}

		})
	}
}

func TestKalenderConverter(t *testing.T) {
	type args struct {
		dateformat DateFormat
		MASDAT     int
	}
	tests := []struct {
		name     string
		args     args
		wantDate string
	}{
		{"testset1", args{DateDEshort, 2}, "02.01.01"},
		{"testset2", args{DateDEshort, 38392}, "10.02.06"},
		{"testset3", args{DateDEshort, 20129}, "10.02.56"},
		{"testset4", args{DateDEshort, 56654}, "10.02.56"},
		{"testset5", args{DateENshort, 2}, "01.02.01"},
		{"testset6", args{DateENshort, 38392}, "02.10.06"},
		{"testset7", args{DateENshort, 20129}, "02.10.56"},
		{"testset8", args{DateENshort, 56654}, "02.10.56"},
		{"testset9", args{DateDElong, 2}, "02.01.1901"},
		{"testset10", args{DateDElong, 38392}, "10.02.2006"},
		{"testset11", args{DateDElong, 20129}, "10.02.1956"},
		{"testset12", args{DateDElong, 56654}, "10.02.2056"},
		{"testset13", args{DateENlong, 2}, "01.02.1901"},
		{"testset14", args{DateENlong, 38392}, "02.10.2006"},
		{"testset15", args{DateENlong, 20129}, "02.10.1956"},
		{"testset16", args{DateENlong, 56654}, "02.10.2056"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := KalenderConverter(tt.args.dateformat, ".")

			if got := converter(tt.args.MASDAT); !reflect.DeepEqual(got, tt.wantDate) {
				t.Errorf("KalenderConverter() = %v, want %v", got, tt.wantDate)
			}
		})
	}
}
