package shield

import (
	"math"
	"sort"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/stats"
)

type endpoint struct {
	pos      int
	end      bool
	interval *stats.ShieldInterval
}
type byPosition []endpoint

func (p byPosition) Len() int      { return len(p) }
func (p byPosition) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p byPosition) Less(i, j int) bool {
	return p[i].pos < p[j].pos || (p[i].pos == p[j].pos && p[i].end)
}

func computeEffective(shields map[string][]stats.ShieldInterval) map[string][]stats.ShieldSingleInterval {
	var n int
	for _, v := range shields {
		n += len(v)
	}

	// populate and sort endpoints
	endpoints := make([]endpoint, 0, n*2)
	for name := range shields {
		for i := range shields[name] {
			start := endpoint{pos: shields[name][i].Start, end: false, interval: &shields[name][i]}
			end := endpoint{pos: shields[name][i].End, end: true, interval: &shields[name][i]}
			endpoints = append(endpoints, start, end)
		}
	}
	sort.Sort(byPosition(endpoints))

	out := make(map[string][]stats.ShieldSingleInterval)
	for _, e := range elements {
		out[e.String()] = make([]stats.ShieldSingleInterval, 0, n)
	}
	out[normalized] = make([]stats.ShieldSingleInterval, 0, n)

	current := make(map[attributes.Element]*stats.ShieldInterval)
	active := make(map[*stats.ShieldInterval]bool)
	handleEnd := func(endpoint endpoint) {
		// if other shields are active, need to elect greatest as new effective
		var normalizedHP float64
		normalizedEnd := math.MaxInt
		for _, e := range elements {
			if current[e] == endpoint.interval {
				var best *stats.ShieldInterval
				for k := range active {
					if best == nil || k.HP[e.String()] > best.HP[e.String()] {
						best = k
					}
				}
				// only add if this new interval is at least 1 frame wide
				if best.End > endpoint.pos {
					current[e] = best
					out[e.String()] = append(out[e.String()], stats.ShieldSingleInterval{
						Start: endpoint.pos,
						End:   best.End,
						HP:    best.HP[e.String()],
					})
					normalizedEnd = min(normalizedEnd, best.End)
				}
			}
			normalizedHP += out[e.String()][len(out[e.String()])-1].HP
		}
		// at least one of the effective elementals got a new interval, so recompute normalized
		if normalizedEnd >= math.MaxInt {
			return
		}
		out[normalized] = append(out[normalized], stats.ShieldSingleInterval{
			Start: endpoint.pos,
			End:   normalizedEnd,
			HP:    normalizedHP / float64(len(elements)),
		})
	}

	for _, endpoint := range endpoints {
		if endpoint.end {
			delete(active, endpoint.interval)
			if len(active) == 0 {
				continue
			}
			handleEnd(endpoint)
			continue
		}

		// start endpoint, add this interval to active set
		active[endpoint.interval] = true

		// only 1 active shield == effective shield
		if len(active) == 1 {
			for _, e := range elements {
				out[e.String()] = append(out[e.String()], stats.ShieldSingleInterval{
					Start: endpoint.pos,
					End:   endpoint.interval.End,
					HP:    endpoint.interval.HP[e.String()],
				})
				current[e] = endpoint.interval
			}

			out[normalized] = append(out[normalized], stats.ShieldSingleInterval{
				Start: endpoint.pos,
				End:   endpoint.interval.End,
				HP:    endpoint.interval.HP[normalized],
			})
			continue
		}

		// for each element if new interval > current, end current early and add new to effective
		var normalizedHP float64
		modified := false
		for _, e := range elements {
			currentIdx := len(out[e.String()]) - 1
			if endpoint.interval.HP[e.String()] > out[e.String()][currentIdx].HP {
				newShieldInterval := stats.ShieldSingleInterval{
					Start: endpoint.pos,
					End:   endpoint.interval.End,
					HP:    endpoint.interval.HP[e.String()],
				}
				current[e] = endpoint.interval
				modified = true

				if out[e.String()][currentIdx].Start == endpoint.pos {
					// this shield is a replacement rather than append
					out[e.String()][currentIdx] = newShieldInterval
				} else {
					out[e.String()][currentIdx].End = endpoint.pos
					out[e.String()] = append(out[e.String()], newShieldInterval)
				}
			}
			normalizedHP += out[e.String()][len(out[e.String()])-1].HP
		}

		if modified {
			newShieldInterval := stats.ShieldSingleInterval{
				Start: endpoint.pos,
				End:   endpoint.interval.End,
				HP:    normalizedHP / float64(len(elements)),
			}
			prevIndex := len(out[normalized]) - 1
			if out[normalized][prevIndex].Start == endpoint.pos {
				out[normalized][prevIndex] = newShieldInterval
			} else {
				out[normalized][prevIndex].End = endpoint.pos
				out[normalized] = append(out[normalized], newShieldInterval)
			}
		}
	}
	return out
}
