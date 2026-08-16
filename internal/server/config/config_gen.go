package config

//type ConfOption[T any] func(*T)
//
//func BuildConfig[T any](opts ...ConfOption[T]) (*T, error) {
//	var cfg T
//	for _, opt := range opts {
//		if err := opt(&cfg); err != nil {
//			return nil, err
//		}
//	}
//	return &cfg, nil
//}
//
//func ApplyDBHost() ConfOption[cfg] {
//	return func(cfg *Config) {
//		var dbName, dbHost, dbUser, dbPass, sslMode string
//		var dbPort int
//
//	}
//}
