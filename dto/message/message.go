package message

const (
	ErrBind           = "Invalid request payload"
	ErrRegister       = "Failed to register user"
	ErrRefresh        = "Failed to refresh access token"
	ErrGetCurrentUser = "Failed to retrieve user identity"
	ErrGetUserTickets = "Failed to retrieve user tickets"
	ErrGetTicket      = "Failed to retrieve ticket"
	ErrGetTickets     = "Failed to retrieve all tickets"
	ErrGetFormFile    = "Failed to retrieve uploaded file"
	ErrCreateTicket   = "Failed to create ticket"
	ErrDeleteTicket   = "Failed to delete ticket"
	ErrUpdateTicket   = "Failed to update ticket"
)

const (
	SuccessRegister       = "User registered successfully"
	SuccessRefresh        = "Access token refreshed successfully"
	SuccessGetCurrentUser = "User identity retrieved successfully"
	SuccessGetUserTickets = "User tickets retrieved successfully"
	SuccessGetTicket      = "Ticket retrieved successfully"
	SuccessGetTickets     = "All tickets retrieved successfully"
	SuccessCreateTicket   = "Ticket created successfully"
	SuccessDeleteTicket   = "Ticket deleted successfully"
	SuccessUpdateTicket   = "Ticket updated successfully"
)
