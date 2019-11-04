package sqlTools

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestCreatePacketQuery(t *testing.T) {
	input := []struct{
		prefix string
		batchSize int
		batchCount int
		postfix []string
	}{
		{ `INSERT INTO Receiver (mailId, email) VALUES`, 2, 3, []string{}},
		{ `INSERT INTO Receiver (mailId, email) VALUES `, 2, 3, []string{`RETURNING id`}},

	}
	expected := []string {
		`INSERT INTO Receiver (mailId, email) VALUES ($1, $2), ($3, $4), ($5, $6);`,
		`INSERT INTO Receiver (mailId, email) VALUES ($1, $2), ($3, $4), ($5, $6) RETURNING id;`,
	}
	for i, test := range input {
		got := CreatePacketQuery(test.prefix, test.batchSize, test.batchCount, test.postfix...)

		if !cmp.Equal(expected[i], got) {
			t.Errorf("Wrong answer. \nGot %s \ninstead %s", got, expected[i])
		}

	}
}

func TestFormatDate(t *testing.T) {
	location, _ := time.LoadLocation("America/New_York")
	input := time.Date(2000, 12, 31, 10, 55, 12, 0, location)
	expected := "2000-12-31 10:55:12"

	if got := FormatDate(BDPostgres, input); got != expected {
		t.Errorf("Wrong result: %s instead %s", got, expected)
	}

}