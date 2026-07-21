package main

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/shizukayuki/excel-hk4e"
	"github.com/urfave/cli/v3"
)

var (
	optSources  []string
	optCacheDir string
	optRemote   bool
	optCache    bool
	optVerbose  bool
	optQuickfix []string

	// excel.Languages
	languages = map[string]string{
		"English":  "EN",
		"Chinese":  "CHS",
		"German":   "DE",
		"Spanish":  "ES",
		"Japanese": "JP",
		"Korean":   "KR",
		"Russian":  "RU",
	}
)

func main() {
	if err := app.Run(context.Background(), os.Args); err != nil {
		Log(slog.LevelError, "%s", err)
		os.Exit(1)
	}
}

var app = &cli.Command{
	Name:  "pipeline",
	Usage: "generate required files for gcsim to function",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "s",
			Usage:       "path to `datamine` directory or a remote repository in the form of 'github:<user>/<repo>/<ref>'",
			HideDefault: true,
			Value: []string{
				"github:iam-akuzihs/excel/live",
				filepath.Join(xdg.Home, "git", "GenshinData"),
				"github:DimbreathBot/AnimeGameData/main",
				"gitlab:Dimbreath/AnimeGameData2/main",
			},
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("DM_REPO"),
				cli.EnvVar("GENSHIN_DATA_REPO"),
			),
			Destination: &optSources,
		},
		&cli.StringFlag{
			Name:        "cache-dir",
			Value:       "./pipeline/_cache",
			Destination: &optCacheDir,
		},
		&cli.BoolWithInverseFlag{
			Name:        "remote",
			Usage:       "use sources that require a network connection",
			Value:       false,
			Destination: &optRemote,
		},
		&cli.BoolWithInverseFlag{
			Name:        "cache",
			Usage:       "use and cache files fetched from remote sources",
			Value:       true,
			Destination: &optCache,
		},
		&cli.BoolFlag{
			Name:        "v",
			Hidden:      true,
			Value:       false,
			Destination: &optVerbose,
		},
		&cli.StringSliceFlag{
			Name:        "qf",
			Usage:       "translate datamine fields at runtime",
			Destination: &optQuickfix,
		},
	},
	Before: func(_ context.Context, _ *cli.Command) (context.Context, error) {
		var err error
		projectRoot, err = os.OpenRoot(".")
		if err != nil {
			return nil, err
		}
		cacheRoot, err = os.OpenRoot(optCacheDir)
		if err != nil {
			return nil, err
		}
		for _, name := range []string{".git", "go.mod"} {
			if _, err := projectRoot.Stat(name); err != nil {
				return nil, errors.New("must be in root directory")
			}
		}
		return nil, nil
	},
	Commands: []*cli.Command{
		{
			Name:   "fetch",
			Usage:  "fetch and cache datamine from multiple sources",
			Action: fetch,
		},
		{
			Name:   "check",
			Usage:  "check for missing implementations",
			Action: check,
		},
		{
			Name:   "run",
			Usage:  "generate files",
			Action: run,
		},
		{
			Name:  "init",
			Usage: "print config.yml template",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "name"},
				&cli.IntFlag{Name: "id"},
				&cli.IntFlag{Name: "depot"},
			},
			Action: configTemplate,
		},
	},
}

func fetch(ctx context.Context, cmd *cli.Command) error {
	var inputs []FetchFunc
	for _, src := range optSources {
		var fn FetchFunc
		if s := strings.SplitN(src, ":", 2); optRemote && len(s) == 2 {
			proto := s[0]
			path := s[1]
			switch proto {
			case "github", "gitlab":
				s := strings.Split(path, "/")
				if len(s) != 3 {
					break
				}
				var urlFormat string
				switch proto {
				case "github":
					urlFormat = "https://raw.githubusercontent.com/%s/%s/%s"
				case "gitlab":
					urlFormat = "https://gitlab.com/%s/%s/-/raw/%s"
				default:
					panic("unreachable")
				}
				path = fmt.Sprintf(urlFormat, s[0], s[1], s[2])
				fn = FetchHTTP(path)
			case "http", "https":
				fn = FetchHTTP(src)
			}
			if fn != nil {
				fn = UseCache(!optCache, fn)
			}
		}
		if fn == nil {
			fn = FetchLocal(src)
		}
		inputs = append(inputs, fn)
	}

	if optCache {
		inputs = append(inputs, FetchLocal(optCacheDir))
	}
	fetch := FetchAny(inputs...)

	var lang []string
	for _, v := range languages {
		lang = append(lang, v)
	}
	quickfix := make(map[string]string)
	for _, in := range optQuickfix {
		s := strings.Split(in, "->")
		if len(s) != 2 {
			return fmt.Errorf("quickfix wants from->to but got=%v", in)
		}
		from, to := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
		from = strconv.Quote(from) + ":"
		to = strconv.Quote(to) + ":"
		quickfix[from] = to
	}
	return excel.LoadResources(excel.LoaderConfig{
		Languages: lang,
		ReadFile: func(root, name string) ([]byte, error) {
			var err error
			defer func() {
				if err == nil {
					Log(slog.LevelInfo, "loaded %s...", name)
				}
			}()
			data, err := fetch(root, name)
			for from, to := range quickfix {
				data = bytes.ReplaceAll(data, []byte(from), []byte(to))
			}
			return data, err
		},
	})
}

func (c *Compiled) build(config *Config) error {
	switch config.Kind {
	case KindNone:
		if config.Path != "pipeline/config.yml" {
			return errors.New("does not have kind set")
		}
	case KindArtifact:
		spec, err := buildArtifactSpec(config)
		if err != nil {
			return err
		}
		for lang, loc := range c.Localization {
			lang := languages[lang]
			loc.ArtifactNames[spec.Model.Key] = spec.ref.Affix(0).NameTextMapHash.Lang(lang)
		}
		config.Artifact = spec
	case KindCharacter:
		spec, err := buildCharacterSpec(config)
		if err != nil {
			return err
		}
		for lang, loc := range c.Localization {
			lang := languages[lang]
			name := spec.ref.NameTextMapHash.Lang(lang)
			if spec.ref.Name() == "Traveler" {
				switch spec.Model.Body {
				case model.BodyType_BODY_BOY:
					name = excel.FindManualTextMap("INFO_MALE_PRONOUN_KONG").Lang(lang)
				case model.BodyType_BODY_GIRL:
					name = excel.FindManualTextMap("INFO_MALE_PRONOUN_YING").Lang(lang)
				default:
					panic("unreachable")
				}
			}
			if len(spec.ref.CandSkillDepotIds) > 0 {
				name = fmt.Sprintf("%s (%s)", name, excel.FindManualTextMap(spec.Model.Element.String()).Lang(lang))
			}
			loc.CharacterNames[spec.Model.Key] = name
		}
		config.Character = spec
	case KindWeapon:
		spec, err := buildWeaponSpec(config)
		if err != nil {
			return err
		}
		for lang, loc := range c.Localization {
			lang := languages[lang]
			loc.WeaponNames[spec.Model.Key] = spec.ref.NameTextMapHash.Lang(lang)
		}
		config.Weapon = spec
	default:
		return fmt.Errorf("unknown config kind: %v", config.Kind)
	}

	slices.SortStableFunc(config.Attributes, func(a, b *AttributeSpec) int {
		aIdx, bIdx := slices.Index(abilities, a.Type), slices.Index(abilities, b.Type)
		if aIdx != -1 && bIdx != -1 {
			return cmp.Compare(aIdx, bIdx)
		}
		return 0
	})

	for i, name := range config.Shortcuts {
		config.Shortcuts[i] = excel.SlugLower(name)
	}
	slices.SortFunc(config.Shortcuts, cmp.Compare)
	config.Shortcuts = slices.Compact(config.Shortcuts)

	for _, name := range config.Combat.ICDGroups {
		att := config.Override.ICDGroup[name]
		if att.Name == "" {
			att.Name = name
		}
		name = excel.Slug(name)
		if _, ok := c.ICDGroup[name]; ok {
			return fmt.Errorf("icd group is used by another config: %v", name)
		}
		att.Name = excel.Slug(att.Name)
		c.ICDGroup[name] = att
	}

	for _, abil := range config.Abilities {
		if !slices.Contains(abilities, abil.Name) {
			return fmt.Errorf("unknown ability name: %v", abil.Name)
		}
		dupe := excel.Find(config.Abilities, func(v *Ability) bool { return abil != v && abil.Name == v.Name })
		if dupe != nil {
			return fmt.Errorf("found duplicate ability: %v", abil.Name)
		}
	}
	slices.SortFunc(config.Abilities, func(a, b *Ability) int {
		return cmp.Compare(
			slices.Index(abilities, a.Name),
			slices.Index(abilities, b.Name),
		)
	})

	return nil
}

func compile(ctx context.Context, cmd *cli.Command) (*Compiled, error) {
	if err := fetch(ctx, cmd); err != nil {
		return nil, err
	}

	configs, err := walkConfigs()
	if err != nil {
		return nil, err
	}

	c := &Compiled{
		Configuration: make([]*Config, 0, len(configs)),
		ElementCoeff:  make([]float64, 0),
		GrowCurveData: make(map[model.GrowCurveType][]float64),
		ICDGroup:      make(map[string]AttackAttenuation),
		Localization:  make(map[string]Localization),
	}
	for k := range languages {
		c.Localization[k] = Localization{
			ArtifactNames:  make(map[string]string),
			CharacterNames: make(map[string]string),
			MonsterNames:   make(map[string]string),
			WeaponNames:    make(map[string]string),
		}
	}
	for _, config := range configs {
		err := c.build(config)
		if err != nil {
			Log(slog.LevelError, "%v: %v", config.Path, err)
			continue
		}
		c.Configuration = append(c.Configuration, config)
	}

	for _, ref := range getMonsters() {
		desc := ref.Describe()
		spec, err := buildMonsterSpec(ref)
		if err != nil {
			Log(slog.LevelError, "id=%v,monster=%v: %v", ref.Id, excel.Slug(desc.Name()), err)
			continue
		}
		for lang, loc := range c.Localization {
			lang := languages[lang]
			loc.MonsterNames[spec.Model.Key] = desc.NameTextMapHash.Lang(lang)
		}
		c.Configuration = append(c.Configuration, &Config{
			Kind:    KindMonster,
			Name:    spec.Model.Key,
			Monster: spec,
		})
	}

	if err := c.buildCurveData(); err != nil {
		return nil, err
	}
	if err := c.buildElementCoeff(); err != nil {
		return nil, err
	}
	if err := c.buildICDGroup(); err != nil {
		return nil, err
	}

	slices.SortFunc(c.Configuration, func(a, b *Config) int { return cmp.Compare(a.SortKey(), b.SortKey()) })
	for _, cfg := range c.Configuration {
		dupe := excel.Find(c.Configuration, func(v *Config) bool { return cfg != v && cfg.SortKey() == v.SortKey() })
		if dupe != nil {
			return nil, fmt.Errorf("found duplicate config: %v %v", cfg.Path, dupe.Path)
		}
	}

	d, err := dumpYAML(c)
	if err != nil {
		return nil, err
	}
	writeFile("pipeline/_compiled.spec.yml", d)
	return c, nil
}

func check(ctx context.Context, cmd *cli.Command) error {
	c, err := compile(ctx, cmd)
	if err != nil {
		return err
	}

	for _, set := range excel.ReliquarySetExcelConfigData {
		if excel.Find(excel.ReliquaryCodexExcelConfigData, func(v *excel.ReliquaryCodex) bool { return v.SuitId == set.SetId }) == nil {
			continue
		}
		if excel.Find(c.Configuration, func(v *Config) bool {
			if spec := v.Artifact; spec != nil {
				return spec.ref.SetId == set.SetId
			}
			return false
		}) == nil {
			kind, name, build := KindArtifact, set.Affix(0).Name(), buildArtifactSpec
			Log(slog.LevelWarn, "%v %v is not registered", kind, excel.Slug(name))
			if _, err := build(&Config{Kind: kind, Name: excel.SlugLower(name)}); err != nil {
				Log(slog.LevelError, "%v %v: %v", kind, excel.Slug(name), err)
			}
		}
	}
	for _, codex := range excel.AvatarCodexExcelConfigData {
		if excel.Find(c.Configuration, func(v *Config) bool {
			if spec := v.Character; spec != nil {
				return spec.ref.Codex() == codex
			}
			return false
		}) == nil {
			kind, name, build := KindCharacter, codex.Avatar().Name(), buildCharacterSpec
			Log(slog.LevelWarn, "%v %v is not registered", kind, excel.Slug(name))
			if _, err := build(&Config{Kind: kind, Name: excel.SlugLower(name)}); err != nil {
				Log(slog.LevelError, "%v %v: %v", kind, excel.Slug(name), err)
			}
		}
	}
	for _, codex := range excel.WeaponCodexExcelConfigData {
		if codex.Weapon().StoryId == 0 {
			continue
		}
		if excel.Find(c.Configuration, func(v *Config) bool {
			if spec := v.Weapon; spec != nil {
				return spec.ref.Codex() == codex
			}
			return false
		}) == nil {
			kind, name, build := KindWeapon, codex.Weapon().Name(), buildWeaponSpec
			Log(slog.LevelWarn, "%v %v is not registered", kind, excel.Slug(name))
			if _, err := build(&Config{Kind: kind, Name: excel.SlugLower(name)}); err != nil {
				Log(slog.LevelError, "%v %v: %v", kind, excel.Slug(name), err)
			}
		}
	}

	return nil
}

func run(ctx context.Context, cmd *cli.Command) error {
	c, err := compile(ctx, cmd)
	if err != nil {
		return err
	}
	c.ClearRef()

	for _, fn := range []func() error{
		c.GenerateArtifacts,
		c.GenerateCharacters,
		c.GenerateWeapons,
		c.GenerateMonsters,
		c.GenerateCurve,
		c.GenerateElementCoeff,
		c.GenerateICDGroup,
		c.GenerateICDTag,
		c.GenerateEditorJS,
		c.GenerateLocalization,
	} {
		name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		name = strings.TrimSuffix(name, "-fm")
		name = strings.Split(name, ".")[2]
		t := time.Now()
		if err := fn(); err != nil {
			return fmt.Errorf("%v: %w", name, err)
		}
		Log(slog.LevelDebug, "%v took %v", name, time.Since(t))
	}

	return nil
}

func configTemplate(ctx context.Context, cmd *cli.Command) error {
	if err := fetch(ctx, cmd); err != nil {
		return err
	}

	cfg := &Config{
		Name: excel.SlugLower(cmd.String("name")),
		Override: Override{
			Id:    uint32(cmd.Int("id")),
			Depot: uint32(cmd.Int("depot")),
		},
	}
	var err error
	if cfg.Artifact, err = buildArtifactSpec(cfg); err != nil {
		Log(slog.LevelError, "%v: %v", KindArtifact, err)
	}
	if cfg.Character, err = buildCharacterSpec(cfg); err != nil {
		Log(slog.LevelError, "%v: %v", KindCharacter, err)
	}
	if cfg.Weapon, err = buildWeaponSpec(cfg); err != nil {
		Log(slog.LevelError, "%v: %v", KindWeapon, err)
	}
	switch {
	case cfg.Artifact != nil:
		cfg.Kind = KindArtifact
		cfg.Name = cfg.Artifact.Model.Key
	case cfg.Character != nil:
		cfg.Kind = KindCharacter
		cfg.Name = cfg.Character.Model.Key
	case cfg.Weapon != nil:
		cfg.Kind = KindWeapon
		cfg.Name = cfg.Weapon.Model.Key
	default:
		Log(slog.LevelError, "no result")
		cli.ShowSubcommandHelp(cmd)
		return nil
	}

	indent := strings.Repeat(" ", 2)
	b := bytes.NewBuffer(nil)
	b.WriteString("use: pipeline\n")
	fmt.Fprintf(b, "kind: %s\n", cfg.Kind)
	fmt.Fprintf(b, "name: %s\n", cfg.Name)
	b.WriteString("\n")
	b.WriteString("talents:\n")
	var seen string
	for _, attr := range cfg.Attributes {
		useParam := attr.ParamDesc == ""
		useParam = useParam || len(excel.Filter(cfg.Attributes, func(s *AttributeSpec) bool {
			return attr.Type == s.Type && attr.Desc == s.Desc
		})) > 1
		for ind := range attr.Index {
			if ind == 0 && seen != attr.Type {
				if seen != "" {
					b.WriteString("\n")
				}
				seen = attr.Type
				if attr.ParamDesc == "" {
					desc := attr.Desc
					desc = wrapText(desc, 80)
					desc = indentText(indent+"# ", desc)
					b.WriteString(desc)
				}
			}

			b.WriteString(indent)
			fmt.Fprintf(b, "- _%[1]s: %[1]s", attr.Type)
			if useParam {
				fmt.Fprintf(b, "#%d", attr.Index[ind])
				switch {
				case attr.ParamDesc != "":
					fmt.Fprintf(b, " # %s | %s", attr.Desc, attr.ParamDesc)
				case attr.Values[ind] != nil:
					fmt.Fprintf(b, " # %v", attr.Values[ind])
				default:
					fmt.Fprintf(b, " # %v", attr.Const[ind])
				}
			} else {
				b.WriteString("(")
				b.WriteString(attr.Desc)
				if ind > 0 {
					fmt.Fprintf(b, "|%d", ind)
				}
				b.WriteString(")")
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
	fmt.Println(b.String())

	return nil
}
