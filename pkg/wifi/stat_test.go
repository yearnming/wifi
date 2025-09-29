package wifi

import "testing"

func TestGetWlanStat(t *testing.T) {
	stat, err := GetWIFIStat()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(stat)
}
