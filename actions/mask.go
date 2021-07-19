package actions

/*
type Mask struct {
	Character string `json:"character"`
}

func (m Mask) Validate(input interface{}, values map[string]interface{}) (interface{}, error) {
	regex := regexp.MustCompile(`[^-/\.]`)
	strValue, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected a string")
	}
	return regex.ReplaceAllString(strValue, "*"), nil
}
*/
