package parse

import (
	"errors"
	"fmt"
)

func parseOptions(p *Parser) (parseFn, error) {
	//option iter=1000 duration=1000 worker=50 debug=true er_calc=true damage_mode=true
	var err error

	//options debug=true iteration=5000 duration=90 workers=24;
	for n := p.next(); n.typ != itemEOF; n = p.next() {

		switch n.typ {
		case itemIdentifier:
			//expecting identifier = some value
			switch n.val {
			case "debug":
				n, err = p.acceptSeqReturnLast(itemEqual, itemBool)
				// every run is going to have a debug from now on so we basically ignore what this flag says
			case "iteration":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Iterations, err = itemNumberToInt(n)
				}
			case "duration":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Duration, err = itemNumberToInt(n)
				}
			case "workers":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.NumberOfWorkers, err = itemNumberToInt(n)
				}
			case "mode":
				n, err = p.acceptSeqReturnLast(itemEqual, itemIdentifier)
				if err == nil {
					//should be either apl or sl
					m, ok := queueModeKeys[n.val]
					if !ok {
						return nil, fmt.Errorf("ln%v: invalid queue mode, got %v", n.line, n.val)
					}
					p.cfg.Settings.QueueMode = m
				}
			case "swap_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Delays.Swap, err = itemNumberToInt(n)
				}
			case "attack_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Delays.Attack, err = itemNumberToInt(n)
				}
			case "charge_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Delays.Charge, err = itemNumberToInt(n)
				}
			case "skill_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Delays.Skill, err = itemNumberToInt(n)
				}
			case "burst_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Delays.Burst, err = itemNumberToInt(n)
				}
			case "jump_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Delays.Jump, err = itemNumberToInt(n)
				}
			case "dash_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Delays.Dash, err = itemNumberToInt(n)
				}
			case "aim_delay":
				n, err = p.acceptSeqReturnLast(itemEqual, itemNumber)
				if err == nil {
					p.cfg.Settings.Delays.Aim, err = itemNumberToInt(n)
				}
			case "frame_defaults":
				n, err = p.acceptSeqReturnLast(itemEqual, itemIdentifier)
				if err == nil {
					switch n.val {
					case "human":
						p.cfg.Settings.Delays.Swap = 8
						p.cfg.Settings.Delays.Attack = 5
						p.cfg.Settings.Delays.Charge = 5
						p.cfg.Settings.Delays.Skill = 5
						p.cfg.Settings.Delays.Burst = 5
						p.cfg.Settings.Delays.Dash = 5
						p.cfg.Settings.Delays.Jump = 5
						p.cfg.Settings.Delays.Aim = 5
					default:
						return nil, fmt.Errorf("ln%v: unrecognized option for frame_defaults specified: %v", n.line, n.val)
					}
				}
			case "er_calc":
				//does nothing thus far...
			default:
				return nil, fmt.Errorf("ln%v: unrecognized option specified: %v", n.line, n.val)
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
