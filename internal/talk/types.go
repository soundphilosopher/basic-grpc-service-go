package talk

type Talk struct {
	Reply func(string) (string, bool)
}

func NewTalk(reply func(string) (string, bool)) *Talk {
	return &Talk{
		Reply: reply,
	}
}
