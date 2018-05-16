package variable

type Variable struct {
	Name string
	ScopeName string
	Type int
	StringValue string
	IntegerValue int
	FloatValue float64
	IsConstant bool
	ArrayValue []Variable
}