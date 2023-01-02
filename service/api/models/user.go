package models

// FullUser - Provides extended information about someone with a WASA Photo account, like total photos count etc.
type FullUser struct {
	BaseUser

	// Total photos count
	TotalPhotos *int `json:"totalPhotos,omitempty"`

	// Total followers count
	TotalFollowers *int `json:"totalFollowers,omitempty"`

	// Total followings count
	TotalFollowings *int `json:"totalFollowings,omitempty"`

	// If a user is banned for the authenticated user
	BannedForUser *bool `json:"bannedForUser,omitempty"`
}
