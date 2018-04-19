package generate

import (
	"fmt"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/segmentio/objconv/objutil"
	"golang.org/x/tools/go/loader"
)

type kind int

const (
	typBasic kind = iota
	typArray
	typStruct
	typAlias
	typTime
)

type generatorType struct {
	kind kind

	name string

	pkgName string
	pkgPath string

	// the type aliased to, or the array element type
	subType *generatorType

	basicKind types.BasicKind

	arrayLength int64

	structFields map[structKey]*generatorType
}

func (g *generatorType) funcName() string {
	switch g.kind {
	case typArray:
		return g.subType.funcName() + "_" + strconv.FormatInt(g.arrayLength, 10)
	case typStruct:
		return g.pkgName + "_" + g.name
	case typAlias:
		return g.pkgName + "_" + g.name
	case typTime:
		return "time_Time"
	case typBasic:
		return g.name
	}
	panic(g)
}

func (g *generatorType) typeName(basePkg string) string {
	prefix := ""
	if g.pkgPath != "" && g.pkgPath != basePkg {
		prefix = g.pkgName + "."
	}
	switch g.kind {
	case typArray:
		return fmt.Sprintf("[%d]%s", g.arrayLength, g.subType.typeName(basePkg))
	case typStruct:
		return prefix + g.name
	case typAlias:
		return prefix + g.name
	case typTime:
		return "time.Time"
	case typBasic:
		return g.name
	}
	panic(g)
}

type structKey struct {
	encodedName string
	fieldName   string
}

func createGeneratorType(path string, typ string) (*generatorType, error) {
	config := loader.Config{}

	config.Import(path)
	program, err := config.Load()
	if err != nil {
		return nil, err
	}

	mainPkg := program.Package(path)

	var typeObj types.Object
	for k, def := range mainPkg.Defs {
		if k.Name == typ {
			typeObj = def
			break
		}
	}
	if typeObj == nil {
		return nil, fmt.Errorf("Could not find type %s in %s", typ, path)
	}

	var walkType func(types.Type) *generatorType

	walkType = func(t types.Type) *generatorType {
		switch typ := t.(type) {
		case *types.Basic:
			return &generatorType{
				kind:      typBasic,
				basicKind: typ.Kind(),
				name:      typ.Name(),
			}
		case *types.Named:
			pkg := typ.Obj().Pkg()
			if st, ok := typ.Underlying().(*types.Struct); ok {
				if pkg.Name()+"."+typ.Obj().Name() == "time.Time" {
					return &generatorType{
						kind: typTime,
					}
				}

				fields := make(map[structKey]*generatorType, st.NumFields())
				for i := 0; i < st.NumFields(); i++ {
					f := st.Field(i)
					fields[structKey{
						fieldName:   f.Name(),
						encodedName: f.Name(), // TODO: support json,objconv tags
					}] = walkType(f.Type())
				}

				return &generatorType{
					kind:         typStruct,
					name:         typ.Obj().Name(),
					pkgName:      pkg.Name(),
					pkgPath:      pkg.Path(),
					structFields: fields,
				}
			}
			// otherwise its an alias
			return &generatorType{
				kind:    typAlias,
				name:    typ.Obj().Name(),
				pkgName: pkg.Name(),
				pkgPath: pkg.Path(),
				subType: walkType(typ.Underlying()),
			}
		case *types.Array:
			return &generatorType{
				kind:        typArray,
				subType:     walkType(typ.Elem()),
				arrayLength: typ.Len(),
			}
		}

		panic("UNKNOWN " + reflect.TypeOf(t).String())
	}

	return walkType(typeObj.Type()), nil
}

func GenerateDecode(path string, typ string) {
	gtyp, err := createGeneratorType(path, typ)
	if err != nil {
		panic(err)
	}

	pkgs := map[string]struct{}{}
	typs := map[string]*generatorType{}

	var walk func(*generatorType)

	walk = func(g *generatorType) {
		typs[g.funcName()] = g
		switch g.kind {
		case typArray:
			walk(g.subType)
		case typStruct:
			pkgs[g.pkgPath] = struct{}{}
			for _, f := range g.structFields {
				walk(f)
			}
		case typAlias:
			pkgs[g.pkgPath] = struct{}{}
			walk(g.subType)
		case typTime:
			pkgs["time"] = struct{}{}
		}
	}

	walk(gtyp)

	decoderName := gtyp.name + "Decoder"

	pkgNames := make([]string, 0, len(pkgs))

	for p := range pkgs {
		if p == gtyp.pkgPath {
			continue
		}
		pkgNames = append(pkgNames, `"`+p+`"`)
	}

	res := fmt.Sprintf(`package %[1]s

import (
	"reflect"
	"unsafe"

	%[2]s

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/generate/util"
)

type %[3]s struct {
	Parser objconv.Parser
}

func (d *%[3]s) Decode() (%[4]s, error) {
	return d.decode_%[5]s()
}`, gtyp.pkgName, strings.Join(pkgNames, "\n"), decoderName, gtyp.name, gtyp.funcName())

	for _, t := range typs {
		res += "\n\n" + generateDecodeFunc(decoderName, path, t)
	}

	fmt.Println(res)
}

func generateDecodeFunc(decoderName string, basePkg string, g *generatorType) string {
	res := fmt.Sprintf("func (d *%[1]s) decode_%[2]s() (%[3]s, error) {", decoderName, g.funcName(), g.typeName(basePkg))

	switch g.kind {
	case typBasic:
		res += generateBasicFunc(g)
	case typArray:
		res += generateArrayFunc(basePkg, g)
	case typStruct:
		res += generateStructFunc(basePkg, g)
	case typAlias:
		res += generateAliasFunc(basePkg, g)
	case typTime:
		res += generateTimeFunc()
	}

	res += "\n}"

	return res
}

func generateBasicFunc(g *generatorType) string {
	switch g.basicKind {
	case types.Bool:
		return `
	return util.ParseBool(d.Parser)`

	case types.Int:
		return `
	v, err := util.ParseInt(d.Parser)
	return int(v), err`

	case types.Int8:
		return fmt.Sprintf(`
	v, err := util.ParseInt(d.Parser)
	if v > %d || v < %d {
		err = fmt.Errorf("%%d does not fit in int8", v)
	}
	return int8(v), err`, objutil.Int8Max, objutil.Int8Min)

	case types.Int16:
		return fmt.Sprintf(`
	v, err := util.ParseInt(d.Parser)
	if v > %d || v < %d {
		err = fmt.Errorf("%%d does not fit in int16", v)
	}
	return int16(v), err`, objutil.Int16Max, objutil.Int16Min)

	case types.Int32:
		return fmt.Sprintf(`
	v, err := util.ParseInt(d.Parser)
	if v > %d || v < %d {
		err = fmt.Errorf("%%d does not fit in int32", v)
	}
	return int32(v), err`, objutil.Int32Max, objutil.Int32Min)

	case types.Int64:
		return `
	return util.ParseInt(d.Parser)`

	case types.Uint:
		return `
	v, err := util.ParseUint(d.Parser)
	return uint(v), err`

	case types.Uint8:
		return fmt.Sprintf(`
	v, err := util.ParseUint(d.Parser)
	if v > %d {
		err = fmt.Errorf("%%d does not fit in uint8", v)
	}
	return uint8(v), err`, objutil.Uint8Max)

	case types.Uint16:
		return fmt.Sprintf(`
	v, err := util.ParseUint(d.Parser)
	if v > %d {
		err = fmt.Errorf("%%d does not fit in uint16", v)
	}
	return uint16(v), err`, objutil.Uint16Max)

	case types.Uint32:
		return fmt.Sprintf(`
	v, err := util.ParseUint(d.Parser)
	if v > %d {
		err = fmt.Errorf("%%d does not fit in uint32", v)
	}
	return uint32(v), err`, objutil.Uint32Max)

	case types.Uint64:
		return `
	return util.ParseUint(d.Parser)`

	case types.Float32:
		return `
	v, err := util.ParseFloat(d.Parser)
	return float32(v), err`

	case types.Float64:
		return `
	return util.ParseFloat(d.Parser)`

	case types.String:
		return `
	v, err := util.ParseString(d.Parser)
	return string(v), err`

	default:
		panic(fmt.Errorf("Unknown basic kind %v (%s)", g.basicKind, g.name))
	}
}

func generateArrayFunc(basePkg string, g *generatorType) string {
	res := fmt.Sprintf(`
	var res %[1]s
	typ, err := d.Parser.ParseType()
	if err != nil {
		return res, err
	}
	if typ != objconv.Array {
		return res, fmt.Errorf("Cannot decode value of type %%s into %[1]s", typ.String())
	}
	count, err := d.Parser.ParseArrayBegin()
	if err != nil {
		return res, err
	}
	if count >= 0 && count != %[2]d {
		return res, fmt.Errorf("Cannot decode array of length %%d into %[1]s", count)
	}
	if count < 0 {
		if err = d.Parser.ParseArrayNext(0); err != nil {
			return res, err
		}
	}`, g.typeName(basePkg), g.arrayLength)
	for i := 0; i < int(g.arrayLength); i++ {
		res += fmt.Sprintf(`
	res[%[1]d], err = d.decode_%[2]s()
	if err != nil {
		return res, err
	}`, i, g.subType.funcName())
	}
	res += fmt.Sprintf(`
	if count < 0 {
		if err = d.Parser.ParseArrayNext(%[1]d); err == nil {
			return res, errors.New("Cannot decode array greater than length %[1]d into %[2]s")
		} else if err != objconv.End {
			return res, err
		}
	}
	return res, d.Parser.ParseArrayEnd(%[1]d)`, g.arrayLength, g.typeName(basePkg))
	return res
}

func generateStructFunc(basePkg string, g *generatorType) string {
	res := fmt.Sprintf(`
	var res %[1]s
	typ, err := d.Parser.ParseType()
	if err != nil {
		return res, err
	}
	if typ != objconv.Map {
		return res, fmt.Errorf("Cannot decode value of type %%s into %[1]s", typ.String())
	}
	count, err := d.Parser.ParseMapBegin()
	if err != nil {
		return res, err
	}
	var i int
	for i = 0; count < 0 || i < count; i++ {
		if count < 0 || i != 0 {
			if err = d.Parser.ParseMapNext(i); err == objconv.End {
				break
			} else if err != nil {
				return res, err
			}
		}
		ktyp, err := d.Parser.ParseType()
		if err != nil {
			return res, err
		}
		if ktyp != objconv.String {
			return res, fmt.Errorf("Cannot decode %[1]s key from type %%s", ktyp.String())
		}
		kb, err := d.Parser.ParseString()
		if err != nil {
			return res, err
		}
		k := *(*string)(unsafe.Pointer(&kb))
		switch k {`, g.typeName(basePkg))
	for f, ftyp := range g.structFields {
		res += fmt.Sprintf(`
		case "%[1]s":
			if err = d.Parser.ParseMapValue(i); err != nil {
				return res, err
			}
			res.%[2]s, err = d.decode_%[3]s()
			if err != nil {
				return res, err
			}`, f.encodedName, f.fieldName, ftyp.funcName())
	}
	res += `
		default:
			if err = util.SkipValue(d.Parser); err != nil {
				return res, err
			}
		}
	}
	err = d.Parser.ParseMapEnd(i)
	return res, err`
	return res
}

func generateAliasFunc(basePkg string, g *generatorType) string {
	return fmt.Sprintf(`
	v, err := d.decode_%s()
	return %s(v), err`, g.subType.funcName(), g.typeName(basePkg))
}

func generateTimeFunc() string {
	return `
	typ, err := d.Parser.ParseType()
	if err != nil {
		return time.Time{}, err
	}
	switch typ {
	case objconv.Time:
		return d.Parser.ParseTime()
	case objconv.String:
		b, err := d.Parser.ParseString()
		if err != nil {
			return time.Time{}, err
		}
		tstr := *(*string)(unsafe.Pointer(&b))
		return time.Parse(time.RFC3339Nano, tstr)
	default:
		return time.Time{}, fmt.Errorf("Cannot parse Time from type %s", typ.String())
	}`
}
