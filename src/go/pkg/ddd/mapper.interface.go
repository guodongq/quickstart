package ddd

type PersistentEntity[ID comparable] interface {
	ToDomain() Entity[ID]
	FromDomain(entity Entity[ID])
}
