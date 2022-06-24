package httpRequest

type Header struct {
	name  string
	value string
}

func NewHeader(name, value string) *Header {
	return &Header{name: name, value: value}
}

func (r *Header) Name() string  { return r.name }
func (r *Header) Value() string { return r.value }

func (r *Header) String() string {
	if r == nil {
		return "<none>"
	}
	return r.name + ": " + r.value
}
