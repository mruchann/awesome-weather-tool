package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type model struct {
	choices    []string
	cursor     int
	selected   map[int]struct{}
	count      int
	isChoosing bool
	isTyping   bool
	isLoading  bool
	weather    Weather
	latitude   float64
	longitude  float64
	textInput1 textinput.Model
	textInput2 textinput.Model
	URL        string
}

type Data struct {
	Temp       float64 `json:"temp"`
	Feels_like float64 `json:"feels_like"`
	Humidity   float64 `json:"humidity"`
}
type Weather struct {
	Current Data   `json:"current"`
	Hourly  []Data `json:"hourly"`
}

func initialModel() model {
	t1 := textinput.New()
	t2 := textinput.New()
	return model{
		choices:    []string{"Current", "Hourly"},
		selected:   make(map[int]struct{}),
		textInput1: t1,
		textInput2: t2,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.count == 0 {
		m.isChoosing = true

	} else {
		m.isChoosing = false
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected[m.cursor] = struct{}{}
			m.count++
			m.isChoosing = false
			switch m.count {
			case 1: // latitude
				m.isTyping = true
			case 2: // longitude
				m.isTyping = true
			case 3: // result
				m.isTyping = false
				m.isLoading = true
			}
		default:
			if m.isTyping {
				var cmd tea.Cmd
				switch m.count {
				case 1: // latitude
					m.textInput1.Focus()
					m.textInput1, cmd = m.textInput1.Update(msg)
					m.textInput1.Blur()
				case 2: // longitude
					m.textInput2.Focus()
					m.textInput2, cmd = m.textInput2.Update(msg)
					m.textInput2.Blur()
				}
				return m, cmd
			}
		}
	}
	if m.isLoading {
		var err error
		m.latitude, err = strconv.ParseFloat(m.textInput1.Value(), 64)
		if err != nil {
			log.Fatalln(err)
		}
		m.longitude, err = strconv.ParseFloat(m.textInput2.Value(), 64)
		if err != nil {
			log.Fatalln(err)
		}
		m.URL = fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&appid=f90f2252ffbb060ba8d8c3bd7e7e500d&units=metric", m.latitude, m.longitude)
		response, err := http.Get(m.URL)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		text := string(body)
		err = json.Unmarshal([]byte(text), &m.weather)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	if m.isTyping {
		s += "The weather for:\nLatitude?\n"
		if m.count >= 1 {
			s += m.textInput1.View()
		}
		if m.count >= 2 {
			s += "\nLongitude?\n"
			s += m.textInput2.View()
		}

	}

	if m.isLoading {
		switch m.cursor {
		case 0:
			s += fmt.Sprintf("current temp: %.2f, feels like: %.2f, humidity: %.2f\n", m.weather.Current.Temp, m.weather.Current.Feels_like, m.weather.Current.Humidity)
		case 1:
			s += fmt.Sprintf("hourly temp: %.2f, feels like: %.2f, humidity: %.2f\n", m.weather.Hourly[0].Temp, m.weather.Hourly[0].Feels_like, m.weather.Hourly[0].Humidity)
		}

	}

	if m.isChoosing {
		s = "Choose a frequency.\n\n"

		for i, choice := range m.choices {

			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}

			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}
	s += "\nPress q to quit.\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
