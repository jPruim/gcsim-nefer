package nefer

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Nefer, NewChar)
}

const (
	isneferkey = "nefer-identifier"
)

type char struct {
	*tmpl.Character
	veilstacks     int
	phantasmStacks int
	seeds          int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	// Set number of skil charges
	c.SetNumCharges(action.ActionSkill, 2)
	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// TODO: cleaner methodology?
	c.AddStatus(isneferkey, -1, true)

	// Start with 0 veil stacks
	c.veilstacks = 0

	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 0
	}
	if k == info.AnimationYelanN0StartDelay {
		return 0
	}
	return c.Character.AnimationStartDelay(k)
}
