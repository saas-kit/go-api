package hash

import "golang.org/x/crypto/bcrypt"

//
const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into New func
)

// New hash from a string
func New(str string, cost ...int) (string, error) {
	cst := DefaultCost
	if len(cost) > 0 {
		cst = cost[0]
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(str), cst)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare hash with string
func Compare(hash, str string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}
