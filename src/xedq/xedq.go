package xedq

// InitTables prepares XED for encoding/decoding requests.
// TODO: specify what operations are safe/unsafe prior to this call.
func InitTables() {
	xedTablesInit()
}

// Ptr describes effective address computation.
type Ptr struct {
	// Base register name. SIB - B.
	// Empty string means "no base".
	Base string

	// Index register name. SIB - I.
	// Empty string means "no scaled index".
	Index string

	// Scaling factor. SIB - S.
	// Values 1, 2, 4 and 8 specify index scaling.
	// Value of 0 means "no explicit scaling factor",
	// which usually implies scaling factor of 1.
	Scale uint8

	// 32bit pointer displacement.
	// Displacement is encoded as 8 or 32 bit immediate value.
	// Exceptions like MOVABS are not handled (yet?).
	Disp int32
}
