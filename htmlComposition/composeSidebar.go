package htmlComposition

import (
	. "PoliSim/componentHelper"
	"fmt"
	"os"
)

var sidebarElements = make([]Node, 0)

func getSidebar() Node {
	return El(DIV, sidebarElements...)
}

func SetupSidebar() {
	sidebarGrouping := ""
	sidebarGroup := make([]Node, 0)

	for _, httpRoute := range LoadingList {

		handler := HandlerList[httpRoute]
		if handler == nil || !handler.HasSidebarButton {
			continue
		}
		node := makeNode(handler)

		if handler.HasSidebarSubMenu && sidebarGrouping == handler.SidebarSubMenuText {
			sidebarGroup = append(sidebarGroup, node)
		} else if handler.HasSidebarSubMenu {
			sidebarGroup = addNewSubMenu(sidebarGroup, sidebarGrouping)
			sidebarGrouping = handler.SidebarSubMenuText
			sidebarGroup = append(sidebarGroup, node)
		} else {
			sidebarElements = append(sidebarElements, node)
		}
	}

	addNewSubMenu(sidebarGroup, sidebarGrouping)
	_, _ = fmt.Fprintf(os.Stdout, "Finished composing sidebar\n")
}

func addNewSubMenu(group []Node, groupName string) []Node {
	return make([]Node, 0)
}

func makeNode(handler *HttpHandling) Node {
	return nil
}
