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
	IsStage   bool       `json:"isStage"`
	Name      string     `json:"name"`
	Blocks    struct{}   `json:"blocks"`
	Variables struct{}   `json:"variables"`
	Costumes  []Costume  `json:"costumes"`
	Sounds    []struct{} `json:"sounds"`
}

type Asset struct {
	Name       string `json:"name"`
	AssetId    string `json:"assetId"`
	DataFormat string `json:"dataFormat"`
}

type Costume struct {
	Asset
}
