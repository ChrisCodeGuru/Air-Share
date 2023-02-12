package detection

import "regexp"

// Determines if text fits regex sensitive data requirements. Lone parameter determines if function should determine sensitive data if it is the whole text passed in or only part of the text
func SensitiveData(text string, lone bool) bool {

	sensitiveDataRegex := [4]string{
		`\b[0-9A-Z]{3}([^ 0-9A-Z]|\s)?[0-9]{4}\b`, // US License plate
		`[0-9]{3}-[0-9]{2}-[0-9]{4}`,              // Social Security number
		`(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})`, // MasterCard, Visa, American Express, Diners Club, Discover, JCB cards
		`[STFG]\d{7}[A-Z]`, // NRIC
	}

	if lone {
		sensitiveDataRegex[0] = `^\b[0-9A-Z]{3}([^ 0-9A-Z]|\s)?[0-9]{4}\b$`                                                                                                                // US License plate
		sensitiveDataRegex[1] = `^[0-9]{3}-[0-9]{2}-[0-9]{4}$`                                                                                                                             // Social Security number
		sensitiveDataRegex[2] = `^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$` // MasterCard, Visa, American Express, Diners Club, Discover, JCB cards
		sensitiveDataRegex[3] = `^[STFG]\d{7}[A-Z]$`                                                                                                                                       // NRIC
	}

	// iterate through sensitive data array
	for _, sensitiveFormat := range sensitiveDataRegex {

		// check for regex match
		matched, _ := regexp.MatchString(sensitiveFormat, text)
		if matched {
			return true
		}
	}

	return false
}

func SensitiveOccurrences(text string) int {

	return len(Sensitive.USLicensePlate.FindAllStringIndex(text, -1)) + len(Sensitive.SocialSecurityNumber.FindAllStringIndex(text, -1)) + len(Sensitive.PaymentCard.FindAllStringIndex(text, -1)) + len(Sensitive.NRIC.FindAllStringIndex(text, -1))

}

type sensitive struct {
	USLicensePlate       *regexp.Regexp
	SocialSecurityNumber *regexp.Regexp
	PaymentCard          *regexp.Regexp
	NRIC                 *regexp.Regexp
}

var Sensitive = sensitive{
	regexp.MustCompile(`\b[0-9A-Z]{3}([^ 0-9A-Z]|\s)?[0-9]{4}\b`),
	regexp.MustCompile(`[0-9]{3}-[0-9]{2}-[0-9]{4}`),
	regexp.MustCompile(`(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})`),
	regexp.MustCompile(`[STFG]\d{7}[A-Z]`),
}
