package groot

// Class represents a ROOT class.
// Class instances are created by a ClassFactory.
type Class interface {
	// GetCheckSum gets the check sum for this ROOT class
	GetCheckSum() int

	// GetMembers returns the list of members for this ROOT class
	GetMembers() []Member

	// GetVersion returns the version number for this ROOT class
	GetVersion() int

	// GetClassName returns the ROOT class name for this ROOT class
	GetClassName() string

	// GetSuperClasses returns the list of super-classes for this ROOT class
	GetSuperClasses() []Class
}

// Member represents a single member of a ROOT class
type Member interface {
	// GetArrayDim returns the dimension of the array (if any)
	GetArrayDim() int

	// GetComment returns the comment associated with this member
	GetComment() string

	// GetName returns the name of this member
	GetName() string

	// GetType returns the class of this member
	GetType() Class

	// GetValue returns the value of this member
	//GetValue(o Object) reflect.Value
}

// Object represents a ROOT object
type Object interface {
	// GetClass returns the ROOT class of this object
	GetClass() Class
}

// ClassFactory creates ROOT classes
type ClassFactory interface {
	Create(name string) Class
}

// EOF
