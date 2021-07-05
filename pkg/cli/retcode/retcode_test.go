package retcode_test

import (
	"fmt"
	"testing"

	"knative.dev/kn-plugin-event/pkg/cli/retcode"
	"knative.dev/kn-plugin-event/pkg/event"
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
		name: "event.ErrCantSentEvent",
		err:  event.ErrCantSentEvent,
		want: 111,
	}, {
		name: "error of wrap caused by 12345",
		err:  fmt.Errorf("%w: 12345", event.ErrCantSentEvent),
		want: 177,
	}}
}

type testCase struct {
	name string
	err  error
	want int
}
