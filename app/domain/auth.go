package domain

type Auth interface {
	Auth(token string) (User, error)
}

type auth struct {
	endpoint     string
	communicator Communicator
}

func NewAuth(ep string, c Communicator) Auth {
	return auth{
		endpoint:     ep,
		communicator: c,
	}

}

func (au auth) Auth(token string) (User, error) {
	return User{}, nil
}
