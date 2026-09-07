package audit

type discarded struct{}

func NewDiscarded() discarded {
	return discarded{}
}

func (discarded) Log(string, int) error {
	return nil
}
