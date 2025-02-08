package core

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 20
const defaultWidth = 20

var (
	titleStyle = lipgloss.
			NewStyle().
			Width(defaultWidth).
			Align(lipgloss.Center).
			Bold(true).
			Border(lipgloss.ThickBorder(), false, false, true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	warningTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
	boldTextStyle = lipgloss.NewStyle().Bold(true)
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))

	Choice ChoiceType = ChoiceType{}
)

type ChoiceType struct {
	Target string
	Type   string
}

type listKeyMap struct {
	editSelection   key.Binding
	deleteSelection key.Binding
	openSelection   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		editSelection: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit ash"),
		),
		deleteSelection: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete ash"),
		),
		openSelection: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open ash"),
		),
	}
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("->  " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list       list.Model
	choice     string
	quitting   bool
	deleting   bool
	configPath string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			Quit = true
			return m, tea.Quit

		case "enter":
			if m.deleting {
				break
			}
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = strings.TrimSpace(string(i))
				Choice.Target = m.choice
				Choice.Type = "open"
			}
			return m, tea.Quit

		case "e":
			if m.deleting {
				break
			}
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = strings.TrimSpace(string(i))
				Choice.Target = m.choice
				Choice.Type = "edit"
			}
			return m, tea.Quit

		case "d":
			if m.deleting {
				break
			}
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.deleting = true
				m.choice = strings.TrimSpace(string(i))
			}

		case "y":
			if !m.deleting {
				break
			}
			m.deleting = false
			idx := m.list.Index()
			Delete(m.configPath, m.choice)
			m.list.RemoveItem(idx)
			m.choice = ""

		case "n":
			if !m.deleting {
				break
			}
			m.deleting = false
			m.choice = ""
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		//TODO: maybe add a nice quitting message
		return ""
	}
	if m.deleting {
		return fmt.Sprintf(
			"\n\n  %s %s %s %s %s\n",
			boldTextStyle.Render("Do you wish to"),
			warningTextStyle.Render("delete"),
			boldTextStyle.Render("the item"),
			keywordStyle.Render(m.choice),
			boldTextStyle.Render("? y/n"),
		)
	}
	if m.choice != "" {
		return ""
	}
	return "\n" + m.list.View()
}

func newList(items []list.Item, configPath string) model {
	listKeys := newListKeyMap()

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Ashes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.openSelection,
			listKeys.editSelection,
			listKeys.deleteSelection,
		}
	}

	return model{
		list:       l,
		configPath: configPath,
	}
}
