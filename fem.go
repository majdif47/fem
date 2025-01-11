package main

import (
	"fmt"
	"os"
	"time"
	"regexp"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)



var docStyle = lipgloss.NewStyle().Margin(2, 2)

type item struct {
	title, desc string
}


func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.desc}


type model struct {
	list list.Model
  message string
  showMsg bool
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  // m.list.SetFilteringEnabled(false)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		} else if msg.String() == "enter" {
			  if selectedItem, ok := m.list.SelectedItem().(item); ok {
				  selected := selectedItem.title
				  err := copyToClipboard(selected)
            if err != nil {
					    m.message = "Failed to copy!"
              m.showMsg = true
              return m, clearMsg()
				    }else{
              msg := fmt.Sprintf("\tCopied %s", selected)
              m.message = msg
              m.showMsg = true
				      return m, clearMsg()
            }
			  }
		  }
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	
  case clearMessageMsg:
    m.message = ""
    m.showMsg = false
  }



	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

type clearMessageMsg struct{}

func clearMsg() tea.Cmd {
  return tea.Tick(time.Second*3, func(_ time.Time) tea.Msg {
    return clearMessageMsg{}
  })
}


func (m model) View() string {
  m.list.SetShowHelp(false)
  m.list.SetShowPagination(false)
  var message string
  if m.showMsg {
    message = fmt.Sprintf("\n\n%s\n", m.message)
  }
	return docStyle.Render(m.list.View() + message)
}

func copyToClipboard(input string) error {
  err := clipboard.WriteAll(input)
  if err != nil {
    return err
  }
  return nil
}

func main() {
  var items []list.Item
  if len(os.Args) <= 1 {
    fmt.Println("\tPlease Enter a search phase! i.e: fem animals")
    os.Exit(0)
  }
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

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
