package proxy

type EmbroideryService interface {
    Embroider(size string)
    EmbroiderCustomized(size, requirements string)
}
