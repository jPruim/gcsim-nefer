package nefer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	chargeHitmark = 65
	windup = 20
)
var (
	chargeFrames            []int
	chargeHitmarks          = []int{14, 25}
	chargeHitlagHaltFrame   = []float64{0, 0.10}
	chargeDefHalt           = []bool{false, true} // TODO Check
	phantasmframes          []int
	phantasmHitmarks        = []int{14, 25}
	phantasmHitlagHaltFrame = []float64{0, 0.10}
	phantasmDefHalt         = []bool{false, true} // TODO Check
)

func init() {
	chargeFrames = frames.InitAbilSlice(35) // CA -> N1/E/Q
	chargeFrames[action.ActionDash] = 31    // CA -> D
	chargeFrames[action.ActionJump] = 31    // CA -> J
	chargeFrames[action.ActionSwap] = 29    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if(c.con)
}

func (c *char) basicChargeAttack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	windup := 0
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		windup = 10
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			info.Point{Y: -3},
			6,
			6,
		),
		chargeHitmark-windup,
		chargeHitmark-windup,
	)
	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionDash] - windup,
		State:           action.ChargeAttackState,
	}, nil
}
