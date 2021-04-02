package runtime

import "github.com/foldsh/fold/logging"

type Option func(*Runtime)

type HandlerT uint8

func WithEnv(env map[string]string) Option {
	return func(r *Runtime) {}
}

func WithSupervisor(supervisor Supervisor) Option {
	return func(r *Runtime) {
		r.supervisor = supervisor
	}
}

func WithClient(client Client) Option {
	return func(r *Runtime) {
		r.client = client
	}
}

func WithSocketFactory(socketFactory SocketFactory) Option {
	return func(r *Runtime) {
		r.socketFactory = socketFactory
	}
}

func WithRouterFactory(routerFactory RouterFactory) Option {
	return func(r *Runtime) {
		r.routerFactory = routerFactory
	}
}

func WatchDirs(dirs ...string) Option {
	return func(r *Runtime) {}
}

type CrashPolicyT uint8

const (
	EXIT CrashPolicyT = iota + 1
	KEEP_ALIVE
)

func CrashPolicy(crashPolicy CrashPolicyT) Option {
	return func(r *Runtime) {
	}
}

func LogLevel(level logging.LogLevel) Option {
	return func(r *Runtime) {}
}
