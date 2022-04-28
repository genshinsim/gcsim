package player

func (m *MasterChar) Tag(key string) int {
	return m.Tags[key]
}

func (m *MasterChar) AddTag(key string, val int) {
	m.Tags[key] = val
}

func (m *MasterChar) RemoveTag(key string) {
	delete(m.Tags, key)
}
