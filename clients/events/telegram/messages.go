package telegram

const msgHelp = `I can save and keep you pages. Also I can offer you them to read any time you online.

In order to save the page, just send me message link or message itself.

In order to get a random saved page , send me command /rnd.
Caution! After that , this page will be removed from your list!`

const msgHello = "Hi there!👋 \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command 🤔"
	msgNoSavedPages   = "You have no saved pages 🤷‍♂️"
	msgSaved          = "Saved! 👌"
	msgAlreadyExists  = "You have already saved this page in your list 🫡"
)
