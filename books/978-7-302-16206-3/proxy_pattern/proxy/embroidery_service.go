package mario

type EmbroideryService interface {
	Embroider(size string) string
	EmbroiderCustomized(customer, size, requirements string) string
}
