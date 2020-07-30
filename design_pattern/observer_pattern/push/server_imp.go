package mario

type Server struct {
    ServerStatus string
    clients      []ClientI
}

var _ ServerI = (*Server)(nil)

func (s *Server) Notify() {
    cs := s.clients

    for i := range cs {
        cs[i].Update()
    }
}

func (s *Server) Attach(c ClientI) {
    s.clients = append(s.clients, c)
}

func (s *Server) Detach(c ClientI) {
    cs := s.clients

    for i := range cs {
        if c.GetName() == cs[i].GetName() {
            cs = append(cs[:i], cs[i+1:]...)
            break
        }
    }

    s.clients = cs
}
