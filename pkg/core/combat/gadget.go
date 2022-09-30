package combat

func (h *Handler) RemoveGadget(i int) {
	if i < 0 || i >= len(h.gadgets) {
		return
	}
	//set to nil for now; we should clean up every so often???
	//TODO: how often do we clean out nil entries? if at all?
	h.gadgets[i] = nil
}

func (h *Handler) AddGadget(t Target) {
	h.gadgets = append(h.gadgets, t)
	t.SetIndex(len(h.gadgets) - 1)
	t.SetKey(h.nextkey())
}

func (h *Handler) ReplaceGadget(key int, t Target) {
	//do nothing if not found
	for i, v := range h.gadgets {
		if v.Key() == key {
			h.gadgets[i] = t
		}
	}
}

func (h *Handler) Gadget(i int) Target {
	return h.gadgets[i]
}

func (h *Handler) GadgetCount() int {
	count := 0
	for _, v := range h.gadgets {
		if v != nil {
			count++
		}
	}

	return count
}
