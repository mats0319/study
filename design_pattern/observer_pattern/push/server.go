package mario

type ServerI interface {
    Notify()
    Attach(ClientI)
    Detach(ClientI)
}
