// Package humanize formats machine values into human-readable English.
//
// It is English-only, stateless, display-side helpers. It does not do
// internationalization, calendar math, unit conversion, or protocol fields.
// Humanized text is for display only; never feed it back into protocols,
// schedulers, billing, or audit.
//
// # API
//
//   - [Bytes] formats with decimal units (kB, MB, GB, ...).
//   - [BinaryBytes] formats with IEC binary units (KiB, MiB, GiB, ...).
//   - [Number] formats integers with comma separators.
//   - [Percent] takes a fraction: Percent(0.333) returns "33.3%".
//   - [Duration] shows at most two descending units: "1h 30m", "2d 5h".
//   - [Relative] formats a target time relative to a reference time.
//   - [Ordinal] adds an English ordinal suffix: "1st", "23rd".
//   - [Count] joins a count and an English noun: "1,000 items".
//
// # Relative time cutovers
//
// [Relative] uses approximate buckets, not calendar math. Months are 30 days,
// years are 365 days, days are 24 hours.
//
//	|diff| < 1s        just now
//	[1s,   60s)        N seconds ago / in N seconds
//	[60s,  60m)        N minutes ago / in N minutes
//	[60m,  24h)        N hours ago   / in N hours
//	[24h,   7d)        N days ago    / in N days
//	[7d,   30d)        N weeks ago   / in N weeks
//	[30d, 345d)        N months ago  / in N months
//	>= 345d            N years ago   / in N years
package humanize
