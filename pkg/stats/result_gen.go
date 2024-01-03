package stats

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ActionEvent) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "action_id":
			z.ActionId, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "action":
			z.Action, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ActionEvent) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "frame"
	err = en.Append(0x83, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Frame)
	if err != nil {
		return
	}
	// write "action_id"
	err = en.Append(0xa9, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.ActionId)
	if err != nil {
		return
	}
	// write "action"
	err = en.Append(0xa6, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Action)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ActionEvent) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "frame"
	o = append(o, 0x83, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	o = msgp.AppendInt(o, z.Frame)
	// string "action_id"
	o = append(o, 0xa9, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64)
	o = msgp.AppendInt(o, z.ActionId)
	// string "action"
	o = append(o, 0xa6, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Action)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ActionEvent) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "action_id":
			z.ActionId, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "action":
			z.Action, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ActionEvent) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 10 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.Action)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ActionFailInterval) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zbai uint32
	zbai, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zbai > 0 {
		zbai--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "end":
			z.End, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "reason":
			z.Reason, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ActionFailInterval) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "start"
	err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Start)
	if err != nil {
		return
	}
	// write "end"
	err = en.Append(0xa3, 0x65, 0x6e, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.End)
	if err != nil {
		return
	}
	// write "reason"
	err = en.Append(0xa6, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Reason)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ActionFailInterval) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "start"
	o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	o = msgp.AppendInt(o, z.Start)
	// string "end"
	o = append(o, 0xa3, 0x65, 0x6e, 0x64)
	o = msgp.AppendInt(o, z.End)
	// string "reason"
	o = append(o, 0xa6, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Reason)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ActionFailInterval) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcmr uint32
	zcmr, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcmr > 0 {
		zcmr--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "end":
			z.End, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "reason":
			z.Reason, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ActionFailInterval) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 4 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.Reason)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ActiveCharacterInterval) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zajw uint32
	zajw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zajw > 0 {
		zajw--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "end":
			z.End, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "character":
			z.Character, err = dc.ReadInt()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ActiveCharacterInterval) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "start"
	err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Start)
	if err != nil {
		return
	}
	// write "end"
	err = en.Append(0xa3, 0x65, 0x6e, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.End)
	if err != nil {
		return
	}
	// write "character"
	err = en.Append(0xa9, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Character)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ActiveCharacterInterval) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "start"
	o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	o = msgp.AppendInt(o, z.Start)
	// string "end"
	o = append(o, 0xa3, 0x65, 0x6e, 0x64)
	o = msgp.AppendInt(o, z.End)
	// string "character"
	o = append(o, 0xa9, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72)
	o = msgp.AppendInt(o, z.Character)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ActiveCharacterInterval) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zwht uint32
	zwht, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zwht > 0 {
		zwht--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "end":
			z.End, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "character":
			z.Character, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ActiveCharacterInterval) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 4 + msgp.IntSize + 10 + msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *CharacterResult) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zrsw uint32
	zrsw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zrsw > 0 {
		zrsw--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "damage_events":
			var zxpk uint32
			zxpk, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.DamageEvents) >= int(zxpk) {
				z.DamageEvents = (z.DamageEvents)[:zxpk]
			} else {
				z.DamageEvents = make([]DamageEvent, zxpk)
			}
			for zhct := range z.DamageEvents {
				err = z.DamageEvents[zhct].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "reaction_events":
			var zdnj uint32
			zdnj, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.ReactionEvents) >= int(zdnj) {
				z.ReactionEvents = (z.ReactionEvents)[:zdnj]
			} else {
				z.ReactionEvents = make([]ReactionEvent, zdnj)
			}
			for zcua := range z.ReactionEvents {
				err = z.ReactionEvents[zcua].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "action_events":
			var zobc uint32
			zobc, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.ActionEvents) >= int(zobc) {
				z.ActionEvents = (z.ActionEvents)[:zobc]
			} else {
				z.ActionEvents = make([]ActionEvent, zobc)
			}
			for zxhx := range z.ActionEvents {
				var zsnv uint32
				zsnv, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zsnv > 0 {
					zsnv--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "frame":
						z.ActionEvents[zxhx].Frame, err = dc.ReadInt()
						if err != nil {
							return
						}
					case "action_id":
						z.ActionEvents[zxhx].ActionId, err = dc.ReadInt()
						if err != nil {
							return
						}
					case "action":
						z.ActionEvents[zxhx].Action, err = dc.ReadString()
						if err != nil {
							return
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
			}
		case "energy_events":
			var zkgt uint32
			zkgt, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.EnergyEvents) >= int(zkgt) {
				z.EnergyEvents = (z.EnergyEvents)[:zkgt]
			} else {
				z.EnergyEvents = make([]EnergyEvent, zkgt)
			}
			for zlqf := range z.EnergyEvents {
				err = z.EnergyEvents[zlqf].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "heal_events":
			var zema uint32
			zema, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.HealEvents) >= int(zema) {
				z.HealEvents = (z.HealEvents)[:zema]
			} else {
				z.HealEvents = make([]HealEvent, zema)
			}
			for zdaf := range z.HealEvents {
				err = z.HealEvents[zdaf].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "failed_actions":
			var zpez uint32
			zpez, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.FailedActions) >= int(zpez) {
				z.FailedActions = (z.FailedActions)[:zpez]
			} else {
				z.FailedActions = make([]ActionFailInterval, zpez)
			}
			for zpks := range z.FailedActions {
				var zqke uint32
				zqke, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zqke > 0 {
					zqke--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "start":
						z.FailedActions[zpks].Start, err = dc.ReadInt()
						if err != nil {
							return
						}
					case "end":
						z.FailedActions[zpks].End, err = dc.ReadInt()
						if err != nil {
							return
						}
					case "reason":
						z.FailedActions[zpks].Reason, err = dc.ReadString()
						if err != nil {
							return
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
			}
		case "energy_status":
			var zqyh uint32
			zqyh, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.EnergyStatus) >= int(zqyh) {
				z.EnergyStatus = (z.EnergyStatus)[:zqyh]
			} else {
				z.EnergyStatus = make([]float64, zqyh)
			}
			for zjfb := range z.EnergyStatus {
				z.EnergyStatus[zjfb], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "health_status":
			var zyzr uint32
			zyzr, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.HealthStatus) >= int(zyzr) {
				z.HealthStatus = (z.HealthStatus)[:zyzr]
			} else {
				z.HealthStatus = make([]float64, zyzr)
			}
			for zcxo := range z.HealthStatus {
				z.HealthStatus[zcxo], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "damage_cumulative_contrib":
			var zywj uint32
			zywj, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.DamageCumulativeContrib) >= int(zywj) {
				z.DamageCumulativeContrib = (z.DamageCumulativeContrib)[:zywj]
			} else {
				z.DamageCumulativeContrib = make([]float64, zywj)
			}
			for zeff := range z.DamageCumulativeContrib {
				z.DamageCumulativeContrib[zeff], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "active_time":
			z.ActiveTime, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "energy_spent":
			z.EnergySpent, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "energy_info":
			err = z.EnergyInfo.DecodeMsg(dc)
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *CharacterResult) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 13
	// write "name"
	err = en.Append(0x8d, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "damage_events"
	err = en.Append(0xad, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.DamageEvents)))
	if err != nil {
		return
	}
	for zhct := range z.DamageEvents {
		err = z.DamageEvents[zhct].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "reaction_events"
	err = en.Append(0xaf, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.ReactionEvents)))
	if err != nil {
		return
	}
	for zcua := range z.ReactionEvents {
		err = z.ReactionEvents[zcua].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "action_events"
	err = en.Append(0xad, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.ActionEvents)))
	if err != nil {
		return
	}
	for zxhx := range z.ActionEvents {
		// map header, size 3
		// write "frame"
		err = en.Append(0x83, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.ActionEvents[zxhx].Frame)
		if err != nil {
			return
		}
		// write "action_id"
		err = en.Append(0xa9, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.ActionEvents[zxhx].ActionId)
		if err != nil {
			return
		}
		// write "action"
		err = en.Append(0xa6, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e)
		if err != nil {
			return err
		}
		err = en.WriteString(z.ActionEvents[zxhx].Action)
		if err != nil {
			return
		}
	}
	// write "energy_events"
	err = en.Append(0xad, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.EnergyEvents)))
	if err != nil {
		return
	}
	for zlqf := range z.EnergyEvents {
		err = z.EnergyEvents[zlqf].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "heal_events"
	err = en.Append(0xab, 0x68, 0x65, 0x61, 0x6c, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.HealEvents)))
	if err != nil {
		return
	}
	for zdaf := range z.HealEvents {
		err = z.HealEvents[zdaf].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "failed_actions"
	err = en.Append(0xae, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.FailedActions)))
	if err != nil {
		return
	}
	for zpks := range z.FailedActions {
		// map header, size 3
		// write "start"
		err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.FailedActions[zpks].Start)
		if err != nil {
			return
		}
		// write "end"
		err = en.Append(0xa3, 0x65, 0x6e, 0x64)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.FailedActions[zpks].End)
		if err != nil {
			return
		}
		// write "reason"
		err = en.Append(0xa6, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e)
		if err != nil {
			return err
		}
		err = en.WriteString(z.FailedActions[zpks].Reason)
		if err != nil {
			return
		}
	}
	// write "energy_status"
	err = en.Append(0xad, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.EnergyStatus)))
	if err != nil {
		return
	}
	for zjfb := range z.EnergyStatus {
		err = en.WriteFloat64(z.EnergyStatus[zjfb])
		if err != nil {
			return
		}
	}
	// write "health_status"
	err = en.Append(0xad, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.HealthStatus)))
	if err != nil {
		return
	}
	for zcxo := range z.HealthStatus {
		err = en.WriteFloat64(z.HealthStatus[zcxo])
		if err != nil {
			return
		}
	}
	// write "damage_cumulative_contrib"
	err = en.Append(0xb9, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.DamageCumulativeContrib)))
	if err != nil {
		return
	}
	for zeff := range z.DamageCumulativeContrib {
		err = en.WriteFloat64(z.DamageCumulativeContrib[zeff])
		if err != nil {
			return
		}
	}
	// write "active_time"
	err = en.Append(0xab, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.ActiveTime)
	if err != nil {
		return
	}
	// write "energy_spent"
	err = en.Append(0xac, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x73, 0x70, 0x65, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.EnergySpent)
	if err != nil {
		return
	}
	// write "energy_info"
	err = en.Append(0xab, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x69, 0x6e, 0x66, 0x6f)
	if err != nil {
		return err
	}
	err = z.EnergyInfo.EncodeMsg(en)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *CharacterResult) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 13
	// string "name"
	o = append(o, 0x8d, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "damage_events"
	o = append(o, 0xad, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.DamageEvents)))
	for zhct := range z.DamageEvents {
		o, err = z.DamageEvents[zhct].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "reaction_events"
	o = append(o, 0xaf, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.ReactionEvents)))
	for zcua := range z.ReactionEvents {
		o, err = z.ReactionEvents[zcua].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "action_events"
	o = append(o, 0xad, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.ActionEvents)))
	for zxhx := range z.ActionEvents {
		// map header, size 3
		// string "frame"
		o = append(o, 0x83, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
		o = msgp.AppendInt(o, z.ActionEvents[zxhx].Frame)
		// string "action_id"
		o = append(o, 0xa9, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64)
		o = msgp.AppendInt(o, z.ActionEvents[zxhx].ActionId)
		// string "action"
		o = append(o, 0xa6, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e)
		o = msgp.AppendString(o, z.ActionEvents[zxhx].Action)
	}
	// string "energy_events"
	o = append(o, 0xad, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.EnergyEvents)))
	for zlqf := range z.EnergyEvents {
		o, err = z.EnergyEvents[zlqf].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "heal_events"
	o = append(o, 0xab, 0x68, 0x65, 0x61, 0x6c, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.HealEvents)))
	for zdaf := range z.HealEvents {
		o, err = z.HealEvents[zdaf].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "failed_actions"
	o = append(o, 0xae, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.FailedActions)))
	for zpks := range z.FailedActions {
		// map header, size 3
		// string "start"
		o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
		o = msgp.AppendInt(o, z.FailedActions[zpks].Start)
		// string "end"
		o = append(o, 0xa3, 0x65, 0x6e, 0x64)
		o = msgp.AppendInt(o, z.FailedActions[zpks].End)
		// string "reason"
		o = append(o, 0xa6, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e)
		o = msgp.AppendString(o, z.FailedActions[zpks].Reason)
	}
	// string "energy_status"
	o = append(o, 0xad, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.EnergyStatus)))
	for zjfb := range z.EnergyStatus {
		o = msgp.AppendFloat64(o, z.EnergyStatus[zjfb])
	}
	// string "health_status"
	o = append(o, 0xad, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.HealthStatus)))
	for zcxo := range z.HealthStatus {
		o = msgp.AppendFloat64(o, z.HealthStatus[zcxo])
	}
	// string "damage_cumulative_contrib"
	o = append(o, 0xb9, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x75, 0x6d, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62)
	o = msgp.AppendArrayHeader(o, uint32(len(z.DamageCumulativeContrib)))
	for zeff := range z.DamageCumulativeContrib {
		o = msgp.AppendFloat64(o, z.DamageCumulativeContrib[zeff])
	}
	// string "active_time"
	o = append(o, 0xab, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65)
	o = msgp.AppendInt(o, z.ActiveTime)
	// string "energy_spent"
	o = append(o, 0xac, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x73, 0x70, 0x65, 0x6e, 0x74)
	o = msgp.AppendFloat64(o, z.EnergySpent)
	// string "energy_info"
	o = append(o, 0xab, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x69, 0x6e, 0x66, 0x6f)
	o, err = z.EnergyInfo.MarshalMsg(o)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *CharacterResult) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zjpj uint32
	zjpj, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zjpj > 0 {
		zjpj--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "damage_events":
			var zzpf uint32
			zzpf, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.DamageEvents) >= int(zzpf) {
				z.DamageEvents = (z.DamageEvents)[:zzpf]
			} else {
				z.DamageEvents = make([]DamageEvent, zzpf)
			}
			for zhct := range z.DamageEvents {
				bts, err = z.DamageEvents[zhct].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "reaction_events":
			var zrfe uint32
			zrfe, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.ReactionEvents) >= int(zrfe) {
				z.ReactionEvents = (z.ReactionEvents)[:zrfe]
			} else {
				z.ReactionEvents = make([]ReactionEvent, zrfe)
			}
			for zcua := range z.ReactionEvents {
				bts, err = z.ReactionEvents[zcua].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "action_events":
			var zgmo uint32
			zgmo, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.ActionEvents) >= int(zgmo) {
				z.ActionEvents = (z.ActionEvents)[:zgmo]
			} else {
				z.ActionEvents = make([]ActionEvent, zgmo)
			}
			for zxhx := range z.ActionEvents {
				var ztaf uint32
				ztaf, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for ztaf > 0 {
					ztaf--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "frame":
						z.ActionEvents[zxhx].Frame, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					case "action_id":
						z.ActionEvents[zxhx].ActionId, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					case "action":
						z.ActionEvents[zxhx].Action, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
			}
		case "energy_events":
			var zeth uint32
			zeth, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.EnergyEvents) >= int(zeth) {
				z.EnergyEvents = (z.EnergyEvents)[:zeth]
			} else {
				z.EnergyEvents = make([]EnergyEvent, zeth)
			}
			for zlqf := range z.EnergyEvents {
				bts, err = z.EnergyEvents[zlqf].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "heal_events":
			var zsbz uint32
			zsbz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.HealEvents) >= int(zsbz) {
				z.HealEvents = (z.HealEvents)[:zsbz]
			} else {
				z.HealEvents = make([]HealEvent, zsbz)
			}
			for zdaf := range z.HealEvents {
				bts, err = z.HealEvents[zdaf].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "failed_actions":
			var zrjx uint32
			zrjx, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.FailedActions) >= int(zrjx) {
				z.FailedActions = (z.FailedActions)[:zrjx]
			} else {
				z.FailedActions = make([]ActionFailInterval, zrjx)
			}
			for zpks := range z.FailedActions {
				var zawn uint32
				zawn, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for zawn > 0 {
					zawn--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "start":
						z.FailedActions[zpks].Start, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					case "end":
						z.FailedActions[zpks].End, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					case "reason":
						z.FailedActions[zpks].Reason, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
			}
		case "energy_status":
			var zwel uint32
			zwel, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.EnergyStatus) >= int(zwel) {
				z.EnergyStatus = (z.EnergyStatus)[:zwel]
			} else {
				z.EnergyStatus = make([]float64, zwel)
			}
			for zjfb := range z.EnergyStatus {
				z.EnergyStatus[zjfb], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "health_status":
			var zrbe uint32
			zrbe, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.HealthStatus) >= int(zrbe) {
				z.HealthStatus = (z.HealthStatus)[:zrbe]
			} else {
				z.HealthStatus = make([]float64, zrbe)
			}
			for zcxo := range z.HealthStatus {
				z.HealthStatus[zcxo], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "damage_cumulative_contrib":
			var zmfd uint32
			zmfd, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.DamageCumulativeContrib) >= int(zmfd) {
				z.DamageCumulativeContrib = (z.DamageCumulativeContrib)[:zmfd]
			} else {
				z.DamageCumulativeContrib = make([]float64, zmfd)
			}
			for zeff := range z.DamageCumulativeContrib {
				z.DamageCumulativeContrib[zeff], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "active_time":
			z.ActiveTime, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "energy_spent":
			z.EnergySpent, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "energy_info":
			bts, err = z.EnergyInfo.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *CharacterResult) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 14 + msgp.ArrayHeaderSize
	for zhct := range z.DamageEvents {
		s += z.DamageEvents[zhct].Msgsize()
	}
	s += 16 + msgp.ArrayHeaderSize
	for zcua := range z.ReactionEvents {
		s += z.ReactionEvents[zcua].Msgsize()
	}
	s += 14 + msgp.ArrayHeaderSize
	for zxhx := range z.ActionEvents {
		s += 1 + 6 + msgp.IntSize + 10 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.ActionEvents[zxhx].Action)
	}
	s += 14 + msgp.ArrayHeaderSize
	for zlqf := range z.EnergyEvents {
		s += z.EnergyEvents[zlqf].Msgsize()
	}
	s += 12 + msgp.ArrayHeaderSize
	for zdaf := range z.HealEvents {
		s += z.HealEvents[zdaf].Msgsize()
	}
	s += 15 + msgp.ArrayHeaderSize
	for zpks := range z.FailedActions {
		s += 1 + 6 + msgp.IntSize + 4 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.FailedActions[zpks].Reason)
	}
	s += 14 + msgp.ArrayHeaderSize + (len(z.EnergyStatus) * (msgp.Float64Size)) + 14 + msgp.ArrayHeaderSize + (len(z.HealthStatus) * (msgp.Float64Size)) + 26 + msgp.ArrayHeaderSize + (len(z.DamageCumulativeContrib) * (msgp.Float64Size)) + 12 + msgp.IntSize + 13 + msgp.Float64Size + 12 + z.EnergyInfo.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DamageEvent) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zelx uint32
	zelx, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zelx > 0 {
		zelx--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "action_id":
			z.ActionId, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "source":
			z.Source, err = dc.ReadString()
			if err != nil {
				return
			}
		case "target":
			z.Target, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "element":
			z.Element, err = dc.ReadString()
			if err != nil {
				return
			}
		case "reaction_modifier":
			{
				var zbal string
				zbal, err = dc.ReadString()
				z.ReactionModifier = ReactionModifier(zbal)
			}
			if err != nil {
				return
			}
		case "crit":
			z.Crit, err = dc.ReadBool()
			if err != nil {
				return
			}
		case "modifiers":
			var zjqz uint32
			zjqz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Modifiers) >= int(zjqz) {
				z.Modifiers = (z.Modifiers)[:zjqz]
			} else {
				z.Modifiers = make([]string, zjqz)
			}
			for zzdc := range z.Modifiers {
				z.Modifiers[zzdc], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "mitigation_modifier":
			z.MitigationModifier, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "damage":
			z.Damage, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *DamageEvent) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 10
	// write "frame"
	err = en.Append(0x8a, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Frame)
	if err != nil {
		return
	}
	// write "action_id"
	err = en.Append(0xa9, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.ActionId)
	if err != nil {
		return
	}
	// write "source"
	err = en.Append(0xa6, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Source)
	if err != nil {
		return
	}
	// write "target"
	err = en.Append(0xa6, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Target)
	if err != nil {
		return
	}
	// write "element"
	err = en.Append(0xa7, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Element)
	if err != nil {
		return
	}
	// write "reaction_modifier"
	err = en.Append(0xb1, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(string(z.ReactionModifier))
	if err != nil {
		return
	}
	// write "crit"
	err = en.Append(0xa4, 0x63, 0x72, 0x69, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.Crit)
	if err != nil {
		return
	}
	// write "modifiers"
	err = en.Append(0xa9, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Modifiers)))
	if err != nil {
		return
	}
	for zzdc := range z.Modifiers {
		err = en.WriteString(z.Modifiers[zzdc])
		if err != nil {
			return
		}
	}
	// write "mitigation_modifier"
	err = en.Append(0xb3, 0x6d, 0x69, 0x74, 0x69, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.MitigationModifier)
	if err != nil {
		return
	}
	// write "damage"
	err = en.Append(0xa6, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Damage)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *DamageEvent) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 10
	// string "frame"
	o = append(o, 0x8a, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	o = msgp.AppendInt(o, z.Frame)
	// string "action_id"
	o = append(o, 0xa9, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64)
	o = msgp.AppendInt(o, z.ActionId)
	// string "source"
	o = append(o, 0xa6, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	o = msgp.AppendString(o, z.Source)
	// string "target"
	o = append(o, 0xa6, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74)
	o = msgp.AppendInt(o, z.Target)
	// string "element"
	o = append(o, 0xa7, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74)
	o = msgp.AppendString(o, z.Element)
	// string "reaction_modifier"
	o = append(o, 0xb1, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72)
	o = msgp.AppendString(o, string(z.ReactionModifier))
	// string "crit"
	o = append(o, 0xa4, 0x63, 0x72, 0x69, 0x74)
	o = msgp.AppendBool(o, z.Crit)
	// string "modifiers"
	o = append(o, 0xa9, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Modifiers)))
	for zzdc := range z.Modifiers {
		o = msgp.AppendString(o, z.Modifiers[zzdc])
	}
	// string "mitigation_modifier"
	o = append(o, 0xb3, 0x6d, 0x69, 0x74, 0x69, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72)
	o = msgp.AppendFloat64(o, z.MitigationModifier)
	// string "damage"
	o = append(o, 0xa6, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65)
	o = msgp.AppendFloat64(o, z.Damage)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DamageEvent) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zkct uint32
	zkct, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zkct > 0 {
		zkct--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "action_id":
			z.ActionId, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "source":
			z.Source, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "target":
			z.Target, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "element":
			z.Element, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "reaction_modifier":
			{
				var ztmt string
				ztmt, bts, err = msgp.ReadStringBytes(bts)
				z.ReactionModifier = ReactionModifier(ztmt)
			}
			if err != nil {
				return
			}
		case "crit":
			z.Crit, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		case "modifiers":
			var ztco uint32
			ztco, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Modifiers) >= int(ztco) {
				z.Modifiers = (z.Modifiers)[:ztco]
			} else {
				z.Modifiers = make([]string, ztco)
			}
			for zzdc := range z.Modifiers {
				z.Modifiers[zzdc], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "mitigation_modifier":
			z.MitigationModifier, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "damage":
			z.Damage, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *DamageEvent) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 10 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.Source) + 7 + msgp.IntSize + 8 + msgp.StringPrefixSize + len(z.Element) + 18 + msgp.StringPrefixSize + len(string(z.ReactionModifier)) + 5 + msgp.BoolSize + 10 + msgp.ArrayHeaderSize
	for zzdc := range z.Modifiers {
		s += msgp.StringPrefixSize + len(z.Modifiers[zzdc])
	}
	s += 20 + msgp.Float64Size + 7 + msgp.Float64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *EnemyResult) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zare uint32
	zare, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zare > 0 {
		zare--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "reaction_status":
			var zljy uint32
			zljy, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.ReactionStatus) >= int(zljy) {
				z.ReactionStatus = (z.ReactionStatus)[:zljy]
			} else {
				z.ReactionStatus = make([]ReactionStatusInterval, zljy)
			}
			for zana := range z.ReactionStatus {
				var zixj uint32
				zixj, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zixj > 0 {
					zixj--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "start":
						z.ReactionStatus[zana].Start, err = dc.ReadInt()
						if err != nil {
							return
						}
					case "end":
						z.ReactionStatus[zana].End, err = dc.ReadInt()
						if err != nil {
							return
						}
					case "type":
						z.ReactionStatus[zana].Type, err = dc.ReadString()
						if err != nil {
							return
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
			}
		case "reaction_uptime":
			var zrsc uint32
			zrsc, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.ReactionUptime == nil && zrsc > 0 {
				z.ReactionUptime = make(map[string]int, zrsc)
			} else if len(z.ReactionUptime) > 0 {
				for key, _ := range z.ReactionUptime {
					delete(z.ReactionUptime, key)
				}
			}
			for zrsc > 0 {
				zrsc--
				var ztyy string
				var zinl int
				ztyy, err = dc.ReadString()
				if err != nil {
					return
				}
				zinl, err = dc.ReadInt()
				if err != nil {
					return
				}
				z.ReactionUptime[ztyy] = zinl
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *EnemyResult) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "reaction_status"
	err = en.Append(0x82, 0xaf, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.ReactionStatus)))
	if err != nil {
		return
	}
	for zana := range z.ReactionStatus {
		// map header, size 3
		// write "start"
		err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.ReactionStatus[zana].Start)
		if err != nil {
			return
		}
		// write "end"
		err = en.Append(0xa3, 0x65, 0x6e, 0x64)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.ReactionStatus[zana].End)
		if err != nil {
			return
		}
		// write "type"
		err = en.Append(0xa4, 0x74, 0x79, 0x70, 0x65)
		if err != nil {
			return err
		}
		err = en.WriteString(z.ReactionStatus[zana].Type)
		if err != nil {
			return
		}
	}
	// write "reaction_uptime"
	err = en.Append(0xaf, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.ReactionUptime)))
	if err != nil {
		return
	}
	for ztyy, zinl := range z.ReactionUptime {
		err = en.WriteString(ztyy)
		if err != nil {
			return
		}
		err = en.WriteInt(zinl)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *EnemyResult) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "reaction_status"
	o = append(o, 0x82, 0xaf, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.ReactionStatus)))
	for zana := range z.ReactionStatus {
		// map header, size 3
		// string "start"
		o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
		o = msgp.AppendInt(o, z.ReactionStatus[zana].Start)
		// string "end"
		o = append(o, 0xa3, 0x65, 0x6e, 0x64)
		o = msgp.AppendInt(o, z.ReactionStatus[zana].End)
		// string "type"
		o = append(o, 0xa4, 0x74, 0x79, 0x70, 0x65)
		o = msgp.AppendString(o, z.ReactionStatus[zana].Type)
	}
	// string "reaction_uptime"
	o = append(o, 0xaf, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65)
	o = msgp.AppendMapHeader(o, uint32(len(z.ReactionUptime)))
	for ztyy, zinl := range z.ReactionUptime {
		o = msgp.AppendString(o, ztyy)
		o = msgp.AppendInt(o, zinl)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *EnemyResult) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zctn uint32
	zctn, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zctn > 0 {
		zctn--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "reaction_status":
			var zswy uint32
			zswy, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.ReactionStatus) >= int(zswy) {
				z.ReactionStatus = (z.ReactionStatus)[:zswy]
			} else {
				z.ReactionStatus = make([]ReactionStatusInterval, zswy)
			}
			for zana := range z.ReactionStatus {
				var znsg uint32
				znsg, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for znsg > 0 {
					znsg--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "start":
						z.ReactionStatus[zana].Start, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					case "end":
						z.ReactionStatus[zana].End, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					case "type":
						z.ReactionStatus[zana].Type, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
			}
		case "reaction_uptime":
			var zrus uint32
			zrus, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.ReactionUptime == nil && zrus > 0 {
				z.ReactionUptime = make(map[string]int, zrus)
			} else if len(z.ReactionUptime) > 0 {
				for key, _ := range z.ReactionUptime {
					delete(z.ReactionUptime, key)
				}
			}
			for zrus > 0 {
				var ztyy string
				var zinl int
				zrus--
				ztyy, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				zinl, bts, err = msgp.ReadIntBytes(bts)
				if err != nil {
					return
				}
				z.ReactionUptime[ztyy] = zinl
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *EnemyResult) Msgsize() (s int) {
	s = 1 + 16 + msgp.ArrayHeaderSize
	for zana := range z.ReactionStatus {
		s += 1 + 6 + msgp.IntSize + 4 + msgp.IntSize + 5 + msgp.StringPrefixSize + len(z.ReactionStatus[zana].Type)
	}
	s += 16 + msgp.MapHeaderSize
	if z.ReactionUptime != nil {
		for ztyy, zinl := range z.ReactionUptime {
			_ = zinl
			s += msgp.StringPrefixSize + len(ztyy) + msgp.IntSize
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *EnergyEvent) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zsvm uint32
	zsvm, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zsvm > 0 {
		zsvm--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "source":
			z.Source, err = dc.ReadString()
			if err != nil {
				return
			}
		case "field_status":
			{
				var zaoz string
				zaoz, err = dc.ReadString()
				z.FieldStatus = FieldStatus(zaoz)
			}
			if err != nil {
				return
			}
		case "gained":
			z.Gained, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "wasted":
			z.Wasted, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "current":
			z.Current, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *EnergyEvent) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 6
	// write "frame"
	err = en.Append(0x86, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Frame)
	if err != nil {
		return
	}
	// write "source"
	err = en.Append(0xa6, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Source)
	if err != nil {
		return
	}
	// write "field_status"
	err = en.Append(0xac, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteString(string(z.FieldStatus))
	if err != nil {
		return
	}
	// write "gained"
	err = en.Append(0xa6, 0x67, 0x61, 0x69, 0x6e, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Gained)
	if err != nil {
		return
	}
	// write "wasted"
	err = en.Append(0xa6, 0x77, 0x61, 0x73, 0x74, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Wasted)
	if err != nil {
		return
	}
	// write "current"
	err = en.Append(0xa7, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Current)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *EnergyEvent) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "frame"
	o = append(o, 0x86, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	o = msgp.AppendInt(o, z.Frame)
	// string "source"
	o = append(o, 0xa6, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	o = msgp.AppendString(o, z.Source)
	// string "field_status"
	o = append(o, 0xac, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	o = msgp.AppendString(o, string(z.FieldStatus))
	// string "gained"
	o = append(o, 0xa6, 0x67, 0x61, 0x69, 0x6e, 0x65, 0x64)
	o = msgp.AppendFloat64(o, z.Gained)
	// string "wasted"
	o = append(o, 0xa6, 0x77, 0x61, 0x73, 0x74, 0x65, 0x64)
	o = msgp.AppendFloat64(o, z.Wasted)
	// string "current"
	o = append(o, 0xa7, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74)
	o = msgp.AppendFloat64(o, z.Current)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *EnergyEvent) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zfzb uint32
	zfzb, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zfzb > 0 {
		zfzb--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "source":
			z.Source, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "field_status":
			{
				var zsbo string
				zsbo, bts, err = msgp.ReadStringBytes(bts)
				z.FieldStatus = FieldStatus(zsbo)
			}
			if err != nil {
				return
			}
		case "gained":
			z.Gained, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "wasted":
			z.Wasted, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "current":
			z.Current, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *EnergyEvent) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.Source) + 13 + msgp.StringPrefixSize + len(string(z.FieldStatus)) + 7 + msgp.Float64Size + 7 + msgp.Float64Size + 8 + msgp.Float64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *EnergyInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var ztls uint32
	ztls, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for ztls > 0 {
		ztls--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "raw_particles_per_burst":
			var zmvo uint32
			zmvo, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.RawParticlesPerBurst) >= int(zmvo) {
				z.RawParticlesPerBurst = (z.RawParticlesPerBurst)[:zmvo]
			} else {
				z.RawParticlesPerBurst = make([]float64, zmvo)
			}
			for zjif := range z.RawParticlesPerBurst {
				z.RawParticlesPerBurst[zjif], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "flat_energy_per_burst":
			var zigk uint32
			zigk, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.FlatEnergyPerBurst) >= int(zigk) {
				z.FlatEnergyPerBurst = (z.FlatEnergyPerBurst)[:zigk]
			} else {
				z.FlatEnergyPerBurst = make([]float64, zigk)
			}
			for zqgz := range z.FlatEnergyPerBurst {
				z.FlatEnergyPerBurst[zqgz], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "weighted_er":
			var zopb uint32
			zopb, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.WeightedER) >= int(zopb) {
				z.WeightedER = (z.WeightedER)[:zopb]
			} else {
				z.WeightedER = make([]float64, zopb)
			}
			for zsnw := range z.WeightedER {
				z.WeightedER[zsnw], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *EnergyInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "raw_particles_per_burst"
	err = en.Append(0x83, 0xb7, 0x72, 0x61, 0x77, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x62, 0x75, 0x72, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.RawParticlesPerBurst)))
	if err != nil {
		return
	}
	for zjif := range z.RawParticlesPerBurst {
		err = en.WriteFloat64(z.RawParticlesPerBurst[zjif])
		if err != nil {
			return
		}
	}
	// write "flat_energy_per_burst"
	err = en.Append(0xb5, 0x66, 0x6c, 0x61, 0x74, 0x5f, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x62, 0x75, 0x72, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.FlatEnergyPerBurst)))
	if err != nil {
		return
	}
	for zqgz := range z.FlatEnergyPerBurst {
		err = en.WriteFloat64(z.FlatEnergyPerBurst[zqgz])
		if err != nil {
			return
		}
	}
	// write "weighted_er"
	err = en.Append(0xab, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x64, 0x5f, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.WeightedER)))
	if err != nil {
		return
	}
	for zsnw := range z.WeightedER {
		err = en.WriteFloat64(z.WeightedER[zsnw])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *EnergyInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "raw_particles_per_burst"
	o = append(o, 0x83, 0xb7, 0x72, 0x61, 0x77, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x62, 0x75, 0x72, 0x73, 0x74)
	o = msgp.AppendArrayHeader(o, uint32(len(z.RawParticlesPerBurst)))
	for zjif := range z.RawParticlesPerBurst {
		o = msgp.AppendFloat64(o, z.RawParticlesPerBurst[zjif])
	}
	// string "flat_energy_per_burst"
	o = append(o, 0xb5, 0x66, 0x6c, 0x61, 0x74, 0x5f, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x62, 0x75, 0x72, 0x73, 0x74)
	o = msgp.AppendArrayHeader(o, uint32(len(z.FlatEnergyPerBurst)))
	for zqgz := range z.FlatEnergyPerBurst {
		o = msgp.AppendFloat64(o, z.FlatEnergyPerBurst[zqgz])
	}
	// string "weighted_er"
	o = append(o, 0xab, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x64, 0x5f, 0x65, 0x72)
	o = msgp.AppendArrayHeader(o, uint32(len(z.WeightedER)))
	for zsnw := range z.WeightedER {
		o = msgp.AppendFloat64(o, z.WeightedER[zsnw])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *EnergyInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zuop uint32
	zuop, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zuop > 0 {
		zuop--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "raw_particles_per_burst":
			var zedl uint32
			zedl, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.RawParticlesPerBurst) >= int(zedl) {
				z.RawParticlesPerBurst = (z.RawParticlesPerBurst)[:zedl]
			} else {
				z.RawParticlesPerBurst = make([]float64, zedl)
			}
			for zjif := range z.RawParticlesPerBurst {
				z.RawParticlesPerBurst[zjif], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "flat_energy_per_burst":
			var zupd uint32
			zupd, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.FlatEnergyPerBurst) >= int(zupd) {
				z.FlatEnergyPerBurst = (z.FlatEnergyPerBurst)[:zupd]
			} else {
				z.FlatEnergyPerBurst = make([]float64, zupd)
			}
			for zqgz := range z.FlatEnergyPerBurst {
				z.FlatEnergyPerBurst[zqgz], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "weighted_er":
			var zome uint32
			zome, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.WeightedER) >= int(zome) {
				z.WeightedER = (z.WeightedER)[:zome]
			} else {
				z.WeightedER = make([]float64, zome)
			}
			for zsnw := range z.WeightedER {
				z.WeightedER[zsnw], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *EnergyInfo) Msgsize() (s int) {
	s = 1 + 24 + msgp.ArrayHeaderSize + (len(z.RawParticlesPerBurst) * (msgp.Float64Size)) + 22 + msgp.ArrayHeaderSize + (len(z.FlatEnergyPerBurst) * (msgp.Float64Size)) + 12 + msgp.ArrayHeaderSize + (len(z.WeightedER) * (msgp.Float64Size))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *FieldStatus) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zrvj string
		zrvj, err = dc.ReadString()
		(*z) = FieldStatus(zrvj)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z FieldStatus) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteString(string(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z FieldStatus) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendString(o, string(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *FieldStatus) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zarz string
		zarz, bts, err = msgp.ReadStringBytes(bts)
		(*z) = FieldStatus(zarz)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z FieldStatus) Msgsize() (s int) {
	s = msgp.StringPrefixSize + len(string(z))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *HealEvent) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zknt uint32
	zknt, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zknt > 0 {
		zknt--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "source":
			z.Source, err = dc.ReadString()
			if err != nil {
				return
			}
		case "target":
			z.Target, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "heal":
			z.Heal, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *HealEvent) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "frame"
	err = en.Append(0x84, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Frame)
	if err != nil {
		return
	}
	// write "source"
	err = en.Append(0xa6, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Source)
	if err != nil {
		return
	}
	// write "target"
	err = en.Append(0xa6, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Target)
	if err != nil {
		return
	}
	// write "heal"
	err = en.Append(0xa4, 0x68, 0x65, 0x61, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Heal)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *HealEvent) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "frame"
	o = append(o, 0x84, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	o = msgp.AppendInt(o, z.Frame)
	// string "source"
	o = append(o, 0xa6, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	o = msgp.AppendString(o, z.Source)
	// string "target"
	o = append(o, 0xa6, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74)
	o = msgp.AppendInt(o, z.Target)
	// string "heal"
	o = append(o, 0xa4, 0x68, 0x65, 0x61, 0x6c)
	o = msgp.AppendFloat64(o, z.Heal)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *HealEvent) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zxye uint32
	zxye, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zxye > 0 {
		zxye--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "source":
			z.Source, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "target":
			z.Target, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "heal":
			z.Heal, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *HealEvent) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.Source) + 7 + msgp.IntSize + 5 + msgp.Float64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ReactionEvent) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zucw uint32
	zucw, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zucw > 0 {
		zucw--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "source":
			z.Source, err = dc.ReadString()
			if err != nil {
				return
			}
		case "target":
			z.Target, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "reaction":
			z.Reaction, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ReactionEvent) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "frame"
	err = en.Append(0x84, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Frame)
	if err != nil {
		return
	}
	// write "source"
	err = en.Append(0xa6, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Source)
	if err != nil {
		return
	}
	// write "target"
	err = en.Append(0xa6, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Target)
	if err != nil {
		return
	}
	// write "reaction"
	err = en.Append(0xa8, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Reaction)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ReactionEvent) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "frame"
	o = append(o, 0x84, 0xa5, 0x66, 0x72, 0x61, 0x6d, 0x65)
	o = msgp.AppendInt(o, z.Frame)
	// string "source"
	o = append(o, 0xa6, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	o = msgp.AppendString(o, z.Source)
	// string "target"
	o = append(o, 0xa6, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74)
	o = msgp.AppendInt(o, z.Target)
	// string "reaction"
	o = append(o, 0xa8, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Reaction)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ReactionEvent) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zlsx uint32
	zlsx, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zlsx > 0 {
		zlsx--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "frame":
			z.Frame, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "source":
			z.Source, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "target":
			z.Target, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "reaction":
			z.Reaction, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ReactionEvent) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 7 + msgp.StringPrefixSize + len(z.Source) + 7 + msgp.IntSize + 9 + msgp.StringPrefixSize + len(z.Reaction)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ReactionModifier) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zbgy string
		zbgy, err = dc.ReadString()
		(*z) = ReactionModifier(zbgy)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ReactionModifier) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteString(string(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ReactionModifier) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendString(o, string(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ReactionModifier) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zrao string
		zrao, bts, err = msgp.ReadStringBytes(bts)
		(*z) = ReactionModifier(zrao)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ReactionModifier) Msgsize() (s int) {
	s = msgp.StringPrefixSize + len(string(z))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ReactionStatusInterval) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zmbt uint32
	zmbt, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zmbt > 0 {
		zmbt--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "end":
			z.End, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "type":
			z.Type, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ReactionStatusInterval) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "start"
	err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Start)
	if err != nil {
		return
	}
	// write "end"
	err = en.Append(0xa3, 0x65, 0x6e, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.End)
	if err != nil {
		return
	}
	// write "type"
	err = en.Append(0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Type)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ReactionStatusInterval) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "start"
	o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	o = msgp.AppendInt(o, z.Start)
	// string "end"
	o = append(o, 0xa3, 0x65, 0x6e, 0x64)
	o = msgp.AppendInt(o, z.End)
	// string "type"
	o = append(o, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendString(o, z.Type)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ReactionStatusInterval) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zvls uint32
	zvls, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zvls > 0 {
		zvls--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "end":
			z.End, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "type":
			z.Type, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ReactionStatusInterval) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 4 + msgp.IntSize + 5 + msgp.StringPrefixSize + len(z.Type)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Result) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zdtr uint32
	zdtr, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zdtr > 0 {
		zdtr--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "seed":
			z.Seed, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "duration":
			z.Duration, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "total_damage":
			z.TotalDamage, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "dps":
			z.DPS, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "damage_buckets":
			var zzqm uint32
			zzqm, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.DamageBuckets) >= int(zzqm) {
				z.DamageBuckets = (z.DamageBuckets)[:zzqm]
			} else {
				z.DamageBuckets = make([]float64, zzqm)
			}
			for zjfj := range z.DamageBuckets {
				z.DamageBuckets[zjfj], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "active_characters":
			var zdqi uint32
			zdqi, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.ActiveCharacters) >= int(zdqi) {
				z.ActiveCharacters = (z.ActiveCharacters)[:zdqi]
			} else {
				z.ActiveCharacters = make([]ActiveCharacterInterval, zdqi)
			}
			for zzak := range z.ActiveCharacters {
				var zyco uint32
				zyco, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zyco > 0 {
					zyco--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "start":
						z.ActiveCharacters[zzak].Start, err = dc.ReadInt()
						if err != nil {
							return
						}
					case "end":
						z.ActiveCharacters[zzak].End, err = dc.ReadInt()
						if err != nil {
							return
						}
					case "character":
						z.ActiveCharacters[zzak].Character, err = dc.ReadInt()
						if err != nil {
							return
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
			}
		case "damage_mitigation":
			var zhgh uint32
			zhgh, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.DamageMitigation) >= int(zhgh) {
				z.DamageMitigation = (z.DamageMitigation)[:zhgh]
			} else {
				z.DamageMitigation = make([]float64, zhgh)
			}
			for zbtz := range z.DamageMitigation {
				z.DamageMitigation[zbtz], err = dc.ReadFloat64()
				if err != nil {
					return
				}
			}
		case "shield_results":
			err = z.ShieldResults.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "characters":
			var zovg uint32
			zovg, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Characters) >= int(zovg) {
				z.Characters = (z.Characters)[:zovg]
			} else {
				z.Characters = make([]CharacterResult, zovg)
			}
			for zsym := range z.Characters {
				err = z.Characters[zsym].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "enemies":
			var zsey uint32
			zsey, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Enemies) >= int(zsey) {
				z.Enemies = (z.Enemies)[:zsey]
			} else {
				z.Enemies = make([]EnemyResult, zsey)
			}
			for zgeu := range z.Enemies {
				err = z.Enemies[zgeu].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "target_overlap":
			z.TargetOverlap, err = dc.ReadBool()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Result) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 11
	// write "seed"
	err = en.Append(0x8b, 0xa4, 0x73, 0x65, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.Seed)
	if err != nil {
		return
	}
	// write "duration"
	err = en.Append(0xa8, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Duration)
	if err != nil {
		return
	}
	// write "total_damage"
	err = en.Append(0xac, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.TotalDamage)
	if err != nil {
		return
	}
	// write "dps"
	err = en.Append(0xa3, 0x64, 0x70, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.DPS)
	if err != nil {
		return
	}
	// write "damage_buckets"
	err = en.Append(0xae, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.DamageBuckets)))
	if err != nil {
		return
	}
	for zjfj := range z.DamageBuckets {
		err = en.WriteFloat64(z.DamageBuckets[zjfj])
		if err != nil {
			return
		}
	}
	// write "active_characters"
	err = en.Append(0xb1, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.ActiveCharacters)))
	if err != nil {
		return
	}
	for zzak := range z.ActiveCharacters {
		// map header, size 3
		// write "start"
		err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.ActiveCharacters[zzak].Start)
		if err != nil {
			return
		}
		// write "end"
		err = en.Append(0xa3, 0x65, 0x6e, 0x64)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.ActiveCharacters[zzak].End)
		if err != nil {
			return
		}
		// write "character"
		err = en.Append(0xa9, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72)
		if err != nil {
			return err
		}
		err = en.WriteInt(z.ActiveCharacters[zzak].Character)
		if err != nil {
			return
		}
	}
	// write "damage_mitigation"
	err = en.Append(0xb1, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x69, 0x74, 0x69, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.DamageMitigation)))
	if err != nil {
		return
	}
	for zbtz := range z.DamageMitigation {
		err = en.WriteFloat64(z.DamageMitigation[zbtz])
		if err != nil {
			return
		}
	}
	// write "shield_results"
	err = en.Append(0xae, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73)
	if err != nil {
		return err
	}
	err = z.ShieldResults.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "characters"
	err = en.Append(0xaa, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Characters)))
	if err != nil {
		return
	}
	for zsym := range z.Characters {
		err = z.Characters[zsym].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "enemies"
	err = en.Append(0xa7, 0x65, 0x6e, 0x65, 0x6d, 0x69, 0x65, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Enemies)))
	if err != nil {
		return
	}
	for zgeu := range z.Enemies {
		err = z.Enemies[zgeu].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "target_overlap"
	err = en.Append(0xae, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.TargetOverlap)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Result) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 11
	// string "seed"
	o = append(o, 0x8b, 0xa4, 0x73, 0x65, 0x65, 0x64)
	o = msgp.AppendUint64(o, z.Seed)
	// string "duration"
	o = append(o, 0xa8, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendInt(o, z.Duration)
	// string "total_damage"
	o = append(o, 0xac, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65)
	o = msgp.AppendFloat64(o, z.TotalDamage)
	// string "dps"
	o = append(o, 0xa3, 0x64, 0x70, 0x73)
	o = msgp.AppendFloat64(o, z.DPS)
	// string "damage_buckets"
	o = append(o, 0xae, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.DamageBuckets)))
	for zjfj := range z.DamageBuckets {
		o = msgp.AppendFloat64(o, z.DamageBuckets[zjfj])
	}
	// string "active_characters"
	o = append(o, 0xb1, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.ActiveCharacters)))
	for zzak := range z.ActiveCharacters {
		// map header, size 3
		// string "start"
		o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
		o = msgp.AppendInt(o, z.ActiveCharacters[zzak].Start)
		// string "end"
		o = append(o, 0xa3, 0x65, 0x6e, 0x64)
		o = msgp.AppendInt(o, z.ActiveCharacters[zzak].End)
		// string "character"
		o = append(o, 0xa9, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72)
		o = msgp.AppendInt(o, z.ActiveCharacters[zzak].Character)
	}
	// string "damage_mitigation"
	o = append(o, 0xb1, 0x64, 0x61, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x69, 0x74, 0x69, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendArrayHeader(o, uint32(len(z.DamageMitigation)))
	for zbtz := range z.DamageMitigation {
		o = msgp.AppendFloat64(o, z.DamageMitigation[zbtz])
	}
	// string "shield_results"
	o = append(o, 0xae, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73)
	o, err = z.ShieldResults.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "characters"
	o = append(o, 0xaa, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Characters)))
	for zsym := range z.Characters {
		o, err = z.Characters[zsym].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "enemies"
	o = append(o, 0xa7, 0x65, 0x6e, 0x65, 0x6d, 0x69, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Enemies)))
	for zgeu := range z.Enemies {
		o, err = z.Enemies[zgeu].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "target_overlap"
	o = append(o, 0xae, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x70)
	o = msgp.AppendBool(o, z.TargetOverlap)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Result) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcjp uint32
	zcjp, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcjp > 0 {
		zcjp--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "seed":
			z.Seed, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "duration":
			z.Duration, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "total_damage":
			z.TotalDamage, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "dps":
			z.DPS, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "damage_buckets":
			var zjhy uint32
			zjhy, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.DamageBuckets) >= int(zjhy) {
				z.DamageBuckets = (z.DamageBuckets)[:zjhy]
			} else {
				z.DamageBuckets = make([]float64, zjhy)
			}
			for zjfj := range z.DamageBuckets {
				z.DamageBuckets[zjfj], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "active_characters":
			var znuf uint32
			znuf, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.ActiveCharacters) >= int(znuf) {
				z.ActiveCharacters = (z.ActiveCharacters)[:znuf]
			} else {
				z.ActiveCharacters = make([]ActiveCharacterInterval, znuf)
			}
			for zzak := range z.ActiveCharacters {
				var znjj uint32
				znjj, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for znjj > 0 {
					znjj--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "start":
						z.ActiveCharacters[zzak].Start, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					case "end":
						z.ActiveCharacters[zzak].End, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					case "character":
						z.ActiveCharacters[zzak].Character, bts, err = msgp.ReadIntBytes(bts)
						if err != nil {
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
			}
		case "damage_mitigation":
			var zhhj uint32
			zhhj, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.DamageMitigation) >= int(zhhj) {
				z.DamageMitigation = (z.DamageMitigation)[:zhhj]
			} else {
				z.DamageMitigation = make([]float64, zhhj)
			}
			for zbtz := range z.DamageMitigation {
				z.DamageMitigation[zbtz], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
			}
		case "shield_results":
			bts, err = z.ShieldResults.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "characters":
			var zuvr uint32
			zuvr, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Characters) >= int(zuvr) {
				z.Characters = (z.Characters)[:zuvr]
			} else {
				z.Characters = make([]CharacterResult, zuvr)
			}
			for zsym := range z.Characters {
				bts, err = z.Characters[zsym].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "enemies":
			var zusq uint32
			zusq, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Enemies) >= int(zusq) {
				z.Enemies = (z.Enemies)[:zusq]
			} else {
				z.Enemies = make([]EnemyResult, zusq)
			}
			for zgeu := range z.Enemies {
				bts, err = z.Enemies[zgeu].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "target_overlap":
			z.TargetOverlap, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Result) Msgsize() (s int) {
	s = 1 + 5 + msgp.Uint64Size + 9 + msgp.IntSize + 13 + msgp.Float64Size + 4 + msgp.Float64Size + 15 + msgp.ArrayHeaderSize + (len(z.DamageBuckets) * (msgp.Float64Size)) + 18 + msgp.ArrayHeaderSize + (len(z.ActiveCharacters) * (21 + msgp.IntSize + msgp.IntSize + msgp.IntSize)) + 18 + msgp.ArrayHeaderSize + (len(z.DamageMitigation) * (msgp.Float64Size)) + 15 + z.ShieldResults.Msgsize() + 11 + msgp.ArrayHeaderSize
	for zsym := range z.Characters {
		s += z.Characters[zsym].Msgsize()
	}
	s += 8 + msgp.ArrayHeaderSize
	for zgeu := range z.Enemies {
		s += z.Enemies[zgeu].Msgsize()
	}
	s += 15 + msgp.BoolSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ShieldInterval) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zpyv uint32
	zpyv, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zpyv > 0 {
		zpyv--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "end":
			z.End, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "hp":
			var zlur uint32
			zlur, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.HP == nil && zlur > 0 {
				z.HP = make(map[string]float64, zlur)
			} else if len(z.HP) > 0 {
				for key, _ := range z.HP {
					delete(z.HP, key)
				}
			}
			for zlur > 0 {
				zlur--
				var zfgq string
				var zvml float64
				zfgq, err = dc.ReadString()
				if err != nil {
					return
				}
				zvml, err = dc.ReadFloat64()
				if err != nil {
					return
				}
				z.HP[zfgq] = zvml
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ShieldInterval) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "start"
	err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Start)
	if err != nil {
		return
	}
	// write "end"
	err = en.Append(0xa3, 0x65, 0x6e, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.End)
	if err != nil {
		return
	}
	// write "hp"
	err = en.Append(0xa2, 0x68, 0x70)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.HP)))
	if err != nil {
		return
	}
	for zfgq, zvml := range z.HP {
		err = en.WriteString(zfgq)
		if err != nil {
			return
		}
		err = en.WriteFloat64(zvml)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ShieldInterval) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "start"
	o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	o = msgp.AppendInt(o, z.Start)
	// string "end"
	o = append(o, 0xa3, 0x65, 0x6e, 0x64)
	o = msgp.AppendInt(o, z.End)
	// string "hp"
	o = append(o, 0xa2, 0x68, 0x70)
	o = msgp.AppendMapHeader(o, uint32(len(z.HP)))
	for zfgq, zvml := range z.HP {
		o = msgp.AppendString(o, zfgq)
		o = msgp.AppendFloat64(o, zvml)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ShieldInterval) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zupi uint32
	zupi, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zupi > 0 {
		zupi--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "end":
			z.End, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "hp":
			var zfvi uint32
			zfvi, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.HP == nil && zfvi > 0 {
				z.HP = make(map[string]float64, zfvi)
			} else if len(z.HP) > 0 {
				for key, _ := range z.HP {
					delete(z.HP, key)
				}
			}
			for zfvi > 0 {
				var zfgq string
				var zvml float64
				zfvi--
				zfgq, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				zvml, bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
				z.HP[zfgq] = zvml
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ShieldInterval) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 4 + msgp.IntSize + 3 + msgp.MapHeaderSize
	if z.HP != nil {
		for zfgq, zvml := range z.HP {
			_ = zvml
			s += msgp.StringPrefixSize + len(zfgq) + msgp.Float64Size
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ShieldResult) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zziv uint32
	zziv, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zziv > 0 {
		zziv--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "shields":
			var zabj uint32
			zabj, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Shields) >= int(zabj) {
				z.Shields = (z.Shields)[:zabj]
			} else {
				z.Shields = make([]ShieldStats, zabj)
			}
			for zzrg := range z.Shields {
				var zmlx uint32
				zmlx, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zmlx > 0 {
					zmlx--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "name":
						z.Shields[zzrg].Name, err = dc.ReadString()
						if err != nil {
							return
						}
					case "intervals":
						var zvbw uint32
						zvbw, err = dc.ReadArrayHeader()
						if err != nil {
							return
						}
						if cap(z.Shields[zzrg].Intervals) >= int(zvbw) {
							z.Shields[zzrg].Intervals = (z.Shields[zzrg].Intervals)[:zvbw]
						} else {
							z.Shields[zzrg].Intervals = make([]ShieldInterval, zvbw)
						}
						for zbmy := range z.Shields[zzrg].Intervals {
							err = z.Shields[zzrg].Intervals[zbmy].DecodeMsg(dc)
							if err != nil {
								return
							}
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
			}
		case "effective_shield":
			var zgvb uint32
			zgvb, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.EffectiveShield == nil && zgvb > 0 {
				z.EffectiveShield = make(map[string][]ShieldSingleInterval, zgvb)
			} else if len(z.EffectiveShield) > 0 {
				for key, _ := range z.EffectiveShield {
					delete(z.EffectiveShield, key)
				}
			}
			for zgvb > 0 {
				zgvb--
				var zarl string
				var zctz []ShieldSingleInterval
				zarl, err = dc.ReadString()
				if err != nil {
					return
				}
				var zqzg uint32
				zqzg, err = dc.ReadArrayHeader()
				if err != nil {
					return
				}
				if cap(zctz) >= int(zqzg) {
					zctz = (zctz)[:zqzg]
				} else {
					zctz = make([]ShieldSingleInterval, zqzg)
				}
				for zljl := range zctz {
					var zexy uint32
					zexy, err = dc.ReadMapHeader()
					if err != nil {
						return
					}
					for zexy > 0 {
						zexy--
						field, err = dc.ReadMapKeyPtr()
						if err != nil {
							return
						}
						switch msgp.UnsafeString(field) {
						case "start":
							zctz[zljl].Start, err = dc.ReadInt()
							if err != nil {
								return
							}
						case "end":
							zctz[zljl].End, err = dc.ReadInt()
							if err != nil {
								return
							}
						case "hp":
							zctz[zljl].HP, err = dc.ReadFloat64()
							if err != nil {
								return
							}
						default:
							err = dc.Skip()
							if err != nil {
								return
							}
						}
					}
				}
				z.EffectiveShield[zarl] = zctz
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ShieldResult) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "shields"
	err = en.Append(0x82, 0xa7, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Shields)))
	if err != nil {
		return
	}
	for zzrg := range z.Shields {
		// map header, size 2
		// write "name"
		err = en.Append(0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
		if err != nil {
			return err
		}
		err = en.WriteString(z.Shields[zzrg].Name)
		if err != nil {
			return
		}
		// write "intervals"
		err = en.Append(0xa9, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x73)
		if err != nil {
			return err
		}
		err = en.WriteArrayHeader(uint32(len(z.Shields[zzrg].Intervals)))
		if err != nil {
			return
		}
		for zbmy := range z.Shields[zzrg].Intervals {
			err = z.Shields[zzrg].Intervals[zbmy].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	// write "effective_shield"
	err = en.Append(0xb0, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.EffectiveShield)))
	if err != nil {
		return
	}
	for zarl, zctz := range z.EffectiveShield {
		err = en.WriteString(zarl)
		if err != nil {
			return
		}
		err = en.WriteArrayHeader(uint32(len(zctz)))
		if err != nil {
			return
		}
		for zljl := range zctz {
			// map header, size 3
			// write "start"
			err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
			if err != nil {
				return err
			}
			err = en.WriteInt(zctz[zljl].Start)
			if err != nil {
				return
			}
			// write "end"
			err = en.Append(0xa3, 0x65, 0x6e, 0x64)
			if err != nil {
				return err
			}
			err = en.WriteInt(zctz[zljl].End)
			if err != nil {
				return
			}
			// write "hp"
			err = en.Append(0xa2, 0x68, 0x70)
			if err != nil {
				return err
			}
			err = en.WriteFloat64(zctz[zljl].HP)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ShieldResult) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "shields"
	o = append(o, 0x82, 0xa7, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Shields)))
	for zzrg := range z.Shields {
		// map header, size 2
		// string "name"
		o = append(o, 0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
		o = msgp.AppendString(o, z.Shields[zzrg].Name)
		// string "intervals"
		o = append(o, 0xa9, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x73)
		o = msgp.AppendArrayHeader(o, uint32(len(z.Shields[zzrg].Intervals)))
		for zbmy := range z.Shields[zzrg].Intervals {
			o, err = z.Shields[zzrg].Intervals[zbmy].MarshalMsg(o)
			if err != nil {
				return
			}
		}
	}
	// string "effective_shield"
	o = append(o, 0xb0, 0x65, 0x66, 0x66, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x73, 0x68, 0x69, 0x65, 0x6c, 0x64)
	o = msgp.AppendMapHeader(o, uint32(len(z.EffectiveShield)))
	for zarl, zctz := range z.EffectiveShield {
		o = msgp.AppendString(o, zarl)
		o = msgp.AppendArrayHeader(o, uint32(len(zctz)))
		for zljl := range zctz {
			// map header, size 3
			// string "start"
			o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
			o = msgp.AppendInt(o, zctz[zljl].Start)
			// string "end"
			o = append(o, 0xa3, 0x65, 0x6e, 0x64)
			o = msgp.AppendInt(o, zctz[zljl].End)
			// string "hp"
			o = append(o, 0xa2, 0x68, 0x70)
			o = msgp.AppendFloat64(o, zctz[zljl].HP)
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ShieldResult) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zakb uint32
	zakb, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zakb > 0 {
		zakb--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "shields":
			var zsdj uint32
			zsdj, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Shields) >= int(zsdj) {
				z.Shields = (z.Shields)[:zsdj]
			} else {
				z.Shields = make([]ShieldStats, zsdj)
			}
			for zzrg := range z.Shields {
				var zsgp uint32
				zsgp, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for zsgp > 0 {
					zsgp--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "name":
						z.Shields[zzrg].Name, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "intervals":
						var zngc uint32
						zngc, bts, err = msgp.ReadArrayHeaderBytes(bts)
						if err != nil {
							return
						}
						if cap(z.Shields[zzrg].Intervals) >= int(zngc) {
							z.Shields[zzrg].Intervals = (z.Shields[zzrg].Intervals)[:zngc]
						} else {
							z.Shields[zzrg].Intervals = make([]ShieldInterval, zngc)
						}
						for zbmy := range z.Shields[zzrg].Intervals {
							bts, err = z.Shields[zzrg].Intervals[zbmy].UnmarshalMsg(bts)
							if err != nil {
								return
							}
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
			}
		case "effective_shield":
			var zwfl uint32
			zwfl, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.EffectiveShield == nil && zwfl > 0 {
				z.EffectiveShield = make(map[string][]ShieldSingleInterval, zwfl)
			} else if len(z.EffectiveShield) > 0 {
				for key, _ := range z.EffectiveShield {
					delete(z.EffectiveShield, key)
				}
			}
			for zwfl > 0 {
				var zarl string
				var zctz []ShieldSingleInterval
				zwfl--
				zarl, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				var zdif uint32
				zdif, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(zctz) >= int(zdif) {
					zctz = (zctz)[:zdif]
				} else {
					zctz = make([]ShieldSingleInterval, zdif)
				}
				for zljl := range zctz {
					var zibu uint32
					zibu, bts, err = msgp.ReadMapHeaderBytes(bts)
					if err != nil {
						return
					}
					for zibu > 0 {
						zibu--
						field, bts, err = msgp.ReadMapKeyZC(bts)
						if err != nil {
							return
						}
						switch msgp.UnsafeString(field) {
						case "start":
							zctz[zljl].Start, bts, err = msgp.ReadIntBytes(bts)
							if err != nil {
								return
							}
						case "end":
							zctz[zljl].End, bts, err = msgp.ReadIntBytes(bts)
							if err != nil {
								return
							}
						case "hp":
							zctz[zljl].HP, bts, err = msgp.ReadFloat64Bytes(bts)
							if err != nil {
								return
							}
						default:
							bts, err = msgp.Skip(bts)
							if err != nil {
								return
							}
						}
					}
				}
				z.EffectiveShield[zarl] = zctz
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ShieldResult) Msgsize() (s int) {
	s = 1 + 8 + msgp.ArrayHeaderSize
	for zzrg := range z.Shields {
		s += 1 + 5 + msgp.StringPrefixSize + len(z.Shields[zzrg].Name) + 10 + msgp.ArrayHeaderSize
		for zbmy := range z.Shields[zzrg].Intervals {
			s += z.Shields[zzrg].Intervals[zbmy].Msgsize()
		}
	}
	s += 17 + msgp.MapHeaderSize
	if z.EffectiveShield != nil {
		for zarl, zctz := range z.EffectiveShield {
			_ = zctz
			s += msgp.StringPrefixSize + len(zarl) + msgp.ArrayHeaderSize + (len(zctz) * (14 + msgp.IntSize + msgp.IntSize + msgp.Float64Size))
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ShieldSingleInterval) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zuff uint32
	zuff, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zuff > 0 {
		zuff--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "end":
			z.End, err = dc.ReadInt()
			if err != nil {
				return
			}
		case "hp":
			z.HP, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ShieldSingleInterval) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "start"
	err = en.Append(0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.Start)
	if err != nil {
		return
	}
	// write "end"
	err = en.Append(0xa3, 0x65, 0x6e, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.End)
	if err != nil {
		return
	}
	// write "hp"
	err = en.Append(0xa2, 0x68, 0x70)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.HP)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ShieldSingleInterval) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "start"
	o = append(o, 0x83, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	o = msgp.AppendInt(o, z.Start)
	// string "end"
	o = append(o, 0xa3, 0x65, 0x6e, 0x64)
	o = msgp.AppendInt(o, z.End)
	// string "hp"
	o = append(o, 0xa2, 0x68, 0x70)
	o = msgp.AppendFloat64(o, z.HP)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ShieldSingleInterval) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zmow uint32
	zmow, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zmow > 0 {
		zmow--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "end":
			z.End, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		case "hp":
			z.HP, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ShieldSingleInterval) Msgsize() (s int) {
	s = 1 + 6 + msgp.IntSize + 4 + msgp.IntSize + 3 + msgp.Float64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ShieldStats) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zslz uint32
	zslz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zslz > 0 {
		zslz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "intervals":
			var zoqj uint32
			zoqj, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Intervals) >= int(zoqj) {
				z.Intervals = (z.Intervals)[:zoqj]
			} else {
				z.Intervals = make([]ShieldInterval, zoqj)
			}
			for zdit := range z.Intervals {
				err = z.Intervals[zdit].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ShieldStats) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "name"
	err = en.Append(0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "intervals"
	err = en.Append(0xa9, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Intervals)))
	if err != nil {
		return
	}
	for zdit := range z.Intervals {
		err = z.Intervals[zdit].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ShieldStats) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "name"
	o = append(o, 0x82, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "intervals"
	o = append(o, 0xa9, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Intervals)))
	for zdit := range z.Intervals {
		o, err = z.Intervals[zdit].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ShieldStats) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zmqr uint32
	zmqr, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zmqr > 0 {
		zmqr--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "intervals":
			var ztic uint32
			ztic, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Intervals) >= int(ztic) {
				z.Intervals = (z.Intervals)[:ztic]
			} else {
				z.Intervals = make([]ShieldInterval, ztic)
			}
			for zdit := range z.Intervals {
				bts, err = z.Intervals[zdit].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ShieldStats) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 10 + msgp.ArrayHeaderSize
	for zdit := range z.Intervals {
		s += z.Intervals[zdit].Msgsize()
	}
	return
}
