package nefer

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames          [][]int
	attackHits            = []int{1, 1, 2, 1}
	attackHitmarks        = []int{12, 9, 10 + 16, 22}
	attackHitlagHaltFrame = []float64{0.03, 0.03, 0.06, 0.10}          // TODO get values
	attackHitboxes        = [][]float64{{1.5}, {1.5}, {6, 3.3}, {2.2}} // TODO get values
	attackOffsets         = []float64{1, 0.6, -0.2, 1.2}               // TODO get values
)

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 40) // N1 -> CA
	attackFrames[0][action.ActionAttack] = 15                             // N1 -> N2

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 38) // N2 -> CA
	attackFrames[1][action.ActionAttack] = 20                             // N2 -> N3

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 62) // N3 -> CA
	attackFrames[2][action.ActionAttack] = 49                             // N3 -> N4

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 51) // N4 -> N1
	attackFrames[3][action.ActionCharge] = 500                            // N4 -> CA, TODO: this action is illegal; need better way to handle it
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	for i := 0; i < attackHits[c.NormalCounter]; i++ {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNormalAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeSlash,
			Element:    attributes.Dendro,
			Durability: 25,
			Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: attackOffsets[c.NormalCounter]},
			attackHitboxes[c.NormalCounter][0],
		)
		if c.NormalCounter == 2 {
			ap = combat.NewBoxHitOnTarget(
				c.Core.Combat.Player(),
				info.Point{Y: attackOffsets[c.NormalCounter]},
				attackHitboxes[c.NormalCounter][0],
				attackHitboxes[c.NormalCounter][1],
			)
		}
		// no multihits so no need for char queue here
		c.Core.QueueAttack(ai, ap, attackHitmarks[c.NormalCounter], attackHitmarks[c.NormalCounter])
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}
