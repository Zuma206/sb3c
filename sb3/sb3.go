package sb3

import (
	"archive/zip"
	"encoding/json"
	"io"

	"github.com/zuma206/sb3c/utils"
)

// Represents a .sb3 file
type SB3 struct {
	project Project
}

func NewSB3() *SB3 {
	return &SB3{
		project: Project{
			Targets: []Target{
				{
					IsStage: true,
					Name:    "Stage",
					Costumes: []Costume{
						{
							Asset: Asset{
								Name:       "backdrop1",
								AssetId:    "cd21514d0531fdffb22204e0ec5ed84a",
								DataFormat: "svg",
							},
						},
					},
					Sounds: []struct{}{},
				},
			},
			Meta: Meta{
				Semver: "3.0.0",
			},
		},
	}
}

func (sb3 *SB3) WriteTo(w io.Writer) (int64, error) {
	counter := utils.NewCounter(w)
	sb3File := zip.NewWriter(counter)
	defer sb3File.Close()
	projectFile, err := sb3File.Create("project.json")
	if err != nil {
		return int64(counter.Written), err
	}
	projectEncoder := json.NewEncoder(projectFile)
	if err := projectEncoder.Encode(sb3.project); err != nil {
		return int64(counter.Written), err
	}
	return int64(counter.Written), nil
}
