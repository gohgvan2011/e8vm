package asm8

import (
	"e8vm.io/e8vm/asm8/ast"
	"e8vm.io/e8vm/asm8/parse"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
)

type pkg struct {
	path  string
	files []*file

	imports *importDecl
}

func resolvePkg(p string, src map[string]*build8.File) (*pkg, []*lex8.Error) {
	log := lex8.NewErrorList()
	ret := new(pkg)
	ret.path = p

	asts := make(map[string]*ast.File)

	// parse all the files first
	var parseErrs []*lex8.Error
	for name, f := range src {
		astFile, es := parse.File(f.Path, f)
		if es != nil {
			parseErrs = append(parseErrs, es...)
		}
		asts[name] = astFile
	}
	if len(parseErrs) > 0 {
		return nil, parseErrs
	}

	for name, astFile := range asts {
		// then resolve the file
		file := resolveFile(log, astFile)
		ret.files = append(ret.files, file)

		// enforce import policy
		if len(src) == 1 || name == "import.s" {
			if ret.imports != nil {
				log.Errorf(file.imports.Kw.Pos,
					"double valid import stmt; two import.s?",
				)
			} else {
				ret.imports = file.imports
			}
		} else if file.imports != nil {
			log.Errorf(file.imports.Kw.Pos,
				"invalid import outside import.s in a multi-file package",
			)
		}
	}

	if es := log.Errs(); es != nil {
		return nil, es
	}
	return ret, nil
}
