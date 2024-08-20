package helpers

import "github.com/google/uuid"

func ParseStringToUuid(uuidAsString string) (uuid.UUID, error) {
	var returnUuid uuid.UUID
	if parsedUuid, err := uuid.Parse(uuidAsString); err == nil {
		returnUuid = parsedUuid
	} else {
		return uuid.UUID{}, err
	}

	return returnUuid, nil
}
