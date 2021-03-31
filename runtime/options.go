package runtime

type option func(*Service)

type HandlerT int

const (
	HTTP HandlerT = iota + 1
	Lambda
)

func Handler(handlerType HandlerT) option {
	return func(s *Service) {}
}

func HotReload(dirs ...string) option {
	return func(s *Service) {}
}

type RestartPolicyT int

const (
	Forever RestartPolicyT = iota + 1
	Never
)

func RestartPolicy(restartPolicy RestartPolicyT) option {
	return func(s *Service) {}
}

func LogLevel(level logging.LogLevel) {
	return func(s *Service) {}
}
