package domain

type Route struct {
	Status RouteStatus
	Days   []RouteDay
}

func (r *Route) Validate() error {
	if err := validateRouteStatus(r.Status); err != nil {
		return err
	}
	if err := validateRouteDays(r.Days); err != nil {
		return err
	}

	return nil
}
