package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea" // Import bubbletea TUI package
)

func main() {

	// Define flags for build and run
	runFlag := flag.String("run", "", "Run a Docker container. Provide the image name and optional arguments (e.g., 'my-image:latest')")
	// buildFlag := flag.String("build", "", "Build a Docker image. Provide the Dockerfile path and tag (e.g., 'Dockerfile:my-image:latest')")
	flag.Parse()

	if *runFlag != "" {
		runDockerContainer(*runFlag)
	} else // if *buildFlag != "" {
	// 	buildDockerImage(*buildFlag)
	// } else
	{
		fmt.Println("Usage:")
		fmt.Println("  -run <image:tag>           Run a Docker container")
		// fmt.Println("  -build <Dockerfile:tag>   Build a Docker image")
	}
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices:  []string{"Run"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Navegation
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// What was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

// UI
func (m model) View() string {
	// The header
	s := "What should we do with Docker?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

// Docker Functions
// Docker run
func runDockerContainer(imageTag string) {
	cmd := exec.Command("docker", "run", imageTag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running Docker container: %s\n", err)
	}
	fmt.Printf("Output: %s\n", output)
}

// Docker build
// func buildDockerImage(dockerfileTag string) {
// 	parts := splitDockerfileTag(dockerfileTag)
// 	if len(parts) != 2 {
// 		fmt.Println("Invalid build argument. Use 'Dockerfile:tag'")
// 		return
// 	}
// 	dockerfile := parts[0]
// 	tag := parts[1]

// 	cmd := exec.Command("docker", "build", "-f", dockerfile, "-t", tag, ".")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("Error building Docker image: %s\n", err)
// 	}
// 	fmt.Printf("Output: %s\n", output)
// }

// func splitDockerfileTag(dockerfileTag string) []string {
// 	return []string{
// 		dockerfileTag[:len(dockerfileTag)-len(":latest")],
// 		dockerfileTag[len(dockerfileTag)-len(":latest"):],
// 	}
// }
