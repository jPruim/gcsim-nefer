package nefer

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type SeedOfDeceit struct {
	srcFrame int
	*gadget.Gadget
}

func newSeedOfDeceit(c *core.Core, p info.Point) *SeedOfDeceit {
	b := &SeedOfDeceit{
		srcFrame: c.F,
	}

	b.Gadget = gadget.New(c, p, 2, info.GadgetTypDendroCore)
	// TODO Check duration, assuming same as bloom block
	b.Duration = 15 * 60

	absorb := func() {
		for _, this := range c.Core.Player.Chars() {
			// TODO: Cleaner methodology?
			if char.StatusIsActive(isneferkey) {
				this.veilstacks++
			}
		}
	}
	b.OnExpiry = nil
	b.OnKill = absorb

	return b
}

func (b *SeedOfDeceit) Tick() {
	// this is needed since gadget tick
	b.Gadget.Tick()
}

func (b *SeedOfDeceit) HandleAttack(atk *info.AttackEvent) float64 {
	b.Core.Events.Emit(event.OnGadgetHit, b, atk)
	return 0
}
func (b *SeedOfDeceit) Attack(*info.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (b *SeedOfDeceit) SetDirection(trg info.Point)                          {}
func (b *SeedOfDeceit) SetDirectionToClosestEnemy()                          {}
func (b *SeedOfDeceit) CalcTempDirection(trg info.Point) info.Point {
	return info.DefaultDirection()
}
