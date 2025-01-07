package main

import (
	"fmt"
	"os"
  // "os/exec"
  // "strings"
  // "runtime"
	"regexp"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
  "github.com/atotto/clipboard"
)



var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }


type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  m.list.SetFilteringEnabled(false)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		} else if msg.String() == "enter" {
			if selectedItem, ok := m.list.SelectedItem().(item); ok {
				selected := selectedItem.title
				err := copyToClipboard(selected)
				if err != nil {
					// Instead of exiting, print an error message and let the user know
					fmt.Println("Error copying to clipboard:", err)
					// Optionally, return a new model state or a command here
					return m, nil
				}
				// Optionally, notify the user the emoji was copied successfully
				fmt.Println("Copied to clipboard:", selected)
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func copyToClipboard(input string) error {
  fmt.Println(input)
  err := clipboard.WriteAll(input)
  if err != nil {
    return err
  }
  return nil
}

func main() {
  var items []list.Item
  for key, value := range emojies {
    for _,arg := range os.Args[1:] {
      pattern := fmt.Sprintf(`(?i).*%s*`,regexp.QuoteMeta(arg))
      regex, err := regexp.Compile(pattern)
    	if err != nil {
		    fmt.Println("Invalid regex:", err)
		    return
	    }else if regex.MatchString(key) {
        items = append(items, item{title: value, desc: key})
      }
    }
  }


	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
