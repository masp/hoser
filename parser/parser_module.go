package parser

import (
	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

func (p *parser) parseModuleHeader() *ast.Module {
	p.eatOnly(token.Module)
	name := p.parseLiteral(p.eat())
	if name.Type != token.String {
		p.expectedError(name.Pos(), "module name as a quoted string")
	}

	return &ast.Module{
		Name: name,
	}
}

func (p *parser) parseModule() (module *ast.Module) {
	module = p.parseModuleHeader()
	for {
		p.eatAll(token.Semicolon)

		keyword := p.eat()
		switch keyword.tok {
		case token.Import:
			imp := p.parseImport()
			module.Imports = append(module.Imports, &imp)
		case token.Pipe:
			pipe := p.parsePipeBlock()
			module.DefinedBlocks = append(module.DefinedBlocks, &pipe)
		case token.Stub:
			stub := p.parseStubBlock()
			module.DefinedBlocks = append(module.DefinedBlocks, &stub)
		case token.Eof:
			return
		default:
			p.expectedError(keyword, "import/pipe/stub")
			return
		}
	}
}

func (p *parser) parseImport() (imp ast.ImportDecl) {
	importArg := p.eat()
	imp.ModuleName = p.parseLiteral(importArg)
	if imp.ModuleName.Type != token.String {
		p.expectedError(importArg, "string module path")
	}
	return
}
