package module_struct
import (
	"github.com/charmbracelet/bubbles/spinner"
)
type errMsg error

type model_loading struct {
	spinner  spinner.Model
	quitting bool
	err      error
	Msg      string
}