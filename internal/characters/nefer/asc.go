package nefer

import (
	"github.com/genshinsim/gcsim/internal/template/dendrocore"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1Status = "nefer-a1"
	a4Mod    = "nefer-a4"
)

// Moonsign: Ascendant Gleam: When she unleashes her Elemental Skill Senet Strategy: Dance of a Thousand Nights,
// any Dendro Cores on the field will be converted to Seeds of Deceit, and any Lunar-Bloom reactions triggered by
// nearby characters in the following 15s that would create Dendro Cores or Bountiful Cores will instead create
// Seeds of Deceit. Seeds of Deceit cannot trigger Hyperbloom or Burgeon reactions and will not burst.
// When Nefer unleashes a Charged Attack or Phantasm Performance, she can absorb Seeds of Deceit within a
// certain range, gaining 1 stack of Veil of Falsehood for every seed absorbed.
// When this effect reaches 3 stacks, or when the third stack's duration is refreshed,
// Nefer's Elemental Mastery will be increased by 100 for 8s.
func (c *char) a1() {
	// Seeds of Deceit
	c.Core.Events.Subscribe(event.OnDendroCore, func(args ...any) bool {
		g, ok := args[0].(*dendrocore.Gadget)
		if !ok {
			return false
		}
		b := newSeedOfDeceit(c.Core, g.Pos())
		b.SetKey(g.Key())
		c.Core.Combat.ReplaceGadget(g.Key(), b)
		// prevent blowing up
		g.OnExpiry = nil
		g.OnKill = nil

		return false
	}, "nefer-a1-cores")

	c.Core.Events.Subscribe(event.OnChargeAttack, func(args ...any) bool {
		charIndex := args[0].(int)
		char := c.Core.Player.ByIndex(charIndex)
		if !char.StatusIsActive(isneferkey) {
			return false
		}

		if c.veilstacks < 3 {
			return false
		}
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 100
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("nefer-a1-em", 8*60),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}, "nilou-a1")
}

// Every 1,000 points of Nilou’s Max HP above 30,000 will cause the DMG dealt by Bountiful Cores created by characters affected
// by Golden Chalice’s Bounty to increase by 9%.
// The maximum increase in Bountiful Core DMG that can be achieved this way is 400%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	for _, this := range c.Core.Player.Chars() {
		// TODO: a4 should be an extra buff
		this.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag(a4Mod, 30*60),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				if ai.AttackTag != attacks.AttackTagBloom {
					return 0, false
				}
				if ai.ICDTag != attacks.ICDTagBountifulCoreDamage {
					return 0, false
				}

				c.Core.Combat.Log.NewEvent("adding nilou a4 bonus", glog.LogCharacterEvent, c.Index()).Write("bonus", c.a4Bonus)
				return c.a4Bonus, false
			},
		})
	}
	c.a4Src = c.Core.F
	c.QueueCharTask(c.updateA4Bonus(c.a4Src), 0.5*60)
}

func (c *char) updateA4Bonus(src int) func() {
	return func() {
		if c.a4Src != src {
			return
		}
		if !c.ReactBonusModIsActive(a4Mod) {
			return
		}

		c.a4Bonus = (c.MaxHP() - 30000) * 0.001 * 0.09
		if c.a4Bonus < 0 {
			c.a4Bonus = 0
		} else if c.a4Bonus > 4 {
			c.a4Bonus = 4
		}

		c.QueueCharTask(c.updateA4Bonus(src), 0.5*60)
	}
}
