package sb3

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"

	"github.com/zuma206/sb3c/utils"
)

// Represents a .sb3 file
type SB3 struct {
	stage   *TargetHnd
	project Project
}

func NewSB3() *SB3 {
	return &SB3{
		project: Project{
			Targets: []*Target{},
			Meta: Meta{
				Semver: "3.0.0",
			},
		},
	}
}

var MissingStageError = errors.New("missing stage")

func (sb3 *SB3) Validate() error {
	if sb3.stage == nil {
		return MissingStageError
	}
	return nil
}

func (sb3 *SB3) WriteTo(w io.Writer) (int64, error) {
	if err := sb3.Validate(); err != nil {
		return 0, err
	}
	counter := utils.NewCounter(w)
	sb3File := zip.NewWriter(counter)
	defer sb3File.Close()
	if err := sb3.writeProjectJson(sb3File); err != nil {
		return int64(counter.Written), err
	}
	return int64(counter.Written), nil
}

func (sb3 *SB3) writeProjectJson(sb3File *zip.Writer) error {
	projectFile, err := sb3File.Create("project.json")
	if err != nil {
		return err
	}
	projectEncoder := json.NewEncoder(projectFile)
	if err := projectEncoder.Encode(sb3.project); err != nil {
		return err
	}
	return nil
}

func (sb3 *SB3) newTarget(name string, isStage bool) *Target {
	target := &Target{
		IsStage: isStage,
		Name:    name,
		Costumes: []Costume{
			{
				Asset: Asset{
					Name:       "backdrop1",
					DataFormat: "svg",
					AssetId:    "cd21514d0531fdffb22204e0ec5ed84a",
					Md5ext:     "cd21514d0531fdffb22204e0ec5ed84a.svg",
				},
			},
		},
		Sounds:    []struct{}{},
		Blocks:    map[string]*Block{},
		Variables: map[string]*Variable{},
	}
	sb3.project.Targets = append(sb3.project.Targets, target)
	return target
}

var StageAlreadyExistsError = errors.New("stage already exists")

func (sb3 *SB3) NewStage() (*TargetHnd, error) {
	if sb3.stage != nil {
		return nil, StageAlreadyExistsError
	}
	sb3.stage = &TargetHnd{
		target: sb3.newTarget("Stage", true),
		sb3:    sb3,
	}
	return sb3.stage, nil
}
