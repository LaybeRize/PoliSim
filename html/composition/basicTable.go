package composition

import (
	. "PoliSim/html/builder"
	"strconv"
)

type Position int

const (
	StartPos Position = iota
	MiddlePos
	EndPos
)

func getStandardTable(id string, children ...Node) Node {
	return TABLE(ID(id), CLASS("table-auto mt-4 w-[800px]"),
		Group(children...))
}

func getTableHeader(p Position, sortPos int, text string) Node {
	var hScript Node = nil
	if sortPos >= 0 {
		hScript = HYPERSCRIPT("on click call sortTable(" + strconv.Itoa(sortPos) + ")")
	}
	switch p {
	case StartPos:
		return TH(CLASS("p-2 border-r-2 border-gray-600"), If(hScript != nil, STYLE("cursor: pointer;")),
			hScript, Text(text))
	case MiddlePos:
		return TH(CLASS("p-2 border-r-2 border-gray-600"), If(hScript != nil, STYLE("cursor: pointer;")),
			hScript, Text(text))
	case EndPos:
		return TH(CLASS("p-2"), If(hScript != nil, STYLE("cursor: pointer;")), hScript,
			Text(text))
	}
	return nil
}
func getTableElement(p Position, rowSpan int, node Node) Node {
	switch p {
	case StartPos:
		return TD(ROWSPAN(strconv.Itoa(rowSpan)), CLASS("p-2 border-r-2 border-t-2 border-gray-600"), node)
	case MiddlePos:
		return TD(ROWSPAN(strconv.Itoa(rowSpan)), CLASS("p-2 border-r-2 border-t-2 border-gray-600"), node)
	case EndPos:
		return TD(ROWSPAN(strconv.Itoa(rowSpan)), CLASS("p-2 border-t-2 border-gray-600"), node)
	}
	return nil
}

var tableNode = Raw(`<script>
        function sortTable(n) {
            var table, rows, switching, i, x, y, shouldSwitch, dir, switchcount = 0;
            table = document.getElementById("sortTable");
            switching = true;
            //Set the sorting direction to ascending:
            dir = "asc";
            /*Make a loop that will continue until
            no switching has been done:*/
            while (switching) {
                //start by saying: no switching is done:
                switching = false;
                rows = table.rows;
                /*Loop through all table rows (except the
                first, which contains table headers):*/
                for (i = 1; i < (rows.length - 1); i++) {
                    //start by saying there should be no switching:
                    shouldSwitch = false;
                    /*Get the two elements you want to compare,
                    one from current row and one from the next:*/
                    x = rows[i].getElementsByTagName("TD")[n];
                    y = rows[i + 1].getElementsByTagName("TD")[n];
                    /*check if the two rows should switch place,
                    based on the direction, asc or desc:*/
                    if (dir == "asc") {
                        if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase()) {
                            //if so, mark as a switch and break the loop:
                            shouldSwitch= true;
                            break;
                        }
                    } else if (dir == "desc") {
                        if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase()) {
                            //if so, mark as a switch and break the loop:
                            shouldSwitch = true;
                            break;
                        }
                    }
                }
                if (shouldSwitch) {
                    /*If a switch has been marked, make the switch
                    and mark that a switch has been done:*/
                    rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
                    switching = true;
                    //Each time a switch is done, increase this count by 1:
                    switchcount ++;
                } else {
                    /*If no switching has been done AND the direction is "asc",
                    set the direction to "desc" and run the while loop again.*/
                    if (switchcount == 0 && dir == "asc") {
                        dir = "desc";
                        switching = true;
                    }
                }
            }
        }
    </script>`)
