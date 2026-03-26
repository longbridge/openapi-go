package util

import "time"

const SimpleDateLayout = "20060102"
const DateLayout = "2006-01-02"
const SimpleMinteLayout = "1504"

func FormatDate(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(DateLayout)
}

func FormatMinuteSimple(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(SimpleMinteLayout)
}

func FormatDateSimple(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(SimpleDateLayout)
}

func ParseDateSimple(v string) (time.Time, error) {
	if v == "" {
		return time.Time{}, nil
	}
	return time.Parse(SimpleDateLayout, v)
}
