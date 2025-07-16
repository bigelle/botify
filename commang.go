package botify

type command struct {
	Name        string
	Description string
	Handler     HandlerFunc
}

type scopeKey struct {
	Scope  string
	ChatID string
	UserID int
}

type commandRegistry struct {
	byScope map[scopeKey]map[string]command

	byCommand map[string]struct {
		Handler HandlerFunc
		Scopes  map[scopeKey]struct{}
	}
}

func (r *commandRegistry) GetCommands(scope scopeKey) []command {
	if r.byScope == nil {
		return nil
	}

	cmds := make([]command, 0, len(r.byScope[scope]))
	for _, cmd := range r.byScope[scope] {
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (r *commandRegistry) GetScopes() []scopeKey {
	if r.byScope == nil {
		return nil
	}

	scopes := make([]scopeKey, 0, len(r.byScope))

	for scope, val := range r.byScope {
		if len(val) != 0 && val != nil {
			scopes = append(scopes, scope)
		}
	}

	return scopes
}

func (r *commandRegistry) GetHandler(name string) (HandlerFunc, bool) {
	if r.byCommand == nil {
		return nil, false
	}

	cmd, ok := r.byCommand[name]
	if !ok {
		return nil, false
	}

	return cmd.Handler, true
}

func (r *commandRegistry) AddCommand(cmd, desc string, handler HandlerFunc, scopes ...scopeKey) {
	if len(scopes) == 0 {
		scopes = []scopeKey{{Scope: "default"}}
	}

	if r.byCommand == nil {
		r.byCommand = make(map[string]struct {
			Handler HandlerFunc
			Scopes  map[scopeKey]struct{}
		})
	}
	if r.byScope == nil {
		r.byScope = make(map[scopeKey]map[string]command)
	}

	var cmdInfo struct {
		Handler HandlerFunc
		Scopes  map[scopeKey]struct{}
	}

	if existing, ok := r.byCommand[cmd]; ok {
		cmdInfo = existing
	} else {
		cmdInfo = struct {
			Handler HandlerFunc
			Scopes  map[scopeKey]struct{}
		}{
			Handler: handler,
			Scopes:  make(map[scopeKey]struct{}),
		}
	}

	for _, scope := range scopes {
		if r.byScope[scope] == nil {
			r.byScope[scope] = make(map[string]command)
		}
		r.byScope[scope][cmd] = command{
			Name:        cmd,
			Description: desc,
			Handler:     handler,
		}

		cmdInfo.Scopes[scope] = struct{}{}
	}

	r.byCommand[cmd] = cmdInfo
}
