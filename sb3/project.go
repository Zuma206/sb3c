package sb3

// Represents a project.json file
type Project struct {
	Meta    Meta      `json:"meta"`
	Targets []*Target `json:"targets"`
}

type Meta struct {
	Semver string `json:"semver"`
}

type Target struct {
	IsStage   bool              `json:"isStage"`
	Name      string            `json:"name"`
	Blocks    map[string]*Block `json:"blocks"`
	Variables struct{}          `json:"variables"`
	Costumes  []Costume         `json:"costumes"`
	Sounds    []struct{}        `json:"sounds"`
}

type Asset struct {
	Name       string `json:"name"`
	AssetId    string `json:"assetId"`
	DataFormat string `json:"dataFormat"`
}

type Costume struct {
	Asset
}

type Block struct {
	Opcode   string            `json:"opcode"`
	Next     Nullable[string]  `json:"next"`
	Parent   Nullable[string]  `json:"parent"`
	Inputs   map[string]*Input `json:"inputs"`
	Fields   struct{}          `json:"fields"`
	Shadow   bool              `json:"shadow"`
	TopLevel bool              `json:"topLevel"`
	X        int               `json:"x"`
	Y        int               `json:"y"`
	Mutation *Mutation         `json:"mutation,omitempty"`
}

type Mutation struct {
	TagName          string     `json:"tagName"`
	Children         []struct{} `json:"children"`
	Proccode         string     `json:"proccode"`
	Argumentids      string     `json:"argumentids"`
	Argumentnames    string     `json:"argumentnames"`
	Argumentdefaults string     `json:"argumentdefaults"`
	Warp             bool       `json:"warp"`
}
