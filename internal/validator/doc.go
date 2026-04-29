// Package validator implements rule-based validation for environment variable sets
// used within envchain deployment contexts.
//
// A Rule can specify:
//   - Required: whether the variable must be present
//   - Pattern:  a regular expression the value must match
//   - AllowedValues: an explicit set of permitted values
//
// Example usage:
//
//	v := validator.New(map[string]validator.Rule{
//		"DATABASE_URL": {Required: true},
//		"PORT":         {Pattern: `^\d+$`},
//		"ENV":          {AllowedValues: []string{"staging", "production"}},
//	})
//
//	if errs := v.Validate(env); errs != nil {
//		for _, e := range errs {
//			fmt.Println(e)
//		}
//	}
package validator
