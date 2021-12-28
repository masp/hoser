package ast

import (
	"github.com/masp/hoser/token"
)

// ModuleSet is a lazily loaded cache of modules that are resolved as needed when they are
// imported by a requesting module. Use Lookup() to fetch and parse a module if it hasn't been
// yet. This guarantees that each file is only parsed once.
type ModuleSet struct {
	Modules map[string]*CachedModule // full module name to module entry
}

type CachedModule struct {
	File *token.File // file holding module
	Mod  *Module     // nil if module is not loaded yet but was indexed
}

func (cm CachedModule) IsLoaded() bool {
	return len(cm.Mod.DefinedBlocks) > 0
}

func EmptyModuleSet() ModuleSet {
	return ModuleSet{
		Modules: make(map[string]*CachedModule),
	}
}

func (ms *ModuleSet) IndexFile(file *token.File, moduleHeader *Module) *CachedModule {
	fullName := moduleHeader.Name.Value
	if cached, ok := ms.Modules[fullName]; ok {
		return cached
	}
	ms.Modules[fullName] = &CachedModule{
		File: file,
		Mod:  moduleHeader,
	}
	return ms.Modules[fullName]
}

func (ms *ModuleSet) LoadModule(module *Module) *CachedModule {
	fullName := module.Name.Value
	ms.Modules[fullName].Mod = module
	return ms.Modules[fullName]
}
