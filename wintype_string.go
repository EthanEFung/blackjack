// Code generated by "stringer -type=WinType"; DO NOT EDIT.

package blackjack

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Undetermined-0]
	_ = x[Lose-1]
	_ = x[Win-2]
	_ = x[Bust-3]
	_ = x[Push-4]
}

const _WinType_name = "UndeterminedLoseWinBustPush"

var _WinType_index = [...]uint8{0, 12, 16, 19, 23, 27}

func (i WinType) String() string {
	if i < 0 || i >= WinType(len(_WinType_index)-1) {
		return "WinType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _WinType_name[_WinType_index[i]:_WinType_index[i+1]]
}
