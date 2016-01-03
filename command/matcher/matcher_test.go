package matcher

import (
	"testing"

	"github.com/botyard/botyard/command"
	"github.com/botyard/botyard/command/parse"
)

func Test_Matcher(t *testing.T) {
	var tests = []struct {
		src     string
		cmd     string
		matched bool
		args    []*command.Argument
	}{
		//		{src: "image me", cmd: "image me", matched: true, args: []*command.Argument{}},
		{src: "imge me", cmd: "image me", matched: false, args: []*command.Argument{}},
		/*		{"send hello to you ", "send <message> to <user>", true,
					[]*command.Argument{{"message", "hello"}, {"user", "you"}},
				},
		*/
	}

	for i, tt := range tests {
		items, err := parse.Parse(tt.cmd)
		if err != nil {
			t.Errorf("%q. err:%v", i, err)
		}
		t.Log(items[0], items[1])

		m := New(tt.src, items)
		args, ok := m.MatchAndReturnArguments()
		if ok != tt.matched {
			t.Errorf("%v. mismatched: want=%v have=%v", i, tt.matched, ok)
		}

		if len(args) != len(tt.args) {
			t.Errorf("%v. args mismatched: want=%v have=%v", i, tt.args, args)
		}

	}

}
