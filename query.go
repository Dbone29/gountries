package gountries

import (
	"strings"
	"sync"
)

// Query holds a reference to the QueryHolder struct
var queryInitOnce sync.Once
var queryInstance *Query

// Query contains queries for countries, cities, etc.
type Query struct {
	Countries           map[string]Country
	NameToAlpha2        map[string]string
	Alpha3ToAlpha2      map[string]string
	NativeNameToAlpha2  map[string]string
	CallingCodeToAlpha2 map[string]string
	CurrencyToAlpha2    map[string][]Country
}

// FindCountryByName finds a country by given name
func (q *Query) FindCountryByName(name string) (result Country, err error) {
	lowerName := strings.ToLower(name)
	alpha2, exists := q.NameToAlpha2[lowerName]
	if !exists {
		return Country{}, makeError("Could not find country with name", name)
	}
	return q.Countries[alpha2], nil
}

// FindCountryByNativeName
func (q *Query) FindCountryByNativeName(name string) (result Country, err error) {
	lowerName := strings.ToLower(name)
	alpha2, exists := q.NativeNameToAlpha2[lowerName]
	if !exists {
		return Country{}, makeError("Could not find country with native name", name)
	}
	return q.Countries[alpha2], nil
}

// FindCountryByAlpha fincs a country by given code
func (q *Query) FindCountryByAlpha(code string) (result Country, err error) {
	codeU := strings.ToUpper(code)
	switch {
	case len(code) == 2:
		country, exists := q.Countries[codeU]
		if !exists {
			return Country{}, makeError("Could not find country with code %s", code)
		}
		return country, nil
	case len(code) == 3:
		alpha2, exists := q.Alpha3ToAlpha2[codeU]
		if !exists {
			return Country{}, makeError("Could not find country with code", code)
		}
		return q.Countries[alpha2], nil
	default:
		return Country{}, makeError("Invalid code format", code)
	}
}

func (q *Query) FindCountryByCallingCode(callingCode string) (result Country, err error) {
	alpha2, exists := q.CallingCodeToAlpha2[callingCode]
	if !exists {
		return Country{}, makeError("Could not find country with callingCode", callingCode)
	}
	return q.Countries[alpha2], nil
}

// FindCountriesByCurrency finds a Country based on the given struct data
func (q *Query) FindCountriesByCurrency(currency string) (results []Country, err error) {
	if len(currency) != 3 {
		return nil, makeError("Invalid currency format", currency)
	}

	alpha2, exists := q.CurrencyToAlpha2[currency]
	if !exists {
		return []Country{}, makeError("Could not find countries with currency name", currency)
	}

	return alpha2, nil
}

// FindAllCountries returns a list of all countries
func (q *Query) FindAllCountries() (countries map[string]Country) {
	return q.Countries
}

// FindCountries finds a Country based on the given struct data
func (q Query) FindCountries(c Country) (countries []Country) {

	for _, country := range q.Countries {

		// Name
		//

		if c.Name.Common != "" && strings.EqualFold(c.Name.Common, country.Name.Common) {
			continue
		}

		// Alpha
		//

		if c.Alpha2 != "" && c.Alpha2 != country.Alpha2 {
			continue
		}

		if c.Alpha3 != "" && c.Alpha3 != country.Alpha3 {
			continue
		}

		// Geo
		//

		if c.Geo.Continent != "" && !strings.EqualFold(c.Geo.Continent, country.Geo.Continent) {
			continue
		}

		if c.Geo.Region != "" && !strings.EqualFold(c.Geo.Region, country.Geo.Region) {
			continue
		}

		if c.Geo.SubRegion != "" && !strings.EqualFold(c.Geo.SubRegion, country.Geo.SubRegion) {
			continue
		}

		// Misc
		//

		if c.InternationalPrefix != "" && !strings.EqualFold(c.InternationalPrefix, country.InternationalPrefix) {
			continue
		}

		// Bordering countries
		//

		allMatch := false

		if len(c.BorderingCountries()) > 0 {

			for _, c1 := range c.BorderingCountries() {

				match := false

				for _, c2 := range country.BorderingCountries() {
					match = c1.Alpha2 == c2.Alpha2

					if match {
						break
					}
				}

				if match {
					allMatch = true
				} else {
					allMatch = false
					break
				}

			}

			if !allMatch {
				continue
			}

		}

		// Append if all matches
		//

		countries = append(countries, country)

	}

	return
}

// FindSubdivisionCountryByName finds the country of a given subdivision name
func (q *Query) FindSubdivisionCountryByName(subdivisionName string) (result Country, err error) {
	for _, country := range q.Countries {
		if _, ok := country.nameToSubdivision[strings.ToLower(subdivisionName)]; ok {
			return country, nil
		}
	}

	return Country{}, makeError("Invalid subdivision name", subdivisionName)
}
