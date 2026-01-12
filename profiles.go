package wise

import (
	"context"
	"fmt"
)

// ProfilesService handles profile-related API calls.
type ProfilesService struct {
	client *Client
}

// Profile represents a Wise profile (personal or business).
type Profile struct {
	ID      int64       `json:"id"`
	Type    ProfileType `json:"type"`
	Details interface{} `json:"details"` // PersonalProfile or BusinessProfile
}

// PersonalProfile represents personal profile details.
type PersonalProfile struct {
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	DateOfBirth      string    `json:"dateOfBirth"` // YYYY-MM-DD
	PhoneNumber      string    `json:"phoneNumber,omitempty"`
	Avatar           string    `json:"avatar,omitempty"`
	Occupation       string    `json:"occupation,omitempty"`
	OccupationFormat string    `json:"occupations,omitempty"`
	PrimaryAddress   *Address  `json:"primaryAddress,omitempty"`
}

// BusinessProfile represents business profile details.
type BusinessProfile struct {
	Name                    string   `json:"name"`
	RegistrationNumber      string   `json:"registrationNumber,omitempty"`
	ACN                     string   `json:"acn,omitempty"`
	ABN                     string   `json:"abn,omitempty"`
	ARBN                    string   `json:"arbn,omitempty"`
	CompanyType             string   `json:"companyType,omitempty"`
	CompanyRole             string   `json:"companyRole,omitempty"`
	DescriptionOfBusiness   string   `json:"descriptionOfBusiness,omitempty"`
	PrimaryAddress          *Address `json:"primaryAddress,omitempty"`
	Webpage                 string   `json:"webpage,omitempty"`
}

// CreatePersonalProfileRequest represents the request to create a personal profile.
type CreatePersonalProfileRequest struct {
	Type    ProfileType      `json:"type"`
	Details *PersonalProfile `json:"details"`
}

// CreateBusinessProfileRequest represents the request to create a business profile.
type CreateBusinessProfileRequest struct {
	Type    ProfileType      `json:"type"`
	Details *BusinessProfile `json:"details"`
}

// List returns all profiles belonging to the authenticated user.
// GET /v1/profiles
func (s *ProfilesService) List(ctx context.Context) ([]Profile, error) {
	var profiles []Profile
	err := s.client.Get(ctx, "/v1/profiles", nil, &profiles)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

// Get returns a profile by ID.
// GET /v1/profiles/{profileId}
func (s *ProfilesService) Get(ctx context.Context, profileID int64) (*Profile, error) {
	var profile Profile
	path := fmt.Sprintf("/v1/profiles/%d", profileID)
	err := s.client.Get(ctx, path, nil, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// CreatePersonal creates a new personal profile.
// POST /v1/profiles
func (s *ProfilesService) CreatePersonal(ctx context.Context, details *PersonalProfile) (*Profile, error) {
	req := CreatePersonalProfileRequest{
		Type:    ProfileTypePersonal,
		Details: details,
	}
	var profile Profile
	err := s.client.Post(ctx, "/v1/profiles", req, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// CreateBusiness creates a new business profile.
// POST /v1/profiles
func (s *ProfilesService) CreateBusiness(ctx context.Context, details *BusinessProfile) (*Profile, error) {
	req := CreateBusinessProfileRequest{
		Type:    ProfileTypeBusiness,
		Details: details,
	}
	var profile Profile
	err := s.client.Post(ctx, "/v1/profiles", req, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
