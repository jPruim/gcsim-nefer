package nefer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

const (
	skillHitmark        = 24
	bloomRemovalKey     = "nefer-e"
	particleICDKey      = "nefer-skill-particle-icd"
	phantasmParticleKey = "nefer-phantasm-particle-icd"
	phantasmCountKey    = "nefer-phantasm-count"
)

func init() {
	skillFrames = frames.InitAbilSlice(52) // E -> Q
	skillFrames[action.ActionAttack] = 26  // E -> N1
	skillFrames[action.ActionDash] = 38    // E -> D
	skillFrames[action.ActionJump] = 38    // E -> J
	skillFrames[action.ActionSwap] = 25    // E -> Swap
	skillFrames[action.ActionCharge] = 29  // E -> CA, doesn't include CA windup (20f)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Senet Strategy",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   30, // TODO: Get Value
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       phantasm_nefer_1_att[c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * phantasm_nefer_1_em[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), info.Point{Y: 2}, 2.25), skillHoldHitmark, skillHoldHitmark)

	c.SetCDWithDelay(action.ActionSkill, 9*60, 23)

	// Add phantasm count
	c.skillCAOverride()

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) initialParticleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.1*60, true)

	count := 2.0
	if c.Core.Rand.Float64() < 0.67 {
		count = 3.0
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Dendro, c.ParticleDelay)
}

// Set CA override count
func (c *char) skillCAOverride() {
	c.AddStatus(phantasmCountKey, 9*60, true) // 9 sec duration
	c.SetTag(phantasmCountKey, 3)             // 3 phantasms per skill
}
