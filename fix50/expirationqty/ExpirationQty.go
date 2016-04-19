package expirationqty

//New returns an initialized ExpirationQty instance
func New() *ExpirationQty {
	var m ExpirationQty
	return &m
}

//NoExpiration is a repeating group in ExpirationQty
type NoExpiration struct {
	//ExpType is a non-required field for NoExpiration.
	ExpType *int `fix:"982"`
	//ExpQty is a non-required field for NoExpiration.
	ExpQty *float64 `fix:"983"`
}

//NewNoExpiration returns an initialized NoExpiration instance
func NewNoExpiration() *NoExpiration {
	var m NoExpiration
	return &m
}

func (m *NoExpiration) SetExpType(v int)    { m.ExpType = &v }
func (m *NoExpiration) SetExpQty(v float64) { m.ExpQty = &v }

//ExpirationQty is a fix50 Component
type ExpirationQty struct {
	//NoExpiration is a non-required field for ExpirationQty.
	NoExpiration []NoExpiration `fix:"981,omitempty"`
}

func (m *ExpirationQty) SetNoExpiration(v []NoExpiration) { m.NoExpiration = v }
