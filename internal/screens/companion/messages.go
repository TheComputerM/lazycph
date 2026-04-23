package companion

import tea "charm.land/bubbletea/v2"

var serverChan = make(chan Data, 1)

func StartServer() tea.Msg {
	go CreateServer(func(d Data) {
		serverChan <- d
	})

	return <-serverChan
}

func requestServer() tea.Msg {
	return <-serverChan
}
