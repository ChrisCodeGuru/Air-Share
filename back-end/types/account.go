package types

type AccountsPermission []AccountPermission

type AccountPermission struct {
	ID         string
	Email      string
	Permission int
}
