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

func (s *AbsString) IsTop() bool {
	return s.Type == TypeTop
}
func (s *AbsString) IsBottom() bool {
	return s.Type == TypeBottom
}
func (s *AbsString) IsConstant() bool {
	return s.Type == TypeConst
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

func (s *AbsString) Meet(other *AbsString) *AbsString {
	if s.IsTop() {
		return other.copy()
	}
	if other.IsTop() {
		return s.copy()
	}
	if s.IsBottom() || other.IsBottom() {
		return s.copy()
	}
	if s.IsConstant() && other.IsConstant() && s.Constant == other.Constant {
		return s.copy()
	}
	return Bottom()
}
func (s *AbsString) Join(other *AbsString) *AbsString {
	if s.IsBottom() {
		return other.copy()
	}
	if other.IsBottom() {
		return s.copy()
	}
	if s.IsTop() || other.IsTop() {
		return s.copy()
	}
	if s.IsConstant() && other.IsConstant() && s.Constant == other.Constant {
		return s.copy()
	}
	return Top()
}

func (s *AbsString) Equals(other *AbsString) bool {
	if s.IsBottom() && other.IsBottom() {
		return true
	}
	if s.IsTop() && other.IsTop() {
		return true
	}
	if s.IsConstant() && other.IsConstant() {
		return s.Constant == other.Constant
	}
	return false
}

func (s *AbsString) String() string {
	if s.IsTop() {
		return "<Top>" // "⊤" //
	}
	if s.IsBottom() {
		return "<Bottom>" // "⊥" //
	}
	return `"`+ s.Constant +`"`
}

func (s *AbsString) copy() *AbsString {
	e2 := new(AbsString)
	e2.Constant = s.Constant
	e2.Type = s.Type
	return e2
}
