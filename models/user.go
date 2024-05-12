package models

type UserProfile struct {
	UserId    string
	Name      string
	Nickname  string
	Picture   string
	Iat       float64
	Exp       float64
	UpdatedAt string
	Admin     bool
}

func GetUserProfile(profile map[string]interface{}) UserProfile {
	return UserProfile{
		UserId:    profile["sub"].(string),
		Name:      profile["name"].(string),
		Picture:   profile["picture"].(string),
		Nickname:  profile["nickname"].(string),
		Iat:       profile["iat"].(float64),
		Exp:       profile["exp"].(float64),
		UpdatedAt: profile["updated_at"].(string),
	}
}
