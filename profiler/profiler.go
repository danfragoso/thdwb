package profiler

import "time"

type Profiler struct {
	profiles []*Profile
}

type Profile struct {
	name      string
	startTime time.Time
	endTime   time.Time
}

func CreateProfiler() *Profiler {
	var profiles []*Profile
	return &Profiler{profiles: profiles}
}

func (p *Profiler) Start(profileName string) {
	for _, profile := range p.profiles {
		if profile.name == profileName {
			profile.startTime = time.Now()
			return
		}
	}

	p.profiles = append(p.profiles, &Profile{
		name:      profileName,
		startTime: time.Now(),
	})
}

func (p *Profiler) Stop(profileName string) {
	for _, profile := range p.profiles {
		if profile.name == profileName {
			profile.endTime = time.Now()
		}
	}
}

func (p *Profiler) GetAllProfiles() []*Profile {
	return p.profiles
}

func (p *Profiler) GetProfile(profileName string) *Profile {
	var rProfile *Profile

	for _, profile := range p.profiles {
		if profile.name == profileName {
			rProfile = profile
		}
	}

	return rProfile
}

func (p *Profile) GetElapsedTime() time.Duration {
	return p.endTime.Sub(p.startTime)
}

func (p *Profile) GetName() string {
	return p.name
}
