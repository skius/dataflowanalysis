package main

type AbsStringType int

const (
	TypeTop AbsStringType = iota
	TypeBottom
	TypeConst
)

/*
	The lattice of abstract strings:
                    Top                          (May be any string)
      /      /      ...    \       \
 ... "a"  "hello"   ...    "42"    "41" ....     (Must be a specific constant)
      \      \      ...     /       /
                   Bottom                        (Uninitialized)
 */

type AbsString struct {
	Type AbsStringType
	Constant string
}

func (e *AbsString) IsTop() bool {
	return e.Type == TypeTop
}
func (e *AbsString) IsBottom() bool {
	return e.Type == TypeBottom
}
func (e *AbsString) IsConstant() bool {
	return e.Type == TypeConst
}

func Bottom() *AbsString {
	bot := new(AbsString)
	bot.Type = TypeBottom
	return bot
}
func Top() *AbsString {
	top := new(AbsString)
	top.Type = TypeTop
	return top
}
func Constant(s string) *AbsString {
	c := new(AbsString)
	c.Type = TypeConst
	c.Constant = s
	return c
}

func (e *AbsString) Meet(other *AbsString) *AbsString {
	if e.IsTop() {
		return other.copy()
	}
	if other.IsTop() {
		return e.copy()
	}
	if e.IsBottom() || other.IsBottom() {
		return e.copy()
	}
	if e.IsConstant() && other.IsConstant() && e.Constant == other.Constant {
		return e.copy()
	}
	return Bottom()
}
func (e *AbsString) Join(other *AbsString) *AbsString {
	if e.IsBottom() {
		return other.copy()
	}
	if other.IsBottom() {
		return e.copy()
	}
	if e.IsTop() || other.IsTop() {
		return e.copy()
	}
	if e.IsConstant() && other.IsConstant() && e.Constant == other.Constant {
		return e.copy()
	}
	return Top()
}

func (e *AbsString) Equals(other *AbsString) bool {
	if e.IsBottom() && other.IsBottom() {
		return true
	}
	if e.IsTop() && other.IsTop() {
		return true
	}
	if e.IsConstant() && other.IsConstant() {
		return e.Constant == other.Constant
	}
	return false
}

func (e *AbsString) String() string {
	if e.IsTop() {
		return "<Top>"
	}
	if e.IsBottom() {
		return "<Bottom>"
	}
	return `"`+ e.Constant +`"`
}

func (e *AbsString) copy() *AbsString {
	e2 := new(AbsString)
	e2.Constant = e.Constant
	e2.Type = e.Type
	return e2
}
