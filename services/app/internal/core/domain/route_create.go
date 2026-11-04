package domain

type RouteCreate struct {
	Status RouteStatus
	Days   []RouteDay
}

func (r *RouteCreate) Validate() error {
	if err := validateRouteStatus(r.Status); err != nil {
		return err
	}
	if err := validateRouteDays(r.Days); err != nil {
		return err
	}

	return nil
}
