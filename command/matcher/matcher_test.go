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
		{src: "image me", cmd: "image me", matched: true, args: []*command.Argument{}},
		{src: "imge me", cmd: "image me", matched: false, args: []*command.Argument{}},
		{src: "image ", cmd: "image me", matched: false, args: []*command.Argument{}},
		{src: "image", cmd: "image me", matched: false, args: []*command.Argument{}},
		{src: "imageme", cmd: "image me", matched: false, args: []*command.Argument{}},
		{"send hello to you ", "send <message> to <user>", true,
			[]*command.Argument{{"message", "hello"}, {"user", "you"}},
		},
		{"test  send hello to you ", "send <message> to <user>", true,
			[]*command.Argument{{"message", "hello"}, {"user", "you"}},
		},
		{"send message hello to you ", "send message <message> to <user>", true,
			[]*command.Argument{{"message", "hello"}, {"user", "you"}},
		},
		{"send hello world ", "send <message1> <message2>", true,
			[]*command.Argument{{"message1", "hello"}, {"message2", "world"}},
		},
	}

	for i, tt := range tests {
		items, err := parse.Parse(tt.cmd)
		if err != nil {
			t.Errorf("%q. err:%v", i, err)
		}
		//t.Log(items[0], items[1])

		m := New(tt.src, items)
		ok, _ := m.Match()
		args := m.Arguments()

		if ok != tt.matched {
			t.Errorf("%v. src='%v',cmd='%v' mismatched: want=%v have=%v", i, tt.src, tt.cmd, tt.matched, ok)
		}

		if len(args) != len(tt.args) {
			t.Errorf("%v. args mismatched: want=%v have=%v cmd=\"%v\" text=\"%v\"", i, tt.args, args, tt.cmd, tt.src)
		}

	}

}
