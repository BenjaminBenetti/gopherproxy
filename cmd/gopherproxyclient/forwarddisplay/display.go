package forwarddisplay

type Display interface {
	// build the display
	Build()
	// start drawing the display
	Start()
}
