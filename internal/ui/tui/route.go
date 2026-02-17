package tui

type route string

const (
	routeDashboard route = "Dashboard"
	routeBooks     route = "Books"
	routeMembers   route = "Members"
	routeLoans     route = "Loans"
	routeReports   route = "Reports"
	routeSettings  route = "Settings"
)

var allRoutes = []route{
	routeDashboard,
	routeBooks,
	routeMembers,
	routeLoans,
	routeReports,
	routeSettings,
}

func nextRoute(current route) route {
	idx := routeIndex(current)
	if idx == -1 {
		return routeDashboard
	}

	idx = (idx + 1) % len(allRoutes)
	return allRoutes[idx]
}

func prevRoute(current route) route {
	idx := routeIndex(current)
	if idx == -1 {
		return routeDashboard
	}

	idx--
	if idx < 0 {
		idx = len(allRoutes) - 1
	}

	return allRoutes[idx]
}

func routeIndex(current route) int {
	for i, r := range allRoutes {
		if r == current {
			return i
		}
	}

	return -1
}
