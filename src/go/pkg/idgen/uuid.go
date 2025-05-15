package idgen

import "github.com/google/uuid"

type UUIDGenerator struct {
	id string
}

func MustUUIDGenerator(idStr string) *UUIDGenerator {
	id := uuid.MustParse(idStr)
	return FromUUID(id)
}

func NewUUIDGenerator(idStr string) (*UUIDGenerator, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return FromUUID(id), nil
}

func FromUUID(uuid uuid.UUID) *UUIDGenerator {
	return &UUIDGenerator{
		id: uuid.String(),
	}
}

func (g *UUIDGenerator) Generate() string {
	return g.id
}
