package apis

// SetInitDefaults sets defaults on the configuration used by init
func SetInitDefaults(config *InitConfiguration) {
}

// SetInitDynamicDefaults sets defaults derived at runtime
func SetInitDynamicDefaults(config *InitConfiguration) error {
	return nil
}

// SetJoinDefaults sets defaults on the configuration used by join
func SetJoinDefaults(config *JoinConfiguration) {
}

// SetJoinDynamicDefaults sets defaults derived at runtime
func SetJoinDynamicDefaults(config *JoinConfiguration) error {
	return nil
}
