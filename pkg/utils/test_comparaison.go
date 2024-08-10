package utils

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
)

func EqualErrors(t *testing.T, expectedErr, receivedErr error) {
	t.Helper()

	if !errors.Is(expectedErr, receivedErr) {
		t.Errorf("expectedErr: %v NOT EQUAL receivedErr: %v", expectedErr, receivedErr)
	}
}

func PlainEqual(t *testing.T, expected, received any) {
	t.Helper()

	if expected != received {
		t.Errorf("expected: %v NOT EQUAL received: %v",
			expected, received)
	}
}

func DeepEqual(t *testing.T, expected, received any) {
	t.Helper()

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("expected: %v NOT DEEP EQUAL received: %v",
			expected, received)
	}
}

func PgxPoolExpectationWereMet(t *testing.T, mock pgxmock.PgxPoolIface) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("not expected err: %s", err)
	}
}
