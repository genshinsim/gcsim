package shield

import (
	"maps"
	"math"
	"slices"
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
	// on the same frame, ends come before starts
	return p[i].pos < p[j].pos || (p[i].pos == p[j].pos && p[i].end && !p[j].end)
}

func computeEffective(shields map[string][]stats.ShieldInterval) map[string][]stats.ShieldSingleInterval {
	var n int
	for _, v := range shields {
		n += len(v)
	}

	// populate and sort endpoints. Names are sorted and the sort is stable so that endpoints on
	// the same frame keep an order that doesn't depend on map iteration
	endpoints := make([]endpoint, 0, n*2)
	for _, name := range slices.Sorted(maps.Keys(shields)) {
		for i := range shields[name] {
			start := endpoint{pos: shields[name][i].Start, end: false, interval: &shields[name][i]}
			end := endpoint{pos: shields[name][i].End, end: true, interval: &shields[name][i]}
			endpoints = append(endpoints, start, end)
		}
	}
	sort.Stable(byPosition(endpoints))

	out := make(map[string][]stats.ShieldSingleInterval)
	for _, e := range elements {
		out[e.String()] = make([]stats.ShieldSingleInterval, 0, n)
	}
	out[normalized] = make([]stats.ShieldSingleInterval, 0, n)

	current := make(map[attributes.Element]*stats.ShieldInterval)
	// active intervals in the order they started
	var active []*stats.ShieldInterval
	handleEnd := func(endpoint endpoint) {
		// if other shields are active, need to elect greatest as new effective
		var normalizedHP float64
		normalizedEnd := math.MaxInt
		for _, e := range elements {
			if current[e] == endpoint.interval {
				// a shield that ends on this frame can't take over; on a tie, the one that started
				// first wins
				var best *stats.ShieldInterval
				for _, k := range active {
					if k.End > endpoint.pos && (best == nil || k.HP[e.String()] > best.HP[e.String()]) {
						best = k
					}
				}
				if best != nil {
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
			if i := slices.Index(active, endpoint.interval); i >= 0 {
				active = slices.Delete(active, i, i+1)
			}
			if len(active) == 0 {
				continue
			}
			handleEnd(endpoint)
			continue
		}

		// start endpoint, add this interval to active set
		active = append(active, endpoint.interval)

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
