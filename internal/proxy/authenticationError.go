package proxy

type AuthenticationError struct {
	Message string
}

// ===========================================
// Constructors
// ===========================================

func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		Message: message,
	}
}

// ===========================================
// Public Methods
// ===========================================

func (e *AuthenticationError) Error() string {
	return e.Message
}
