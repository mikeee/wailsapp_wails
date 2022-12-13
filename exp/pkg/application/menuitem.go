package application

import (
	"sync"
	"sync/atomic"
)

type menuItemType int

const (
	text menuItemType = iota
	separator
	checkbox
	radio
	submenu
)

var menuItemID uintptr
var menuItemMap = make(map[uint]*MenuItem)
var menuItemMapLock sync.Mutex

func addToMenuItemMap(menuItem *MenuItem) {
	menuItemMapLock.Lock()
	menuItemMap[menuItem.id] = menuItem
	menuItemMapLock.Unlock()
}

func getMenuItemByID(id uint) *MenuItem {
	menuItemMapLock.Lock()
	defer menuItemMapLock.Unlock()
	return menuItemMap[id]
}

type menuItemImpl interface {
	setTooltip(s string)
	setLabel(s string)
	setDisabled(disabled bool)
	setChecked(checked bool)
}

type MenuItem struct {
	id       uint
	label    string
	tooltip  string
	disabled bool
	checked  bool
	submenu  *Menu
	callback func(*Context)
	itemType menuItemType

	impl menuItemImpl
}

func newMenuItem(label string) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		itemType: text,
	}
	addToMenuItemMap(result)
	return result
}

func newMenuItemSeperator() *MenuItem {
	result := &MenuItem{
		itemType: separator,
	}
	return result
}

func newMenuItemCheckbox(label string, checked bool) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		checked:  checked,
		itemType: checkbox,
	}
	addToMenuItemMap(result)
	return result
}

func newSubMenuItem(label string) *MenuItem {
	result := &MenuItem{
		id:       uint(atomic.AddUintptr(&menuItemID, 1)),
		label:    label,
		itemType: submenu,
		submenu:  &Menu{},
	}
	addToMenuItemMap(result)
	return result
}

func (m *MenuItem) handleClick() {
	var ctx = newContext().withClickedMenuItem(m)
	if m.itemType == checkbox {
		m.checked = !m.checked
		ctx.withChecked(m.checked)
		m.impl.setChecked(m.checked)
	}
	if m.callback != nil {
		go m.callback(ctx)
	}
}

func (m *MenuItem) SetTooltip(s string) *MenuItem {
	m.tooltip = s
	if m.impl != nil {
		m.impl.setTooltip(s)
	}
	return m
}

func (m *MenuItem) SetLabel(s string) *MenuItem {
	m.label = s
	if m.impl != nil {
		m.impl.setLabel(s)
	}
	return m
}

func (m *MenuItem) SetEnabled(enabled bool) *MenuItem {
	m.disabled = !enabled
	if m.impl != nil {
		m.impl.setDisabled(m.disabled)
	}
	return m
}

func (m *MenuItem) OnClick(f func(*Context)) *MenuItem {
	m.callback = f
	return m
}

func (m *MenuItem) Label() string {
	return m.label
}

func (m *MenuItem) Tooltip() string {
	return m.tooltip
}

func (m *MenuItem) Enabled() bool {
	return !m.disabled
}
