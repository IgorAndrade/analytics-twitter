package usecase

type Search interface {
	Find(map[string]string) error
}

func (t Twitter) Find(map[string]string) error {
	return nil
}
