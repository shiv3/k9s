package render

const (
	// NonResource represents a custom resource.
	NonResource = "*"
)

const (
	// Terminating represents a pod terminating status.
	Terminating = "Terminating"

	// Running represents a pod running status.
	Running = "Running"

	// Initialized represents a pod intialized status.
	Initialized = "Initialized"

	// Completed represents a pod completed status.
	Completed = "Completed"

	// ContainerCreating represents a pod container status.
	ContainerCreating = "ContainerCreating"

	// PodInitializing represents a pod initializing status.
	PodInitializing = "PodInitializing"
)

const (
	// MissingValue indicates an unset value.
	MissingValue = "<none>"

	// NAValue indicates a value that does not pertain.
	NAValue = "n/a"

	// UnknownValue represents an unknown.
	UnknownValue = "<unknown>"
)
