package diary

type (
	Day   string
	Month string
	Year  string

	DayElem struct {
		Day
		Path string
	}

	MonthElem struct {
		Month
		Days []DayElem
	}
	YearElem struct {
		Year
		Months []MonthElem
	}
	TopElem struct {
		Base  string
		Years []YearElem
	}
)
