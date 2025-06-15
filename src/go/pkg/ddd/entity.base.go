package ddd

type Entity[ID comparable] interface {
	ID() ID
	Equals(Entity[ID]) bool
}

type BaseEntity[ID comparable] struct {
	Id ID `json:"id"`
}

func NewBaseEntity[ID comparable](id ID) *BaseEntity[ID] {
	return &BaseEntity[ID]{
		Id: id,
	}
}

func (e *BaseEntity[ID]) ID() ID {
	return e.Id
}

func (e *BaseEntity[ID]) Equals(other Entity[ID]) bool {
	return e.ID() == other.ID()
}
