package groot

// Class represents a ROOT class.
// Class instances are created by a ClassFactory.
type Class interface {
	// GetCheckSum gets the check sum for this ROOT class
	//CheckSum() int

	// Members returns the list of members for this ROOT class
	Members() []Member

	// Version returns the version number for this ROOT class
	Version() int

	// Name returns the ROOT class name for this ROOT class
	Name() string
}

// Member represents a single member of a ROOT class
type Member interface {
	// GetArrayDim returns the dimension of the array (if any)
	ArrayDim() int

	// GetComment returns the comment associated with this member
	Comment() string

	// Name returns the name of this member
	Name() string

	// Type returns the class of this member
	Type() Class

	// GetValue returns the value of this member
	//GetValue(o Object) reflect.Value
}

// Object represents a ROOT object
type Object interface {
	// Class returns the ROOT class of this object
	Class() string

	// Name returns the name of this ROOT object
	Name() string

	// Title returns the title of this ROOT object
	Title() string
}

// ClassFactory creates ROOT classes
type ClassFactory interface {
	Create(name string) Class
}

// ROOTStreamer is the interface describing data that provides its own
// routine for encoding and decoding transmitted values sent to a ROOT file.
type ROOTStreamer interface {
	// de-serialize into the current value using 'buf' as input
	ROOTDecode(buf *Buffer) error
	// serialize the current value using 'buf' as output
	ROOTEncode(buf *Buffer) error
}

// FileSetter is the interface allowing to reset the ownership of a ROOT object
// w.r.t a given file
type FileSetter interface {
	SetFile(f *File) error
}

// EOF
