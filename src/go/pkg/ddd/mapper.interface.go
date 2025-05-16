package ddd

type PersistentObject[ID comparable] interface {
	ToDomain() Entity[ID]
	FromDomain(entity Entity[ID])
}
