package domain

import "errors"

var (
	ErrInternal 		= errors.New("internal error")
	ErrConflictingData 	= errors.New("conflicting data")
	ErrDataNotFound    	= errors.New("data not found")

)
