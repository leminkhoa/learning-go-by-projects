package main

// Accounts
// Products

type queryResolver struct {
	server *Server,
}


func (r *queryResolver) Accounts(
	ctx context.Context, 
	pagination *PaginationInput, 
	id *string
) ([]*Account ,error) {

}

