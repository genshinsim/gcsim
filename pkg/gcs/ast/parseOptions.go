package ast

import (
	"errors"
	"fmt"
)

func parseOptions(p *Parser) (parseFn, error) {
	// option iter=1000 duration=1000 worker=50 debug=true er_calc=true damage_mode=true
	var err error

	// options debug=true iteration=5000 duration=90 workers=24;
	for n := p.next(); n.Typ != itemEOF; n = p.next() {
		switch n.Typ {
		case itemIdentifier:
			// expecting identifier = some value
			switch n.Val {
			case "debug":
				_, err = p.acceptSeqReturnLast(itemAssign, itemBool)
				// every run is going to have a debug from now on so we basically ignore what this flag says
			case "defhalt":
				n, err = p.acceptSeqReturnLast(itemAssign, itemBool)
				p.res.Settings.DefHalt = n.Val == TrueVal
			case "hitlag":
				n, err = p.acceptSeqReturnLast(itemAssign, itemBool)
				p.res.Settings.EnableHitlag = n.Val == TrueVal
			case "iteration":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Iterations, err = itemNumberToInt(n)
				}
			case "duration":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Duration, err = itemNumberToFloat64(n)
				}
			case "workers":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.NumberOfWorkers, err = itemNumberToInt(n)
				}
			case "mode":
				//TODO: this is for backward compatibility for now
				_, err = p.acceptSeqReturnLast(itemAssign, itemIdentifier)
			case "swap_delay":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Delays.Swap, err = itemNumberToInt(n)
				}
			case "attack_delay":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Delays.Attack, err = itemNumberToInt(n)
				}
			case "charge_delay":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Delays.Charge, err = itemNumberToInt(n)
				}
			case "skill_delay":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Delays.Skill, err = itemNumberToInt(n)
				}
			case "burst_delay":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Delays.Burst, err = itemNumberToInt(n)
				}
			case "jump_delay":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Delays.Jump, err = itemNumberToInt(n)
				}
			case "dash_delay":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Delays.Dash, err = itemNumberToInt(n)
				}
			case "aim_delay":
				n, err = p.acceptSeqReturnLast(itemAssign, itemNumber)
				if err == nil {
					p.res.Settings.Delays.Aim, err = itemNumberToInt(n)
				}
			case "frame_defaults":
				n, err = p.acceptSeqReturnLast(itemAssign, itemIdentifier)
				if err == nil {
					switch n.Val {
					case "human":
						p.res.Settings.Delays.Swap = 8
						p.res.Settings.Delays.Attack = 5
						p.res.Settings.Delays.Charge = 5
						p.res.Settings.Delays.Skill = 5
						p.res.Settings.Delays.Burst = 5
						p.res.Settings.Delays.Dash = 5
						p.res.Settings.Delays.Jump = 5
						p.res.Settings.Delays.Aim = 5
					default:
						return nil, fmt.Errorf("ln%v: unrecognized option for frame_defaults specified: %v", n.line, n.Val)
					}
				}
			case "er_calc":
				// Current there is no method to easily show this data to the user so it's removed for now
				// n, err = p.acceptSeqReturnLast(itemAssign, itemBool)
				// p.res.Settings.ErCalc = n.Val == TrueVal
			default:
				return nil, fmt.Errorf("ln%v: unrecognized option specified: %v", n.line, n.Val)
			}
		case itemTerminateLine:
			return parseRows, nil
		default:
			return nil, fmt.Errorf("ln%v: unrecognized token parsing options: %v", n.line, n)
		}
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("unexpected end of line while parsing options")
}
