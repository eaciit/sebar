package sebar

type SebarServer struct {
	Protocol, Address, Secret string
}

func (s *SebarServer) Start() error {
	return nil
}

func (s *SebarServer) Stop() error {
	return nil
}
