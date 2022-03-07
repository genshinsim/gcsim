package testhelper

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func TestSkillCooldownSingleCharge(c *core.Core, char core.Character, cd int) error {
	setupChar(c, char)
	p := make(map[string]int)
	char.Skill(p)
	SkipFrames(c, cd-1)
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("skill shouldn't be ready yet")
	}
	x := char.Cooldown(core.ActionSkill)
	if x != 1 {
		return fmt.Errorf("expecting cooldown to be 1, got %v", x)
	}
	SkipFrames(c, 1)
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("skill should be ready now")
	}
	x = char.Cooldown(core.ActionSkill)
	if x != 0 {
		return fmt.Errorf("expecting cooldown to be 0, got %v", x)
	}
	return nil
}

func TestSkillCooldownDoubleCharge(c *core.Core, char core.Character, cd []int) error {
	setupChar(c, char)
	p := make(map[string]int)
	char.Skill(p)
	//should have one charge left
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("used skill once, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("used skill once, skill should be ready still after using one charge")
	}
	if char.Cooldown(core.ActionSkill) > 0 {
		return fmt.Errorf("used skill once, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//use skill again without skipping
	char.Skill(p)
	if char.Charges(core.ActionSkill) != 0 {
		return fmt.Errorf("used skill twice, expecting to have 0 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("used skill twice, skill shouldn't be ready")
	}
	if char.Cooldown(core.ActionSkill) == 0 {
		return fmt.Errorf("used skill twice, expecting cooldown to be > 0, got %v", char.Cooldown(core.ActionSkill))
	}
	SkipFrames(c, cd[0]-1)
	if char.Charges(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 1st charge, expecting to have 0 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 1st charge, skill shouldn't be ready")
	}
	if char.Cooldown(core.ActionSkill) == 0 {
		return fmt.Errorf("checking 1st charge, expecting cooldown to be > 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//1 more frame for first charge
	SkipFrames(c, 1)
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("checking 1st charge, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 1st charge, skill should be ready still after using one charge")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 1st charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//n[1] more frames for 2nd charge
	SkipFrames(c, cd[1]-1)
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("checking 2nd charge, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 2nd charge, skill should be ready")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 2nd charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//1 more frame for second charge
	SkipFrames(c, 1)
	if char.Charges(core.ActionSkill) != 2 {
		return fmt.Errorf("checking 2nd charge, expecting to have 2 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 2nd charge, skill should be ready")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 2nd charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	return nil
}

func TestFlatCDReductionSingleCharge(c *core.Core, char core.Character, cd int) error {
	setupChar(c, char)
	p := make(map[string]int)
	char.Skill(p)
	SkipFrames(c, cd-10)
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("skill shouldn't be ready yet")
	}
	if char.Cooldown(core.ActionSkill) != 10 {
		return fmt.Errorf("expecting cooldown to be %v, got %v", 10, char.Cooldown(core.ActionSkill))
	}
	//reduce by 5
	char.ReduceActionCooldown(core.ActionSkill, 5)
	SkipFrames(c, 4)
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("skill shouldn't be ready yet")
	}
	if char.Cooldown(core.ActionSkill) != 1 {
		return fmt.Errorf("expecting cooldown to be %v, got %v", 1, char.Cooldown(core.ActionSkill))
	}
	SkipFrames(c, 1)
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("skill should be ready now")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("expecting cooldown to be %v, got %v", 0, char.Cooldown(core.ActionSkill))
	}
	return nil
}

func TestFlatCDReductionDoubleCharge(c *core.Core, char core.Character, cd []int) error {
	setupChar(c, char)
	p := make(map[string]int)
	char.Skill(p)
	//should have one charge left
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("used skill once, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("used skill once, skill should be ready still after using one charge")
	}
	if char.Cooldown(core.ActionSkill) > 0 {
		return fmt.Errorf("used skill once, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//use skill again without skipping
	char.Skill(p)
	if char.Charges(core.ActionSkill) != 0 {
		return fmt.Errorf("used skill twice, expecting to have 0 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("used skill twice, skill shouldn't be ready")
	}
	if char.Cooldown(core.ActionSkill) == 0 {
		return fmt.Errorf("used skill twice, expecting cooldown to be > 0, got %v", char.Cooldown(core.ActionSkill))
	}
	SkipFrames(c, cd[0]-10)
	if char.Charges(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 1st charge, expecting to have 0 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 1st charge, skill shouldn't be ready")
	}
	if char.Cooldown(core.ActionSkill) != 10 {
		return fmt.Errorf("checking 1st charge, expecting cooldown to be 10, got %v", char.Cooldown(core.ActionSkill))
	}
	//skip 5
	char.ReduceActionCooldown(core.ActionSkill, 5)
	SkipFrames(c, 4)
	if char.Charges(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 1st charge, expecting to have 0 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 1st charge, skill shouldn't be ready")
	}
	if char.Cooldown(core.ActionSkill) != 1 {
		return fmt.Errorf("checking 1st charge, expecting cooldown to be 1, got %v", char.Cooldown(core.ActionSkill))
	}
	//1 more frame for first charge
	SkipFrames(c, 1)
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("checking 1st charge, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 1st charge, skill should be ready still after using one charge")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 1st charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//skip 10 frames
	SkipFrames(c, cd[1]-1)
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("checking 2nd charge, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 2nd charge, skill should be ready")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 2nd charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//1 more frame for second charge
	SkipFrames(c, 1)
	if char.Charges(core.ActionSkill) != 2 {
		return fmt.Errorf("checking 2nd charge, expecting to have 2 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 2nd charge, skill should be ready")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 2nd charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	return nil
}

func TestResetSkillCDSingleCharge(c *core.Core, char core.Character, cd int) error {
	setupChar(c, char)
	p := make(map[string]int)
	char.Skill(p)
	SkipFrames(c, cd-10)
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("skill shouldn't be ready yet")
	}
	if char.Cooldown(core.ActionSkill) != 10 {
		return fmt.Errorf("expecting cooldown to be %v, got %v", 10, char.Cooldown(core.ActionSkill))
	}
	//reset
	char.ResetActionCooldown(core.ActionSkill)
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("skill should be ready now")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("expecting cooldown to be %v, got %v", 0, char.Cooldown(core.ActionSkill))
	}
	return nil
}

func TestResetSkillCDDoubleCharge(c *core.Core, char core.Character, cd []int) error {
	setupChar(c, char)
	p := make(map[string]int)
	char.Skill(p)
	//should have one charge left
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("used skill once, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("used skill once, skill should be ready still after using one charge")
	}
	if char.Cooldown(core.ActionSkill) > 0 {
		return fmt.Errorf("used skill once, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//use skill again without skipping
	char.Skill(p)
	if char.Charges(core.ActionSkill) != 0 {
		return fmt.Errorf("used skill twice, expecting to have 0 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("used skill twice, skill shouldn't be ready")
	}
	if char.Cooldown(core.ActionSkill) == 0 {
		return fmt.Errorf("used skill twice, expecting cooldown to be > 0, got %v", char.Cooldown(core.ActionSkill))
	}
	SkipFrames(c, cd[0]-10)
	if char.Charges(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 1st charge, expecting to have 0 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 1st charge, skill shouldn't be ready")
	}
	if char.Cooldown(core.ActionSkill) != 10 {
		return fmt.Errorf("checking 1st charge, expecting cooldown to be 10, got %v", char.Cooldown(core.ActionSkill))
	}
	//reset
	char.ResetActionCooldown(core.ActionSkill)
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("checking 1st charge, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 1st charge, skill should be ready still after using one charge")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 1st charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	SkipFrames(c, cd[1]-1)
	if char.Charges(core.ActionSkill) != 1 {
		return fmt.Errorf("checking 2nd charge, expecting to have 1 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 2nd charge, skill should be ready")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 2nd charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	//1 more frame for second charge
	SkipFrames(c, 1)
	if char.Charges(core.ActionSkill) != 2 {
		return fmt.Errorf("checking 2nd charge, expecting to have 2 charge left, got %v", char.Charges(core.ActionSkill))
	}
	if !char.ActionReady(core.ActionSkill, p) {
		return errors.New("checking 2nd charge, skill should be ready")
	}
	if char.Cooldown(core.ActionSkill) != 0 {
		return fmt.Errorf("checking 2nd charge, expecting cooldown to be 0, got %v", char.Cooldown(core.ActionSkill))
	}
	return nil
}

func TestSkillCDSingleCharge(c *core.Core, char core.Character, cd int) error {
	var err error
	err = TestSkillCooldownSingleCharge(c, char, cd)
	if err != nil {
		return err
	}
	err = TestFlatCDReductionSingleCharge(c, char, cd)
	if err != nil {
		return err
	}
	err = TestResetSkillCDSingleCharge(c, char, cd)
	if err != nil {
		return err
	}
	err = TestResetSkillCDSingleCharge(c, char, cd)
	if err != nil {
		return err
	}
	return nil
}

func TestSkillCDDoubleCharge(c *core.Core, char core.Character, cd []int) error {
	var err error
	err = TestSkillCooldownDoubleCharge(c, char, cd)
	if err != nil {
		return err
	}
	err = TestFlatCDReductionDoubleCharge(c, char, cd)
	if err != nil {
		return err
	}
	err = TestResetSkillCDDoubleCharge(c, char, cd)
	if err != nil {
		return err
	}
	return nil
}
