package golang_mysql

import "github.com/google/uuid"

func GenerateCustomUUID() string {
	u := uuid.New()
	uuidStr := u.String()

	// Membuang tanda "-" dari UUID dan menghasilkan 36 karakter
	uuidStr = uuidStr[:8] + uuidStr[9:13] + uuidStr[14:18] + uuidStr[19:23] + uuidStr[24:]

	return uuidStr
}
