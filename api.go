package teq

type Teq struct{}

func (teq *Teq) Equal(t TestingT, a, b interface{}) bool {
	ok := a == b
	if !ok {
		t.Errorf("expected %v, got %v", a, b)
	}
	return ok
}
