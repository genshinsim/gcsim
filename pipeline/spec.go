package main

import (
	"path"
	"path/filepath"

	"github.com/shizukayuki/excel-hk4e"
)

//go:generate go tool github.com/dmarkham/enumer -text -json -linecomment -type=Kind $GOFILE
type Kind int

const (
	KindNone      Kind = iota //
	KindArtifact              // artifact
	KindCharacter             // character
	KindMonster               // monster
	KindWeapon                // weapon
	KindInvalid               // invalid
)

type Compiled struct {
	Configuration []*Config                         `yaml:"configuration,omitempty"`
	ElementCoeff  []float64                         `yaml:"element_coeff,omitempty"`
	GrowCurveData map[excel.GrowCurveType][]float64 `yaml:"grow_curve_data,omitempty"`
	ICDGroup      map[string]AttackAttenuation      `yaml:"icd_group,omitempty"`
	Localization  map[string]Localization           `yaml:"localization,omitempty"`
}

func (c *Compiled) ClearRef() {
	for _, cfg := range c.Configuration {
		cfg.Artifact.ClearRef()
		cfg.Character.ClearRef()
		cfg.Monster.ClearRef()
		cfg.Weapon.ClearRef()
	}
}

type Config struct {
	Path      string           `yaml:"path,omitempty"`
	Use       string           `yaml:"use,omitempty"`
	Kind      Kind             `yaml:"kind,omitempty"`
	Name      string           `yaml:"name,omitempty"`
	Override  Override         `yaml:"override,omitempty"`
	Shortcuts []string         `yaml:"shortcuts,omitempty"`
	Combat    CombatInfo       `yaml:"combat,omitempty"`
	Talents   []map[string]any `yaml:"talents,omitempty"`
	Docs      Documentation    `yaml:"docs,omitempty"`
	Abilities []*Ability       `yaml:"abilities,omitempty"`

	Artifact   *ArtifactSpec    `yaml:"artifact,omitempty"`
	Character  *CharacterSpec   `yaml:"character,omitempty"`
	Monster    *MonsterSpec     `yaml:"monster,omitempty"`
	Weapon     *WeaponSpec      `yaml:"weapon,omitempty"`
	Attributes []*AttributeSpec `yaml:"attributes,omitempty"`
}

func (c *Config) Dir() string {
	return path.Dir(filepath.ToSlash(c.Path))
}

func (c *Config) SortKey() string {
	switch c.Kind {
	case KindArtifact:
		return c.Artifact.Model.Key
	case KindCharacter:
		return c.Character.Model.Key
	case KindMonster:
		return c.Monster.Model.Key
	case KindWeapon:
		return c.Weapon.Model.Key
	default:
		return c.Path
	}
}

func (c *Config) ClearSpec() {
	c.Artifact = nil
	c.Character = nil
	c.Monster = nil
	c.Weapon = nil
	c.Attributes = nil
}

type Override struct {
	Id       uint32                       `yaml:"id,omitempty"`
	Depot    uint32                       `yaml:"depot,omitempty"`
	ICDGroup map[string]AttackAttenuation `yaml:"icd_group,omitempty"`
}

type AttackAttenuation struct {
	Name       string    `yaml:"name,omitempty"`
	Damage     []float64 `yaml:"damage,flow,omitempty"`
	Durability []float64 `yaml:"durability,flow,omitempty"`
	Timer      float64   `yaml:"timer,omitempty"`
}

type CombatInfo struct {
	ICDTags   []string `yaml:"icd_tags,omitempty"`
	ICDGroups []string `yaml:"icd_groups,omitempty"`
}

type Ability struct {
	Name   string            `yaml:"name,omitempty"`
	Note   string            `yaml:"note,omitempty"`
	Params map[string]string `yaml:"params,omitempty"`
	Hitbox []Hitbox          `yaml:"hitbox,omitempty"`
}

type Hitbox struct {
	Name     string    `yaml:"name,omitempty"`
	Note     string    `yaml:"note,omitempty"`
	Shape    string    `yaml:"shape,omitempty"`
	Center   string    `yaml:"center,omitempty"`
	Offset   []float64 `yaml:"offset,flow,omitempty"`
	Box      []float64 `yaml:"box,flow,omitempty"`
	FanAngle float64   `yaml:"fan_angle,omitempty"`
	Radius   float64   `yaml:"radius,omitempty"`
	Hitlag   Hitlag    `yaml:"hitlag,flow,omitempty"`
}

type Hitlag struct {
	Time        float64 `yaml:"time,omitempty"`
	Scale       float64 `yaml:"scale,omitempty"`
	DefenseHalt bool    `yaml:"defense_halt,omitempty"`
	Deployable  bool    `yaml:"deployable,omitempty"`
}

func (h Hitlag) IsZero() bool {
	if h.Deployable {
		return false
	}
	return !h.HasHitlag()
}

func (h Hitlag) HasHitlag() bool {
	return h.Scale != 1 && (h.Time != 0 || h.DefenseHalt)
}

type Documentation struct {
	Issues []string          `yaml:"issues,omitempty"`
	Fields map[string]string `yaml:"fields,omitempty"`
	Frames []Frames          `yaml:"frames,omitempty"`
}

type Frames struct {
	Name        string `yaml:"name,omitempty"`
	Count       string `yaml:"count,omitempty"`
	CountCredit string `yaml:"count_credit,omitempty"`
	Video       string `yaml:"video,omitempty"`
	VideoCredit string `yaml:"video_credit,omitempty"`
}

type Localization struct {
	ArtifactNames  map[string]string `json:"artifact_names,omitempty"  yaml:"artifact_names,omitempty"`
	CharacterNames map[string]string `json:"character_names,omitempty" yaml:"character_names,omitempty"`
	MonsterNames   map[string]string `json:"monster_names,omitempty"   yaml:"monster_names,omitempty"`
	WeaponNames    map[string]string `json:"weapon_names,omitempty"    yaml:"weapon_names,omitempty"`
}
