package bed

type Bed interface {
	Upload(fname string) (string, error)
}
