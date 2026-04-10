package sb3

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/zuma206/sb3c/utils"
)

// Represents a .sb3 file
type SB3 struct {
	assets  []StoredAsset
	stage   *TargetHnd
	project Project
	fs      fs.FS
}

type StoredAsset struct {
	*Asset
	content []byte
}

func NewSB3(fileSystem fs.FS) *SB3 {
	return &SB3{
		project: Project{
			Targets: []*Target{},
			Meta: Meta{
				Semver: "3.0.0",
			},
		},
		fs: fileSystem,
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
	if err := sb3.writeAssetFiles(sb3File); err != nil {
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

func (sb3 *SB3) writeAssetFiles(sb3File *zip.Writer) error {
	for _, asset := range sb3.assets {
		file, err := sb3File.Create(asset.Md5ext)
		if err != nil {
			return err
		}
		if _, err := file.Write(asset.content); err != nil {
			return err
		}
	}
	return nil
}

func (sb3 *SB3) newTarget(name string, isStage bool) *Target {
	target := &Target{
		IsStage:   isStage,
		Name:      name,
		Costumes:  []*Costume{},
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

var InvalidCostumePathErr = errors.New("invalid costume path")

func (sb3 *SB3) loadAsset(name string, path string) (*Asset, error) {
	content, err := fs.ReadFile(sb3.fs, path)
	if err != nil {
		return nil, err
	}
	sum := md5.Sum(content)
	hexSum := hex.EncodeToString(sum[:])
	pathParts := strings.Split(path, ".")
	if len(pathParts) < 2 {
		err := fmt.Errorf("%q has no file extension", path)
		return nil, errors.Join(InvalidCostumePathErr, err)
	}
	format := pathParts[len(pathParts)-1]
	asset := &Asset{
		Name:       name,
		AssetId:    hexSum,
		DataFormat: format,
		Md5ext:     hexSum + "." + format,
	}
	sb3.assets = append(sb3.assets, StoredAsset{asset, content})
	return asset, nil
}
