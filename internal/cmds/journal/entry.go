package journal

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nirenjan.org/overlord/internal/config"
	"nirenjan.org/overlord/internal/terminal"
	"nirenjan.org/overlord/internal/util"
)

// Entry holds a single entry on disk
type Entry struct {
	ID    string    `json:"-"`
	Title string    `json:"title"`
	Body  string    `json:"text"`
	Date  time.Time `json:"timestamp"`
	Tags  []string  `json:"tags,omitempty"`
	Path  string    `json:"-"`
}

func entryFromFile(file string) (Entry, error) {
	f, err1 := os.Open(file)
	entry := Entry{Path: file}
	if err1 != nil {
		return entry, err1
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		text := scanner.Text()
		switch line {
		case 0:
			// Date
			var err2 error
			entry.Date, err2 = time.Parse(time.RFC1123Z, text)
			if err2 != nil {
				return entry, err2
			}

		case 1:
			// Tags
			entry.Tags = strings.Fields(text)

		case 2:
			// Title
			entry.Title = text

		default:
			// Body
			entry.Body += text + "\n"
		}

		line++
	}
	if err3 := scanner.Err(); err3 != nil {
		return entry, err3
	}

	entry.UpdateID()
	return entry, nil
}

func newEntry(tags []string) (Entry, error) {
	entry := Entry{
		Tags: tags,
		Date: time.Now(),
	}

	if err := entry.UpdatePath(); err != nil {
		return Entry{}, err
	}

	entry.Write()
	entry.UpdateID()

	return entry, nil
}

func (e *Entry) UpdatePath() error {
	dpath, err := config.ModuleDir("journal", e.Date.Format("2006"))
	if err != nil {
		return err
	}
	e.Path = filepath.Join(dpath, e.Date.Format("0102-150405.entry"))
	return nil
}

func (e *Entry) Write() error {
	efile, err := os.Create(e.Path)
	if err != nil {
		return err
	}
	defer efile.Close()

	// Write the date to the file
	date_line := []byte(e.Date.Format(time.RFC1123Z + "\n"))
	efile.Write(date_line)

	// Write the tags to the file
	tags_line := []byte(strings.Join(e.Tags, " ") + "\n")
	efile.Write(tags_line)

	// Write the title to the file
	title_line := []byte(e.Title + "\n")
	efile.Write(title_line)

	// Write the rest of the body to the file
	efile.Write([]byte(e.Body))

	return nil
}

func (e *Entry) UpdateID() string {
	// ID is the Unix date in hex, followed by the SHA256 checksum of
	// the date and title
	hash_inp := fmt.Sprintf("%v %v", e.Date.Format(time.RFC3339), e.Title)
	hash := sha256.Sum256([]byte(hash_inp))

	id := fmt.Sprintf("%08x-%x", e.Date.Unix(), hash[:5])

	e.ID = id
	return id
}

func (entry *Entry) Display() {
	out := os.Stdout

	out.WriteString(terminal.Foreground(terminal.Yellow))
	out.WriteString(entry.Date.Format(time.RFC1123))

	// Stardate calculation
	stardate := entry.Date.Unix()/864 + 4058750
	stardate_str := fmt.Sprintf(" (Stardate %v.%v)\n",
		stardate/100, stardate%100)
	out.WriteString(terminal.Foreground(terminal.Red))
	out.WriteString(stardate_str)

	// Title
	out.WriteString(terminal.Reset())
	out.WriteString(entry.Title + "\n")
	out.WriteString(strings.Repeat("=", len(entry.Title)) + "\n")

	// Body
	out.WriteString(entry.Body)
	out.WriteString("\n")

	// Tags
	if len(entry.Tags) > 0 {
		out.WriteString("Tags:\t")
		out.WriteString(terminal.Foreground(terminal.Cyan))

		for _, t := range entry.Tags {
			out.WriteString(t + " ")
		}
		out.WriteString("\n" + terminal.Reset())
	}

	out.WriteString(terminal.HorizontalLine())
	out.WriteString("\n")
}

func (entry *Entry) Edit() error {
	// Create a temporary file, and call the editor to edit the file
	tempfile, err := ioutil.TempFile("", "journal*")
	if err != nil {
		return err
	}
	tempname := tempfile.Name()
	defer os.Remove(tempname)

	if entry.Title == "" {
		tempfile.WriteString(`
# Enter your journal entry here. Lines beginning with # are deleted
# from the journal. The first line of the message is the title.
`)
	} else {
		tempfile.WriteString(entry.Title + "\n")
		tempfile.WriteString(entry.Body)
	}
	tempfile.Close()

	// Call the editor
	err1 := util.Editor(tempname)
	if err1 != nil {
		return err1
	}

	// Read the file contents
	file, err2 := os.Open(tempname)
	if err2 != nil {
		return err2
	}
	defer file.Close()

	content, err3 := ioutil.ReadAll(file)
	if err3 != nil {
		return err3
	}

	var title string
	var body []string
	for i, line := range strings.Split(string(content), "\n") {
		if i == 0 {
			title = line
		} else if len(line) == 0 || line[0] != '#' {
			body = append(body, line)
		}
	}

	if len(body) == 0 {
		return errors.New("No body in journal entry")
	}

	entry.Title = title
	entry.Body = strings.Join(body, "\n")
	entry.Write()
	entry.UpdateID()

	return nil
}
