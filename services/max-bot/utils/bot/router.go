package bot

type HandlerRegistration struct {
	Handler HandlerFunc
	Filter  FilterFunc
}

type Router struct {
	handlers     map[string][]HandlerRegistration
	middlewres   []MiddlewareFunc
	childRouters []*Router
}

func NewRouter() *Router {
	return &Router{
		handlers:     make(map[string][]HandlerRegistration),
		middlewres:   make([]MiddlewareFunc, 0),
		childRouters: make([]*Router, 0),
	}
}

func (r *Router) IncludeRouter(child *Router) {
	r.childRouters = append(r.childRouters, child)
}

func (r *Router) registerHandler(updateType string, handler HandlerFunc, filter FilterFunc) {
	registration := HandlerRegistration{
		Handler: handler,
		Filter: filter,
	}

	r.handlers[updateType] = append(r.handlers[updateType], registration)
}

func (r *Router) Message(handler HandlerFunc, filter FilterFunc) {
	r.registerHandler("message", handler, filter)
}

func (r *Router) CallbackQuery(handler HandlerFunc, filter FilterFunc) {
	r.registerHandler("callback_query", handler, filter)
}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.middlewres = append(r.middlewres, middleware)
}

func (r *Router) getHandlers(updateType string) []HandlerRegistration {
	allHandlers := []HandlerRegistration{}

	if handlers, ok := r.handlers[updateType]; ok {
		allHandlers = append(allHandlers, handlers...)
	}

	for _, child := range r.childRouters {
		allHandlers = append(allHandlers, child.getHandlers(updateType)...)
	}

	for i := len(r.middlewres) - 1; i >=0 ; i-- {
		for j := range allHandlers {
			allHandlers[j].Handler = r.middlewres[i](allHandlers[j].Handler)
		}
	}

	return allHandlers
}
