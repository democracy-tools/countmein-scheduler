package ds

func GetKaplanDemonstration(dsc Client) (*Demonstration, error) {

	var demonstration Demonstration
	err := dsc.Get(KindDemonstration, DemonstrationKaplan, &demonstration)
	if err != nil {
		return nil, err
	}

	return &demonstration, nil
}

func GetVolunteers(dsc Client, demonstration string) ([]Volunteer, error) {

	var volunteers []Volunteer
	err := dsc.GetFilter(KindVolunteer,
		[]FilterField{{Name: "demonstration_id", Operator: "=", Value: demonstration}},
		&volunteers)
	if err != nil {
		return nil, err
	}

	return volunteers, nil
}
