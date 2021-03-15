package retcode_test

import (
	"fmt"
	"testing"

	"github.com/cardil/kn-event/internal/cli/retcode"
	"github.com/cardil/kn-event/internal/sender"
)

func TestCalc(t *testing.T) {
	cases := testCases()
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := retcode.Calc(tt.err); got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testCases() []testCase {
	return []testCase{{
		name: "nil",
		err:  nil,
		want: 0,
	}, {
		name: "sender.ErrCouldntBeSent",
		err:  sender.ErrCouldntBeSent,
		want: 157,
	}, {
		name: "error of wrap caused by 12345",
		err:  fmt.Errorf("%w: 12345", sender.ErrCouldntBeSent),
		want: 193,
	}}
}

type testCase struct {
	name string
	err  error
	want int
}
