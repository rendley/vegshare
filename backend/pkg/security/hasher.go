package security

type PasswordHasher interface {
	Hash(password string) (string, error)
	Check(hashedPassword, inputPassword string) bool
}
