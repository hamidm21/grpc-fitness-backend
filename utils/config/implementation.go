package config

import (
	"strings"
)

func (m *Manager) GetDefaultString(key, def string) string {
	m.SetDefault(key, def)
	return m.GetString(key)
}

func (m *Manager) GetDefaultInt(key string, def int) int {
	m.SetDefault(key, def)
	return m.GetInt(key)
}

func (m *Manager) GetStringArray(key string) []string {
	strValue := m.GetString(key)
	strArr := make([]string, 0, 10)
	for _, str := range strings.Split(strings.TrimSpace(strValue), ",") {
		if len(str) > 0 {
			strArr = append(strArr, str)
		}
	}
	return strArr
}
