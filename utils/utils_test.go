package utils

import (
	"testing"

	"cloud.google.com/go/civil"
)

func TestDateToString(t *testing.T) {
	// Uncomment for failure
	//now, err := civil.ParseDate("12.10.2022")
	now, err := civil.ParseDate("2002-10-02")
	if err != nil {
		t.Errorf("civil.ParseDate returned error")
		t.FailNow()
	}
	string_date := DateToString(&now)

	if string_date != "2.10.2002" {
		t.Errorf("DateToString cannot convert 2002-10-02 to string")
		t.FailNow()
	}
}

func TestFloatToMoney(t *testing.T) {
	float_val := 10.2566
	money_float := FloatToMoney(float_val)

	if money_float != 10.26 {
		t.Errorf("FloatToMoney cannot convert 10.2566")
		t.FailNow()
	}
}
