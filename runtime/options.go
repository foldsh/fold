package runtime

type option func(*Runtime)

type HandlerT int

const (
	HTTP HandlerT = iota + 1
	Lambda
)

func Handler(handlerType HandlerT) option {
	return func(r *Runtime) {}
}

func HotReload(dirs ...string) option {
	return func(r *Runtime) {}
}

type RestartPolicyT int

const (
	Forever RestartPolicyT = iota + 1
	Never
)

func RestartPolicy(restartPolicy RestartPolicyT) option {
	return func(r *Runtime) {
	}
}

func LogLevel(level logging.LogLevel) {
	return func(r *Runtime) {}
}
